package db

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type sadb struct {
	db *gorm.DB
}

type SADB interface {
	AccountAuthToken
	AccountAuthSimple
	AccountStore
	AccountAudit
	EnableLogging(enable bool)
	IsAlive() bool
}

func New(driver string, args string) SADB {
	logrus.Infof("Connecting to %s at %s...", driver, args)

	db, err := gorm.Open(driver, args)
	if err != nil {
		logrus.Fatal(err)
	}

	db.SetLogger(logrus.StandardLogger())

	db.AutoMigrate(&Account{})
	db.AutoMigrate(&AccountAuditRecord{})
	db.AutoMigrate(&accountAuthSimple{})
	db.AutoMigrate(&accountAuthSessionToken{})
	db.AutoMigrate(&accountAuthVerificationToken{})

	return &sadb{db}
}

func (s *sadb) IsAlive() bool {
	return s.db.DB().Ping() == nil
}

func (s *sadb) EnableLogging(enable bool) {
	s.db.LogMode(enable)
}
