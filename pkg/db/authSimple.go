package db

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AccountAuthSimple interface {
	CreateAccountAuthSimple(belongsTo *Account, username, password string) error
	FindAndVerifySimpleAuth(username, password string) (*Account, error)
	FindAccountForSimpleAuth(username string) (*Account, error)
	FindSimpleAuthUsername(account *Account) (string, error)
	UpdatePasswordForUsername(username string, newPassword string) error
}

type accountAuthSimple struct {
	gorm.Model
	AccountID      uint   `gorm:"index;not null"`
	Username       string `gorm:"type:varchar(256);unique_index;not null"`
	PasswordBcrypt string `gorm:"not null"`
}

var ErrorAccountInactive = errors.New("Inactive Account")

// verifyPassword checks a password against the bcrypt entry
func (s *accountAuthSimple) verifyPassword(against string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.PasswordBcrypt), []byte(against)) == nil
}

// CreateAccountAuthSimple creates a new account simple auth with a crypted password
func (s *sadb) CreateAccountAuthSimple(belongsTo *Account, username, password string) error {
	if belongsTo == nil {
		return errors.New("Invalid account")
	}
	if !belongsTo.Active {
		return errors.New("Unable to associate with deactivated account")
	}

	username = strings.TrimSpace(strings.ToLower(username))
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return errors.New("Invalid username or password")
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

	return s.db.Create(auth).Error
}

func (s *sadb) UpdatePasswordForUsername(username string, newPassword string) error {
	auth, _, err := s.resolveSimpleAuthForUser(username)
	if err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 0)
	if err != nil {
		return err
	}

	return s.db.Model(auth).Update(accountAuthSimple{PasswordBcrypt: string(hashed)}).Error
}

func (s *sadb) resolveSimpleAuthForUser(username string) (*accountAuthSimple, *Account, error) {
	if username == "" {
		return nil, nil, errors.New("Invalid username")
	}
	username = strings.ToLower(username)

	var simpleAuth accountAuthSimple
	if err := s.db.Where(&accountAuthSimple{Username: username}).First(&simpleAuth).Error; err != nil {
		return nil, nil, err
	}

	var account Account
	if err := s.db.Model(&simpleAuth).Related(&account).Error; err != nil {
		return nil, nil, err
	}

	if !account.Active {
		return nil, nil, ErrorAccountInactive
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

func (s *sadb) FindAndVerifySimpleAuth(username, password string) (*Account, error) {
	if username == "" || password == "" {
		return nil, errors.New("Invalid arg")
	}

	auth, account, err := s.resolveSimpleAuthForUser(username)
	if err != nil {
		return nil, err
	}

	if !auth.verifyPassword(password) {
		s.CreateAuditRecord(account, AuditModuleSimple, AuditLevelWarn, "Login failed")
		return nil, UserVerificationFailed
	}

	s.CreateAuditRecord(account, AuditModuleSimple, AuditLevelInfo, "Login Successful")

	return account, nil
}
