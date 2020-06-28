package db

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AccountAuthToken interface {
}

type accountAuthSessionToken struct {
	gorm.Model
	AccountID    uint   `gorm:"index;not null"`
	SessionToken string `gorm:"index"`
}

type accountAuthVerificationToken struct {
	gorm.Model
	AccountAuthSessionTokenID uint `gorm:"index; not null"`
}

// AssertCreateSessionToken checks the username, and password
// and upon acceptance, issues a session
func (db *sadb) AssertCreateSessionToken(username, password string) (string, error) {
	account, err := db.FindAndVerifySimpleAuth(username, password)
	if err != nil {
		return "", err
	}

	sessionToken := &accountAuthSessionToken{
		AccountID:    account.ID,
		SessionToken: uuid.New().String(),
	}
	db.db.Create(sessionToken)

	return sessionToken.SessionToken, nil
}

func (db *sadb) CreateVerificationToken(sessionToken string) (string, error) {
	return "", nil
}

func (db *sadb) AssertVerificationToken(username, verificationToken string) error {
	return nil
}
