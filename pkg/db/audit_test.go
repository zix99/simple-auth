package db_test

import (
	"simple-auth/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

var auditAccount *db.Account

func createAuditMock() {
	auditAccount, _ = sadb.CreateAccount("test", "audit-test@asdf.com")
}

func hasEvent(records []db.AccountAuditRecord, title string) bool {
	for _, val := range records {
		if val.Message == title {
			return true
		}
	}
	return false
}

func TestAddAudit(t *testing.T) {
	err := sadb.CreateAuditRecord(auditAccount, db.AuditModuleUI, db.AuditLevelInfo, "Test")
	assert.NoError(t, err)
}

func TestFetchAuditRecords(t *testing.T) {
	sadb.CreateAuditRecord(auditAccount, db.AuditModuleUI, db.AuditLevelInfo, "Test1")
	sadb.CreateAuditRecord(auditAccount, db.AuditModuleUI, db.AuditLevelInfo, "Test2")

	records, err := sadb.GetAuditTrailForAccount(auditAccount, 0, 5)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(records), 2)
	assert.True(t, hasEvent(records, "Test1"))
	assert.True(t, hasEvent(records, "Test2"))
}
