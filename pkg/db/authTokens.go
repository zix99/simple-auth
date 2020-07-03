package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// AccountAuthToken exposes token authentication against simple auth
type AccountAuthToken interface {
	AssertCreateSessionToken(username, password string, expires time.Duration) (string, error)
	InvalidateSession(sessionToken string) error
	CreateVerificationToken(username, sessionToken string) (string, error)
	AssertVerificationToken(username, verificationToken string) (*Account, error)
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
	AccountID                 uint   `gorm:"index; not null"`
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
	s.db.Where(accountAuthSessionToken{
		AccountID:   account.ID,
		Invalidated: false,
	}).Updates(accountAuthSessionToken{Invalidated: true})

	// Create new session tokens
	sessionToken := &accountAuthSessionToken{
		AccountID:    account.ID,
		SessionToken: uuid.New().String(),
		Expires:      time.Now().UTC().Add(expires),
		Invalidated:  false,
	}
	s.db.Create(sessionToken)

	return sessionToken.SessionToken, nil
}

func (s *sadb) InvalidateSession(sessionToken string) error {
	return s.db.Where(accountAuthSessionToken{
		SessionToken: sessionToken,
		Invalidated:  false,
	}).Updates(accountAuthSessionToken{Invalidated: true}).Error
}

// CreateVerificationToken takes a session token and converts it into a short-lived verification token
func (s *sadb) CreateVerificationToken(username, sessionToken string) (string, error) {
	account, err := s.FindAccountForSimpleAuth(username)
	if err != nil {
		return "", fmt.Errorf("Account not found: %w", err)
	}

	var session accountAuthSessionToken
	logrus.Infof("Looking for %v(%v) and token=%v", username, account, sessionToken)
	if err := s.db.Where("account_id = ? AND session_token = ? AND not invalidated", account.ID, sessionToken).First(&session).Error; err != nil {
		return "", fmt.Errorf("Session not found: %w", err)
	}

	if session.Invalidated {
		return "", errors.New("Session invalidated")
	}
	if time.Now().After(session.Expires) {
		return "", errors.New("Session expired")
	}

	verificationToken := &accountAuthVerificationToken{
		AccountID:                 session.AccountID,
		AccountAuthSessionTokenID: session.ID,
		VerificationToken:         uuid.New().String(),
		Consumed:                  false,
		Expires:                   time.Now().UTC().Add(10 * time.Second),
	}
	s.db.Create(verificationToken)

	return verificationToken.VerificationToken, nil
}

func (s *sadb) AssertVerificationToken(username, verificationToken string) (*Account, error) {
	account, err := s.FindAccountForSimpleAuth(username)
	if err != nil {
		return nil, err
	}

	var token accountAuthVerificationToken
	if err := s.db.Where("account_id = ? AND verification_token = ?", account.ID, verificationToken).First(&token).Error; err != nil {
		return nil, err
	}

	if token.Consumed {
		return nil, errors.New("Verification token already consumed")
	}

	if time.Now().After(token.Expires) {
		return nil, errors.New("Verification token expired")
	}

	if err := s.db.Model(token).Update(accountAuthVerificationToken{Consumed: true}).Error; err != nil {
		return nil, err
	}

	if token.VerificationToken != verificationToken {
		return nil, errors.New("Invalid verification token")
	}

	return account, nil
}
