package db_test

import (
	"simple-auth/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

var authSimpleAccount *db.Account

const authSimpleEmail = "auth-simple-test@asdf.com"
const authSimpleUsername = "auth"
const authSimplePassword = "blablabla22"

func createAuthSimpleMock() {
	authSimpleAccount, _ = sadb.CreateAccount(authSimpleEmail)
	sadb.CreateAccountAuthSimple(authSimpleAccount, authSimpleUsername, authSimplePassword)
}

func TestCreateAuthSimple(t *testing.T) {
	assert.NotNil(t, authSimpleAccount)
}

func TestCreateInvalidAccount(t *testing.T) {
	err := sadb.CreateAccountAuthSimple(authSimpleAccount, "", "")
	assert.Error(t, err)
}

func TestSimpleAuthLookupUsername(t *testing.T) {
	username, err := sadb.FindSimpleAuthUsername(authSimpleAccount)
	assert.Equal(t, authSimpleUsername, username)
	assert.NoError(t, err)
}

func TestSimpleAuthLookupUsernameAccountNil(t *testing.T) {
	username, err := sadb.FindSimpleAuthUsername(nil)
	assert.Empty(t, username)
	assert.Error(t, err)
}

func TestSimpleAuthLookupUsernameAccountNoLink(t *testing.T) {
	account, _ := sadb.CreateAccount("no-simpleauth-account@asdf.com")
	username, err := sadb.FindSimpleAuthUsername(account)
	assert.Empty(t, username)
	assert.Error(t, err)
}

func TestFindByUsername(t *testing.T) {
	account, err := sadb.FindAccountForSimpleAuth(authSimpleUsername)
	assert.NoError(t, err)
	assert.Equal(t, authSimpleEmail, account.Email)
	assert.Equal(t, authSimpleAccount.UUID, account.UUID)
}

func TestFindByUsernameFail(t *testing.T) {
	account, err := sadb.FindAccountForSimpleAuth("not-exist")
	assert.Error(t, err)
	assert.Nil(t, account)
}

func TestAssertLoginSuccess(t *testing.T) {
	account, err := sadb.AssertSimpleAuth(authSimpleUsername, authSimplePassword)
	assert.NotNil(t, account)
	assert.NoError(t, err)
	assert.Equal(t, authSimpleEmail, account.Email)
	assert.Equal(t, authSimpleAccount.UUID, account.UUID)
}

func TestAssertLoginFail(t *testing.T) {
	account, err := sadb.AssertSimpleAuth(authSimpleUsername, "made-up")
	assert.Nil(t, account)
	assert.Error(t, err)
}

func TestCreateDupeUsername(t *testing.T) {
	err := sadb.CreateAccountAuthSimple(authSimpleAccount, authSimpleUsername, authSimplePassword)
	assert.Error(t, err)
}

func TestUpdatePassword(t *testing.T) {
	const changePassUsername = "change-pass-uname"
	account, _ := sadb.CreateAccount("change-pass@asdf.com")
	assert.NotNil(t, account)
	sadb.CreateAccountAuthSimple(account, changePassUsername, authSimplePassword)

	sadb.UpdatePasswordForUsername(changePassUsername, "new-password")

	verifiedAccount, err := sadb.AssertSimpleAuth(changePassUsername, "new-password")
	assert.NoError(t, err)
	assert.Equal(t, account.UUID, verifiedAccount.UUID)
}
