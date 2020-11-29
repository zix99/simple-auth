package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	account, err := sadb.CreateAccount("test", "me@asdf.com")
	assert.NoError(t, err)
	assert.Equal(t, "me@asdf.com", account.Email)
	assert.True(t, account.Active)
	assert.NotEmpty(t, account.UUID)
}

func TestCreateAccountDupeEmails(t *testing.T) {
	sadb.CreateAccount("test", "dupe@asdf.com")
	account, err := sadb.CreateAccount("test", "dupe@asdf.com")
	assert.Error(t, err)
	assert.Nil(t, account)
}

func TestFindAccount(t *testing.T) {
	newAccount, _ := sadb.CreateAccount("test", "findme@asdf.com")
	account, err := sadb.FindAccount(newAccount.UUID)
	assert.NoError(t, err)
	assert.Equal(t, newAccount.ID, account.ID)
	assert.Equal(t, "findme@asdf.com", account.Email)
}

func TestFindAccountFail(t *testing.T) {
	account, err := sadb.FindAccount("not-exist")
	assert.Nil(t, account)
	assert.Error(t, err)
}
