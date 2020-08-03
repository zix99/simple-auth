package db_test

import (
	"simple-auth/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

var oidcAccount *db.Account

const oidcEmail = "oidc-test@asdf.com"

func createOIDCMock() {
	oidcAccount, _ = sadb.CreateAccount(oidcEmail)
}

func TestCreateOIDCOnAccount(t *testing.T) {
	err := sadb.CreateOIDCForAccount(oidcAccount, "test", "abcd")
	assert.NoError(t, err)

	getAccount, err := sadb.FindAccountForOIDC("test", "abcd")
	assert.NoError(t, err)
	assert.Equal(t, oidcAccount.UUID, getAccount.UUID)
}

func TestCreateTwoProviders(t *testing.T) {
	assert.NoError(t, sadb.CreateOIDCForAccount(oidcAccount, "p1", "abcd"))
	assert.NoError(t, sadb.CreateOIDCForAccount(oidcAccount, "p2", "abcd"))
}

func TestCreateOIDCDupe(t *testing.T) {
	err1 := sadb.CreateOIDCForAccount(oidcAccount, "dupe", "quack")
	assert.NoError(t, err1)
	err2 := sadb.CreateOIDCForAccount(oidcAccount, "dupe", "quack")
	assert.Error(t, err2)
}

func TestMissingOIDCAccount(t *testing.T) {
	find, err := sadb.FindAccountForOIDC("laksdjf", "nope")
	assert.Nil(t, find)
	assert.Error(t, err)
}
