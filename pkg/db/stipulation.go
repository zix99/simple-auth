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
	}

	return s.db.Create(stip).Error
}

func (s *sadb) findStipulations(account *Account, t StipulationType) ([]accountStipulation, error) {
	var stips []accountStipulation
	if err := s.db.Where("account_id = ? AND type = ?", account.ID, t).Find(&stips).Error; err != nil {
		return nil, err
	}
	return stips, nil
}

func (s *sadb) SatisfyStipulation(account *Account, satisfy IStipulation) error {
	if account == nil || satisfy == nil {
		return errors.New("invalid argument")
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
			err := s.db.Delete(&st).Error
			if err != nil {
				return err
			}
			s.CreateAuditRecord(account, AuditModuleAccount, AuditLevelInfo, "Stipulation %s is satisfied", spec.Type())
			return nil
		}
	}

	s.CreateAuditRecord(account, AuditModuleAccount, AuditLevelWarn, "Unable to validate stipulation %s", satisfy.Type())

	return errors.New("no stipulation satisfied")
}

func (s *sadb) AccountHasUnsatisfiedStipulations(account *Account) bool {
	var count int
	if err := s.db.Model(&accountStipulation{}).Where("account_id = ?", account.ID).Count(&count).Error; err != nil {
		return true
	}

	if count > 0 {
		return true
	}

	return false
}

func (s *sadb) ForceSatisfyStipulations(account *Account) error {
	return s.db.Where("account_id = ?", account.ID).Delete(&accountStipulation{}).Error
}
