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
	AccountAuthLocal
	AccountStore
	AccountAudit
	AccountOIDC
	AccountAuthOneTime
	AccountStipulations
	WithLogger(logger logrus.FieldLogger) SADB
	EnableLogging(enable bool)
	IsAlive() bool
	BeginTransaction() SADBTransaction
}

type SADBTransaction interface {
	SADB
	Commit() error
	Rollback() error
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
	db.AutoMigrate(&accountAuthLocal{})
	db.AutoMigrate(&accountAuthSessionToken{})
	db.AutoMigrate(&accountAuthVerificationToken{})
	db.AutoMigrate(&accountAuthOneTime{})
	db.AutoMigrate(&accountStipulation{})

	db.AutoMigrate(&accountOIDC{})
	db.Model(&accountOIDC{}).AddUniqueIndex("idx_provider_subject", "provider", "subject")

	return &sadb{db}
}

func (s *sadb) WithLogger(logger logrus.FieldLogger) SADB {
	wl := &sadb{
		s.db.New(),
	}
	wl.db.SetLogger(logger)
	return wl
}

func (s *sadb) IsAlive() bool {
	return s.db.DB().Ping() == nil
}

func (s *sadb) EnableLogging(enable bool) {
	s.db.LogMode(enable)
}

func (s *sadb) BeginTransaction() SADBTransaction {
	return &sadb{s.db.Begin()}
}

func (s *sadb) Commit() error {
	return s.db.Commit().Error
}

func (s *sadb) Rollback() error {
	return s.db.Rollback().Error
}
