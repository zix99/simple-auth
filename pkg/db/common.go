package db

import (
	"errors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var UserNotFound = errors.New("User not found")

var UserVerificationFailed = errors.New("User verification failed")
