package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type AccountAudit interface {
	CreateAuditRecord(account AccountProvider, module AuditModule, level AuditLevel, message string, params ...interface{}) error
	GetAuditTrailForAccount(account *Account, offset, count int) ([]AccountAuditRecord, error)
}

type (
	AuditLevel  string
	AuditModule string
)

const (
	AuditLevelDebug = "debug"
	AuditLevelInfo  = "info"
	AuditLevelWarn  = "warn"
	AuditLevelAlert = "alert"
)

const (
	AuditModuleAccount = "account"
	AuditModuleUI      = "ui"
	AuditModuleLocal   = "auth:simple"
	AuditModuleToken   = "auth:token"
	AuditModuleOAuth2  = "auth:oauth2"
	AuditModuleOIDC    = "login:oidc"
	AuditModuleOneTime = "auth:onetime"
)

type AccountAuditRecord struct {
	gorm.Model
	AccountID uint `form:"index;not null"`
	Module    AuditModule
	Level     AuditLevel
	Message   string
}

func (s *sadb) CreateAuditRecord(account AccountProvider, module AuditModule, level AuditLevel, message string, params ...interface{}) error {
	record := &AccountAuditRecord{
		AccountID: account.Account().ID,
		Module:    module,
		Level:     level,
		Message:   fmt.Sprintf(message, params...),
	}
	err := s.db.Create(record).Error
	if err != nil {
		logrus.Warnf("Failed to create audit log: %s", err)
	}
	return err
}

func (s *sadb) GetAuditTrailForAccount(account *Account, offset, count int) ([]AccountAuditRecord, error) {
	var auditRecords []AccountAuditRecord
	err := s.db.Model(account).Order("created_at desc").Offset(offset).Limit(count).Related(&auditRecords).Error
	if err != nil {
		return nil, err
	}
	return auditRecords, nil
}
