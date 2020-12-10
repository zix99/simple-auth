package db_test

import (
	"simple-auth/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sadb db.SADB

func init() {
	sadb = db.New("sqlite3", "file::memory:?cache=shared")
	sadb.EnableLogging(true)
	createAuthSimpleMock()
	createAuthTokenMock()
	createAuditMock()
	createOIDCMock()
}

func TestIsAlive(t *testing.T) {
	assert.True(t, sadb.IsAlive())
}

func TestCommit(t *testing.T) {
	trx := sadb.BeginTransaction()
	trx.CreateAccount("test", "commit-tran@example.com")
	trx.Commit()

	account, err := sadb.FindAccountByEmail("commit-tran@example.com")
	assert.NotNil(t, account)
	assert.NoError(t, err)
}

func TestRollback(t *testing.T) {
	trx := sadb.BeginTransaction()
	trx.CreateAccount("test", "rollback-tran@example.com")
	trx.Rollback()

	account, err := sadb.FindAccountByEmail("rollback-tran@example.com")
	assert.Nil(t, account)
	assert.Error(t, err)
}
