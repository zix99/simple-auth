package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AccountAuthOneTime interface {
	CreateAccountOneTimeToken(account *Account, maxAge time.Duration) (string, error)
	AssertOneTimeToken(token string) (*Account, error)
}

type accountAuthOneTime struct {
	gorm.Model
	AccountID uint   `gorm:"index;not null"`
	Token     string `gorm:"index;not null"`
	Expires   time.Time
	Consumed  bool
}

func (s *sadb) CreateAccountOneTimeToken(account *Account, maxAge time.Duration) (string, error) {
	if account == nil {
		return "", InvalidAccount.New()
	}
	if !account.Active {
		return "", InactiveAccount.New()
	}

	token := accountAuthOneTime{
		AccountID: account.ID,
		Token:     uuid.New().String(),
		Expires:   time.Now().Add(maxAge),
		Consumed:  false,
	}

	if err := s.db.Create(&token).Error; err != nil {
		return "", err
	}

	s.CreateAuditRecord(account, AuditModuleOneTime, AuditLevelInfo, "One time token issued for account, expires in %s", maxAge.String())
	return token.Token, nil
}

func (s *sadb) AssertOneTimeToken(token string) (*Account, error) {
	if token == "" {
		return nil, SAOneTimeInvalidToken.New()
	}

	var oneTimeToken accountAuthOneTime
	if err := s.db.Where(&accountAuthOneTime{Token: token}).First(&oneTimeToken).Error; err != nil {
		return nil, err
	}

	if oneTimeToken.Consumed {
		return nil, SAOneTimeConsumed.New()
	}
	if time.Now().After(oneTimeToken.Expires) {
		return nil, SAOneTimeExpired.New()
	}

	// consume the token
	if err := s.db.Model(&oneTimeToken).Update(accountAuthOneTime{Consumed: true}).Error; err != nil {
		return nil, InternalError.Wrapf(err, "Error consuming token")
	}

	// Gain access to account
	var account Account
	if err := s.db.Model(&oneTimeToken).Related(&account).Error; err != nil {
		return nil, InternalError.Wrapf(err, "Unable to find account")
	}

	if !account.Active {
		return nil, InactiveAccount.New()
	}

	s.CreateAuditRecord(&account, AuditModuleOneTime, AuditLevelInfo, "One time token consumed for login")

	return &account, nil
}
