package db_test

import (
	"simple-auth/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sadb db.SADB

func init() {
	sadb = db.New("sqlite3", "file::memory:?cache=shared")
	createAuthSimpleMock()
	createAuthTokenMock()
}

func TestIsAlive(t *testing.T) {
	assert.True(t, sadb.IsAlive())
}
