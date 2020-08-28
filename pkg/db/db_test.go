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
