package db

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

type StipulationType string

type IStipulation interface {
	IsSatisfiedBy(data IStipulation) bool
	Type() StipulationType
}

type AccountStipulations interface {
	SatisfyStipulation(account *Account, satisfy IStipulation) error
	AddStipulation(account *Account, spec IStipulation) error
	AccountHasUnsatisfiedStipulations(account *Account) bool
	ForceSatisfyStipulations(account *Account) error
}

type accountStipulation struct {
	gorm.Model
	AccountID     uint `gorm:"index;not null"`
	Type          StipulationType
	Specification string
	Satisfied     bool
}

func (s *sadb) AddStipulation(account *Account, spec IStipulation) error {
	specBytes, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	stip := &accountStipulation{
		AccountID:     account.ID,
		Type:          spec.Type(),
		Specification: string(specBytes),
		Satisfied:     false,
	}

	return s.db.Create(stip).Error
}

func (s *sadb) findStipulations(account *Account, t StipulationType) ([]accountStipulation, error) {
	var stips []accountStipulation
	if err := s.db.Where("account_id = ? AND type = ? AND satisfied != 1", account.ID, t).Find(&stips).Error; err != nil {
		return nil, err
	}
	return stips, nil
}

func (s *sadb) SatisfyStipulation(account *Account, satisfy IStipulation) error {
	if account == nil || satisfy == nil {
		return errors.New("Invalid argument")
	}

	stips, err := s.findStipulations(account, satisfy.Type())
	if err != nil {
		return err
	}

	spec := reflect.New(reflect.ValueOf(satisfy).Elem().Type()).Interface().(IStipulation)

	for _, st := range stips {
		if err := json.Unmarshal([]byte(st.Specification), spec); err != nil {
			logrus.Error(err)
			continue
		}
		if spec.IsSatisfiedBy(satisfy) {
			err := s.db.Model(&st).Update(accountStipulation{Satisfied: true}).Error
			if err != nil {
				return err
			}
			s.CreateAuditRecord(account, AuditModuleAccount, AuditLevelInfo, "Stipulation %s is satisfied", spec.Type())
			return nil
		}
	}

	s.CreateAuditRecord(account, AuditModuleAccount, AuditLevelWarn, "Unable to validate stipulation %s", satisfy.Type())

	return errors.New("No stipulation satisfied")
}

func (s *sadb) AccountHasUnsatisfiedStipulations(account *Account) bool {
	var count int
	if err := s.db.Model(&accountStipulation{}).Where("account_id = ? and not satisfied", account.ID).Count(&count).Error; err != nil {
		return true
	}

	if count > 0 {
		return true
	}

	return false
}

func (s *sadb) ForceSatisfyStipulations(account *Account) error {
	return s.db.Model(&accountStipulation{AccountID: account.ID}).Update(&accountStipulation{Satisfied: true}).Error
}
