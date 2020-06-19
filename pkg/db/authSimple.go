package db

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AccountAuthSimple interface {
	CreateAccountAuthSimple(belongsTo *Account, username, password string) error
	VerifySimpleAuth(account *Account, username, password string) error
	FindAndVerifySimpleAuth(username, password string) (*Account, error)
}

type accountAuthSimple struct {
	gorm.Model
	AccountID      uint
	Username       string
	PasswordBcrypt string
}

// verifyPassword checks a password against the bcrypt entry
func (s *accountAuthSimple) verifyPassword(against string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.PasswordBcrypt), []byte(against)) != nil
}

// CreateAccountAuthSimple creates a new account simple auth with a crypted password
func (db *sadb) CreateAccountAuthSimple(belongsTo *Account, username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}

	auth := &accountAuthSimple{
		AccountID:      belongsTo.ID,
		Username:       username,
		PasswordBcrypt: string(hashed),
	}

	db.db.Create(&auth)
	return nil
}

func (db *sadb) VerifySimpleAuth(account *Account, username, password string) error {
	var auth accountAuthSimple
	if result := db.db.Model(account).Related(&auth); result.Error != nil {
		return UserNotFound
	}
	if auth.Username != username || !auth.verifyPassword(password) {
		return UserVerificationFailed
	}
	return nil
}

func (db *sadb) FindAndVerifySimpleAuth(username, password string) (*Account, error) {
	var auth accountAuthSimple
	if result := db.db.Where(&accountAuthSimple{Username: username}).First(&auth); result.Error != nil {
		return nil, result.Error
	}

	if !auth.verifyPassword(password) {
		return nil, UserVerificationFailed
	}

	var account Account
	db.db.Model(&auth).Related(&account) // TODO: Error check
	return &account, nil
}
