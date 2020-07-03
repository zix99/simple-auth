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

// verifyPassword checks a password against the bcrypt entry
func (s *accountAuthSimple) verifyPassword(against string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.PasswordBcrypt), []byte(against)) == nil
}

// CreateAccountAuthSimple creates a new account simple auth with a crypted password
func (db *sadb) CreateAccountAuthSimple(belongsTo *Account, username, password string) error {
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

	db.db.Create(auth)
	return nil
}

func (db *sadb) FindAndVerifySimpleAuth(username, password string) (*Account, error) {
	if username == "" || password == "" {
		return nil, errors.New("Invalid arg")
	}

	var auth accountAuthSimple
	if err := db.db.Where(&accountAuthSimple{Username: username}).First(&auth).Error; err != nil {
		return nil, err
	}

	if !auth.verifyPassword(password) {
		return nil, UserVerificationFailed
	}

	var account Account
	if err := db.db.Model(&auth).Related(&account).Error; err != nil {
		return nil, UserVerificationFailed
	}

	if !account.Active {
		return nil, errors.New("Account Inactive")
	}

	return &account, nil
}

func (s *sadb) FindAccountForSimpleAuth(username string) (*Account, error) {
	if username == "" {
		return nil, errors.New("Invalid username")
	}

	var simpleAuth accountAuthSimple
	if err := s.db.Where(&accountAuthSimple{Username: username}).First(&simpleAuth).Error; err != nil {
		return nil, err
	}

	var account Account
	if err := s.db.Model(&simpleAuth).Related(&account).Error; err != nil {
		return nil, err
	}

	if !account.Active {
		return nil, errors.New("Account Inactive")
	}

	return &account, nil
}
