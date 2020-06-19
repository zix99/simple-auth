package db

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type DB struct {
	db *gorm.DB
}

func New(driver string, args string) *DB {
	logrus.Infof("Connecting to %s at %s...", driver, args)

	db, err := gorm.Open(driver, args)
	if err != nil {
		logrus.Fatal(err)
	}

	db.AutoMigrate(&Account{})
	db.AutoMigrate(&AccountAuthSimple{})

	return &DB{db}
}

func (s *DB) IsAlive() bool {
	return s.db.DB().Ping() == nil
}
