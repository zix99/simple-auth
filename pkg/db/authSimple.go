package db

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AccountAuthSimple struct {
	gorm.Model
	AccountID      uint
	Username       string
	PasswordBcrypt string
}

// VerifyPassword checks a password against the bcrypt entry
func (s *AccountAuthSimple) VerifyPassword(against string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.PasswordBcrypt), []byte(against)) != nil
}

// NewAccountAuthSimple creates a new account simple auth with a crypted password
func CreateAccountAuthSimple(belongsTo *Account, username, password string) (*AccountAuthSimple, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return nil, err
	}

	return &AccountAuthSimple{
		AccountID:      belongsTo.ID,
		Username:       username,
		PasswordBcrypt: string(hashed),
	}, nil
}
