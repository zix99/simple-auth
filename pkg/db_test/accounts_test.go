package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	account, err := sadb.CreateAccount("me@asdf.com")
	assert.NoError(t, err)
	assert.Equal(t, "me@asdf.com", account.Email)
	assert.True(t, account.Active)
	assert.NotEmpty(t, account.UUID)
}

func TestCreateAccountDupeEmails(t *testing.T) {
	sadb.CreateAccount("dupe@asdf.com")
	account, err := sadb.CreateAccount("dupe@asdf.com")
	assert.Error(t, err)
	assert.Nil(t, account)
}
