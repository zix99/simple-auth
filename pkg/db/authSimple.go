package db

import (
	"errors"
	"simple-auth/pkg/lib/totp"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AccountAuthSimple interface {
	// Safe function
	AssertSimpleAuth(username, password string, totpCode *string) (*Account, error)

	// Unsafe functions
	CreateAccountAuthSimple(belongsTo *Account, username, password string) error
	FindAccountForSimpleAuth(username string) (*Account, error)
	FindSimpleAuthUsername(account *Account) (string, error)
	UpdatePasswordForUsername(username string, newPassword string) error

	// TOTP
	HasTOTPEnabled(account *Account) bool
	SetAuthSimpleTOTP(account *Account, totpURL *string) error
	ValidateTOTP(account *Account, totpCode string) bool
}

type accountAuthSimple struct {
	gorm.Model
	AccountID      uint   `gorm:"index;not null"`
	Username       string `gorm:"type:varchar(256);unique_index;not null"`
	PasswordBcrypt string `gorm:"not null"`
	TOTPSpec       *string
}

// verifyPassword checks a password against the bcrypt entry
func (s *accountAuthSimple) verifyPassword(against string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.PasswordBcrypt), []byte(against)) == nil
}

// CreateAccountAuthSimple creates a new account simple auth with a crypted password
func (s *sadb) CreateAccountAuthSimple(belongsTo *Account, username, password string) error {
	if belongsTo == nil {
		return InvalidAccount.New()
	}
	if !belongsTo.Active {
		return InactiveAccount.Newf("Unable to associate with deactivated account")
	}

	username = strings.TrimSpace(strings.ToLower(username))
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return SAInvalidCredentials.New()
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}

	auth := &accountAuthSimple{
		AccountID:      belongsTo.ID,
		Username:       username,
		PasswordBcrypt: string(hashed),
	}

	s.CreateAuditRecord(belongsTo, AuditModuleSimple, AuditLevelInfo, "Associated username: %s", username)

	return s.db.Create(auth).Error
}

func (s *sadb) UpdatePasswordForUsername(username string, newPassword string) error {
	auth, account, err := s.resolveSimpleAuthForUser(username)
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 0)
	if err != nil {
		return InternalError.Wrap(err)
	}

	s.CreateAuditRecord(account, AuditModuleSimple, AuditLevelInfo, "Password updated")

	err = s.db.Model(auth).Update(accountAuthSimple{PasswordBcrypt: string(hashed)}).Error
	if err != nil {
		return InternalError.Wrap(err)
	}
	return nil
}

func (s *sadb) resolveSimpleAuthForUser(username string) (*accountAuthSimple, *Account, error) {
	if username == "" {
		return nil, nil, SAInvalidCredentials.Newf("Missing username")
	}
	username = strings.ToLower(username)

	var simpleAuth accountAuthSimple
	if err := s.db.Where(&accountAuthSimple{Username: username}).First(&simpleAuth).Error; err != nil {
		return nil, nil, SAInvalidCredentials.Wrap(err)
	}

	var account Account
	if err := s.db.Model(&simpleAuth).Related(&account).Error; err != nil {
		return nil, nil, InternalError.Wrap(err)
	}

	if !account.Active {
		return nil, nil, InactiveAccount.New()
	}

	return &simpleAuth, &account, nil
}

func (s *sadb) FindSimpleAuthUsername(account *Account) (string, error) {
	if account == nil {
		return "", errors.New("No account")
	}

	var simpleAuth accountAuthSimple
	if err := s.db.Model(account).Related(&simpleAuth).Error; err != nil {
		return "", errors.New("No simple-auth linked")
	}

	return simpleAuth.Username, nil
}

func (s *sadb) FindAccountForSimpleAuth(username string) (*Account, error) {
	_, account, err := s.resolveSimpleAuthForUser(username)
	return account, err
}

func (s *sadb) AssertSimpleAuth(username, password string, totpCode *string) (*Account, error) {
	if username == "" || password == "" {
		return nil, SAInvalidCredentials.Newf("Invalid username/password")
	}

	auth, account, err := s.resolveSimpleAuthForUser(username)
	if err != nil {
		return nil, SAInvalidCredentials.Wrap(err)
	}

	// Password
	if !auth.verifyPassword(password) {
		s.CreateAuditRecord(account, AuditModuleSimple, AuditLevelWarn, "Login failed")
		return nil, SAInvalidCredentials.New()
	}

	// Stipulations
	if s.AccountHasUnsatisfiedStipulations(account) {
		return nil, SAUnsatisfiedStipulations.New()
	}

	// TOTP
	if auth.TOTPSpec != nil {
		if totpCode == nil || *totpCode == "" {
			return nil, SATOTPMissing.New()
		}

		otp, err := totp.ParseTOTP(*auth.TOTPSpec)
		if err != nil {
			return nil, SATOTPFailed.Wrap(err)
		}
		if !otp.Validate(*totpCode, 1) {
			s.CreateAuditRecord(account, AuditModuleSimple, AuditLevelWarn, "TOTP Rejected")
			return nil, SATOTPFailed.New()
		}
	}

	// Success
	s.CreateAuditRecord(account, AuditModuleSimple, AuditLevelInfo, "Login Successful")

	return account, nil
}

func (s *sadb) resolveSimpleAuthForAccount(account *Account) (*accountAuthSimple, error) {
	if account == nil {
		return nil, InvalidAccount.New()
	}
	if !account.Active {
		return nil, InactiveAccount.New()
	}

	var auth accountAuthSimple
	if err := s.db.Model(account).Related(&auth).Error; err != nil {
		return nil, err
	}

	return &auth, nil
}

func (s *sadb) SetAuthSimpleTOTP(account *Account, totpURL *string) error {
	simpleAuth, err := s.resolveSimpleAuthForAccount(account)
	if err != nil {
		return err
	}

	// Disable
	if totpURL == nil {
		s.CreateAuditRecord(account, AuditModuleSimple, AuditLevelInfo, "Disabled TOTP")
		return s.db.Model(&simpleAuth).Update("TOTPSpec", nil).Error
	}

	s.CreateAuditRecord(account, AuditModuleSimple, AuditLevelInfo, "Activated TOTP")

	return s.db.Model(&simpleAuth).Update(accountAuthSimple{TOTPSpec: totpURL}).Error
}

func (s *sadb) HasTOTPEnabled(account *Account) bool {
	simpleAuth, err := s.resolveSimpleAuthForAccount(account)
	if err != nil {
		return false
	}
	return simpleAuth.TOTPSpec != nil
}

func (s *sadb) ValidateTOTP(account *Account, code string) bool {
	auth, err := s.resolveSimpleAuthForAccount(account)
	if err != nil {
		return false
	}

	if auth.TOTPSpec == nil {
		return true
	}

	if code == "" {
		return false
	}

	otp, err := totp.ParseTOTP(*auth.TOTPSpec)
	if err != nil {
		return false
	}
	return otp.Validate(code, 1) // FIXME: Drift
}
