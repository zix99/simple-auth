package db

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AccountStore interface {
	CreateAccount(email string) (*Account, error)
}

// Account represents a user
type Account struct {
	gorm.Model
	UUID   string `gorm:"type:varchar(64);unique_index;not null"`
	Email  string `gorm:"type:varchar(256);unique_index;not null"`
	Active bool   `gorm:"not null"`
}

func (s *sadb) CreateAccount(email string) (*Account, error) {
	account := &Account{
		UUID:   uuid.New().String(),
		Email:  email,
		Active: true,
	}
	if result := s.db.Create(&account); result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}
