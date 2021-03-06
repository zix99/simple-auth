package db

import (
	"errors"
	"simple-auth/pkg/lib/totp"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AccountAuthLocal interface {
	FindAuthLocal(account *Account) (*AuthLocal, error)
	FindAuthLocalByUsername(username string) (*AuthLocal, error)
	FindAuthLocalByEmail(email string) (*AuthLocal, error)

	CreateAuthLocal(belongsTo *Account, username, password string) (*AuthLocal, error)

	UpdateAuthLocalPassword(authLocal *AuthLocal, newPassword string) error
	UpdateAuthLocalTOTP(authLocal *AuthLocal, totpURL *string) error
}

type accountAuthLocal struct {
	gorm.Model
	AccountID      uint   `gorm:"index;not null"`
	Username       string `gorm:"type:varchar(256);unique_index;not null"`
	PasswordBcrypt string `gorm:"not null"`
	TOTPSpec       *string
}

type AuthLocal struct {
	auth    *accountAuthLocal
	account *Account
}

func (s *AuthLocal) Username() string {
	return s.auth.Username
}

func (s *AuthLocal) VerifyPassword(against string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.auth.PasswordBcrypt), []byte(against)) == nil
}

func (s *AuthLocal) VerifyTOTP(against string, drift int) bool {
	if s.auth.TOTPSpec == nil {
		return true
	}

	tfa, err := totp.ParseTOTP(*s.auth.TOTPSpec)
	if err != nil {
		return false
	}

	return tfa.Validate(against, drift)
}

func (s *AuthLocal) HasTOTP() bool {
	return s.auth.TOTPSpec != nil
}

func (s *AuthLocal) Account() *Account {
	return s.account
}

func (s *sadb) FindAuthLocal(account *Account) (*AuthLocal, error) {
	if account == nil {
		return nil, errors.New("no account")
	}

	var localAuth accountAuthLocal
	if err := s.db.Model(account).Related(&localAuth).Error; err != nil {
		return nil, errors.New("no local-auth linked")
	}

	return &AuthLocal{
		auth:    &localAuth,
		account: account,
	}, nil
}

func (s *sadb) FindAuthLocalByUsername(username string) (*AuthLocal, error) {
	if username == "" {
		return nil, AuthInvalidUsername.Newf("Empty username")
	}
	username = strings.ToLower(username)

	var localAuth accountAuthLocal
	if err := s.db.Where(&accountAuthLocal{Username: username}).First(&localAuth).Error; err != nil {
		return nil, AuthInvalidUsername.Wrap(err)
	}

	var account Account
	if err := s.db.Model(&localAuth).Related(&account).Error; err != nil {
		return nil, InternalError.Wrap(err)
	}

	return &AuthLocal{
		auth:    &localAuth,
		account: &account,
	}, nil
}

func (s *sadb) FindAuthLocalByEmail(email string) (*AuthLocal, error) {
	account, err := s.FindAccountByEmail(email)
	if err != nil {
		return nil, InvalidAccount.Wrap(err)
	}

	var localAuth accountAuthLocal
	if err := s.db.Model(account).Related(&localAuth).Error; err != nil {
		return nil, InvalidAccount.Wrap(err)
	}

	return &AuthLocal{
		auth:    &localAuth,
		account: account,
	}, nil
}

// CreateAuthLocal creates a new account local auth with a crypted password
func (s *sadb) CreateAuthLocal(belongsTo *Account, username, password string) (*AuthLocal, error) {
	if belongsTo == nil {
		return nil, InvalidAccount.New()
	}
	if !belongsTo.Active {
		return nil, InactiveAccount.Newf("Unable to associate with deactivated account")
	}

	username = strings.TrimSpace(strings.ToLower(username))
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return nil, AuthInvalidUsername.New()
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return nil, err
	}

	auth := &accountAuthLocal{
		AccountID:      belongsTo.ID,
		Username:       username,
		PasswordBcrypt: string(hashed),
	}

	s.CreateAuditRecord(belongsTo, AuditModuleLocal, AuditLevelInfo, "Associated username: %s", username)

	if err := s.db.Create(auth).Error; err != nil {
		return nil, err
	}
	return &AuthLocal{
		auth:    auth,
		account: belongsTo,
	}, nil
}

func (s *sadb) UpdateAuthLocalPassword(authLocal *AuthLocal, newPassword string) error {
	if authLocal == nil {
		return InternalError.Newf("Auth nil")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 0)
	if err != nil {
		return InternalError.Wrap(err)
	}

	s.CreateAuditRecord(authLocal, AuditModuleLocal, AuditLevelInfo, "Password updated")

	err = s.db.Model(authLocal.auth).Update(accountAuthLocal{PasswordBcrypt: string(hashed)}).Error
	if err != nil {
		return InternalError.Wrap(err)
	}
	return nil
}

func (s *sadb) UpdateAuthLocalTOTP(authLocal *AuthLocal, totpURL *string) error {
	if authLocal == nil {
		return InternalError.Newf("Auth nil")
	}

	// Disable
	if totpURL == nil {
		s.CreateAuditRecord(authLocal, AuditModuleLocal, AuditLevelInfo, "Disabled TOTP")
		return s.db.Model(authLocal.auth).Update("TOTPSpec", nil).Error
	}

	s.CreateAuditRecord(authLocal, AuditModuleLocal, AuditLevelInfo, "Activated TOTP")

	return s.db.Model(authLocal.auth).Update(accountAuthLocal{TOTPSpec: totpURL}).Error
}
