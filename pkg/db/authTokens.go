package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AccountAuthToken interface {
}

type accountAuthSessionToken struct {
	gorm.Model
	AccountID    uint   `gorm:"index;not null"`
	SessionToken string `gorm:"index"`
	Invalidated  bool
	Expires      time.Time
}

type accountAuthVerificationToken struct {
	gorm.Model
	AccountAuthSessionTokenID uint   `gorm:"index; not null"`
	VerificationToken         string `gorm:"index"`
	Consumed                  bool
	Expires                   time.Time
}

// AssertCreateSessionToken checks the username, and password
// and upon acceptance, issues a session token
func (s *sadb) AssertCreateSessionToken(username, password string, expires time.Duration) (string, error) {
	account, err := s.FindAndVerifySimpleAuth(username, password)
	if err != nil {
		return "", err
	}

	// Invalidate all existing tokens
	s.db.Where(accountAuthSessionToken{AccountID: account.ID, Invalidated: false}).Updates(accountAuthSessionToken{Invalidated: true})

	// Create new session tokens
	sessionToken := &accountAuthSessionToken{
		AccountID:    account.ID,
		SessionToken: uuid.New().String(),
		Expires:      time.Now().UTC().Add(expires),
	}
	s.db.Create(sessionToken)

	return sessionToken.SessionToken, nil
}

// CreateVerificationToken takes a session token and converts it into a verification token
func (s *sadb) CreateVerificationToken(sessionToken string) (string, error) {
	var session accountAuthSessionToken
	if err := s.db.Where("SessionToken = ? AND not Invalidated", sessionToken).Order("created_at desc").First(&session).Error; err != nil {
		return "", err
	}

	verificationToken := accountAuthVerificationToken{
		AccountAuthSessionTokenID: session.ID,
		VerificationToken:         uuid.New().String(),
		Consumed:                  false,
		Expires:                   time.Now().UTC().Add(10 * time.Second),
	}
	s.db.Create(verificationToken)

	return verificationToken.VerificationToken, nil
}

func (s *sadb) AssertVerificationToken(username, verificationToken string) error {
	return nil
}
