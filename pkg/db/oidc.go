package db

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type AccountOIDC interface {
	FindAccountForOIDC(provider, subject string) (*Account, error)
	CreateOIDCForAccount(account *Account, provider, subject string) error
	FindOIDCForAccount(account *Account) ([]OIDCDescriptor, error)
}

// NOTE: extra index created in db.go
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

type OIDCDescriptor struct {
	Provider string
	Subject  string
}

func (s *sadb) FindOIDCForAccount(account *Account) ([]OIDCDescriptor, error) {
	var providers []*accountOIDC
	err := s.db.Model(account).Related(&providers).Error
	if err != nil {
		return nil, err
	}

	ret := make([]OIDCDescriptor, len(providers))
	for i, provider := range providers {
		ret[i] = OIDCDescriptor{
			Provider: provider.Provider,
			Subject:  provider.Subject,
		}
	}
	return ret, nil
}

func (s *sadb) CreateOIDCForAccount(account *Account, provider, subject string) error {
	if account == nil {
		return errors.New("invalid account")
	}
	if provider == "" || subject == "" {
		return errors.New("bad arguments")
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
