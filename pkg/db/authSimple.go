package db

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AccountAuthSimple interface {
	CreateAccountAuthSimple(belongsTo *Account, username, password string) error
	FindAndVerifySimpleAuth(username, password string) (*Account, error)
	FindAccountForSimpleAuth(username string) (*Account, error)
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
	if !belongsTo.Active {
		return errors.New("Unable to associate with deactivated account")
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

	s.db.Create(auth)
	return nil
}

func (s *sadb) resolveSimpleAuthForUser(username string) (*accountAuthSimple, *Account, error) {
	if username == "" {
		return nil, nil, errors.New("Invalid username")
	}

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
		return nil, UserVerificationFailed
	}

	return account, nil
}
