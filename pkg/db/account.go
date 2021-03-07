package db

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AccountStore interface {
	CreateAccount(name, email string) (*Account, error)
	FindAccount(uuid string) (*Account, error)
	FindAccountByEmail(email string) (*Account, error)

	GetAllAccounts(itr func(account *Account) bool) error
}

type AccountProvider interface {
	Account() *Account
}

// Account represents a user
type Account struct {
	gorm.Model
	UUID   string `gorm:"type:varchar(64);unique_index;not null"`
	Name   string `gorm:"type:varchar(256);not null"`
	Email  string `gorm:"type:varchar(256);unique_index;not null"`
	Active bool   `gorm:"not null"`
}

func (s *Account) Account() *Account {
	return s
}

func (s *sadb) CreateAccount(name, email string) (*Account, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return nil, errors.New("invalid email")
	}

	account := &Account{
		UUID:   uuid.New().String(),
		Name:   name,
		Email:  email,
		Active: true,
	}

	if result := s.db.Create(&account); result.Error != nil {
		return nil, result.Error
	}

	s.CreateAuditRecord(account, AuditModuleAccount, AuditLevelInfo, "Account created")

	return account, nil
}

func (s *sadb) FindAccount(uuid string) (*Account, error) {
	if uuid == "" {
		return nil, errors.New("missing UUID")
	}
	var account Account
	err := s.db.Where(&Account{UUID: uuid}).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, err
}

func (s *sadb) FindAccountByEmail(email string) (*Account, error) {
	if email == "" {
		return nil, errors.New("missing email")
	}
	var account Account
	err := s.db.Where(&Account{Email: email}).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, err
}

func (s *sadb) GetAllAccounts(itr func(account *Account) bool) error {
	rows, err := s.db.Model(&Account{}).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var account Account
		s.db.ScanRows(rows, &account)
		if itr(&account) {
			break
		}
	}

	return nil
}
