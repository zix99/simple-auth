package db

import (
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
	account, err := s.AssertSimpleAuth(username, password, nil)
	if err != nil {
		return "", err
	}

	// Invalidate all existing tokens
	err = s.db.Model(&accountAuthSessionToken{}).Where("account_id = ? and not invalidated", account.ID).Update("invalidated", true).Error
	if err != nil {
		return "", InternalError.Wrapf(err, "Error invalidating tokens")
	}

	// Create new session tokens
	sessionToken := &accountAuthSessionToken{
		AccountID:    account.ID,
		SessionToken: uuid.New().String(),
		Expires:      time.Now().UTC().Add(expires),
		Invalidated:  false,
	}
	if err := s.db.Create(sessionToken).Error; err != nil {
		return "", InternalError.Wrap(err)
	}

	s.CreateAuditRecord(account, AuditModuleToken, AuditLevelInfo, "Created session token")

	return sessionToken.SessionToken, nil
}

func (s *sadb) InvalidateSession(sessionToken string) error {
	return s.db.Model(&accountAuthSessionToken{}).Where(&accountAuthSessionToken{
		SessionToken: sessionToken,
		Invalidated:  false,
	}).Updates(accountAuthSessionToken{Invalidated: true}).Error
}

// CreateVerificationToken takes a session token and converts it into a short-lived verification token
func (s *sadb) CreateVerificationToken(username, sessionToken string) (string, error) {
	account, err := s.FindAccountForSimpleAuth(username)
	if err != nil {
		return "", InvalidAccount.Wrapf(err, "Account not found")
	}

	var session accountAuthSessionToken
	logrus.Infof("Looking for %v(%v) and token=%v", username, account, sessionToken)
	if err := s.db.Where("account_id = ? AND session_token = ? AND not invalidated", account.ID, sessionToken).First(&session).Error; err != nil {
		s.CreateAuditRecord(account, AuditModuleToken, AuditLevelWarn, "Failed to create verification token on undefined session")
		return "", SessionNotFound.Wrap(err)
	}

	if session.Invalidated {
		s.CreateAuditRecord(account, AuditModuleToken, AuditLevelWarn, "Failed to create verification token on invalidated session")
		return "", SessionInvalidated.New()
	}
	if time.Now().After(session.Expires) {
		s.CreateAuditRecord(account, AuditModuleToken, AuditLevelWarn, "Failed to create verification token on expired session")
		return "", SessionExpired.New()
	}

	verificationToken := &accountAuthVerificationToken{
		AccountID:                 session.AccountID,
		AccountAuthSessionTokenID: session.ID,
		VerificationToken:         uuid.New().String(),
		Consumed:                  false,
		Expires:                   time.Now().UTC().Add(10 * time.Second),
	}
	if err := s.db.Create(verificationToken).Error; err != nil {
		return "", InternalError.Wrap(err)
	}

	s.CreateAuditRecord(account, AuditModuleToken, AuditLevelInfo, "Created verification token")

	return verificationToken.VerificationToken, nil
}

func (s *sadb) AssertVerificationToken(username, verificationToken string) (*Account, error) {
	account, err := s.FindAccountForSimpleAuth(username)
	if err != nil {
		return nil, err
	}

	var token accountAuthVerificationToken
	if err := s.db.Where("account_id = ? AND verification_token = ?", account.ID, verificationToken).First(&token).Error; err != nil {
		return nil, VerificationMissing.Wrap(err)
	}

	if token.Consumed {
		return nil, VerificationConsumed.Newf("Verification token already consumed")
	}

	if time.Now().After(token.Expires) {
		return nil, VerificationExpired.Newf("Verification token expired")
	}

	if err := s.db.Model(token).Update(accountAuthVerificationToken{Consumed: true}).Error; err != nil {
		return nil, InternalError.Wrap(err)
	}

	if token.VerificationToken != verificationToken {
		return nil, VerificationInvalid.Newf("Invalid verification token")
	}

	s.CreateAuditRecord(account, AuditModuleToken, AuditLevelDebug, "Verification Token Validated")

	return account, nil
}
