package db

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type AccountOIDC interface {
	FindAccountForOIDC(provider, subject string) (*Account, error)
	CreateOIDCForAccount(account *Account, provider, subject string) error
}

type accountOIDC struct {
	gorm.Model
	AccountID uint   `gorm:"index;not null"`
	Provider  string `gorm:"not null"`
	Subject   string `gorm:"not null"`
}

func (s *sadb) FindAccountForOIDC(provider, subject string) (*Account, error) {
	var oidc accountOIDC
	if err := s.db.Where(&accountOIDC{Provider: provider, Subject: subject}).First(&oidc).Error; err != nil {
		return nil, err
	}

	var account Account
	if err := s.db.Model(&oidc).Related(&account).Error; err != nil {
		return nil, err
	}

	s.CreateAuditRecord(&account, AuditModuleOIDC, AuditLevelInfo, "OIDC lookup succeeded for %s", provider)

	return &account, nil
}

func (s *sadb) CreateOIDCForAccount(account *Account, provider, subject string) error {
	if account == nil {
		return errors.New("Invalid account")
	}
	if provider == "" || subject == "" {
		return errors.New("Bad arguments")
	}

	oidc := &accountOIDC{
		AccountID: account.ID,
		Provider:  provider,
		Subject:   subject,
	}
	err := s.db.Create(oidc).Error
	if err == nil {
		s.CreateAuditRecord(account, AuditModuleOIDC, AuditLevelInfo, "OIDC provider %s linked to account", provider)
	}
	return err
}
