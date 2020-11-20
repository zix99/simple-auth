package db_test

import (
	"simple-auth/pkg/db"
	"simple-auth/pkg/lib/totp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var authSimpleAccount *db.Account

const authSimpleEmail = "auth-simple-test@asdf.com"
const authSimpleUsername = "auth"
const authSimplePassword = "blablabla22"

func createAuthSimpleMock() {
	authSimpleAccount, _ = sadb.CreateAccount(authSimpleEmail)
	sadb.CreateAuthLocal(authSimpleAccount, authSimpleUsername, authSimplePassword)
}

func TestCreateAuthSimple(t *testing.T) {
	assert.NotNil(t, authSimpleAccount)
}

func TestCreateInvalidAccount(t *testing.T) {
	_, err := sadb.CreateAuthLocal(authSimpleAccount, "", "")
	assert.Error(t, err)
}

func TestSimpleAuthLookup(t *testing.T) {
	authLocal, err := sadb.FindAuthLocal(authSimpleAccount)
	assert.Equal(t, authSimpleUsername, authLocal.Username())
	assert.NoError(t, err)

	assert.True(t, authLocal.VerifyPassword(authSimplePassword))
	assert.False(t, authLocal.VerifyPassword("bad-pass"))
	assert.False(t, authLocal.HasTOTP())
}

func TestSimpleAuthLookupAccountNil(t *testing.T) {
	authLocal, err := sadb.FindAuthLocal(nil)
	assert.Empty(t, authLocal)
	assert.Error(t, err)
}

func TestSimpleAuthLookupAccountNoLink(t *testing.T) {
	account, _ := sadb.CreateAccount("no-simpleauth-account@asdf.com")
	authLocal, err := sadb.FindAuthLocal(account)
	assert.Empty(t, authLocal)
	assert.Error(t, err)
}

func TestFindByUsername(t *testing.T) {
	authLocal, err := sadb.FindAuthLocalByUsername(authSimpleUsername)
	assert.NoError(t, err)

	account := authLocal.Account()
	assert.Equal(t, authSimpleEmail, account.Email)
	assert.Equal(t, authSimpleAccount.UUID, account.UUID)
}

func TestFindByUsernameFail(t *testing.T) {
	authLocal, err := sadb.FindAuthLocalByUsername("not-exist")
	assert.Error(t, err)
	assert.Nil(t, authLocal)
}

func TestCreateDupeUsername(t *testing.T) {
	authLocal, err := sadb.CreateAuthLocal(authSimpleAccount, authSimpleUsername, authSimplePassword)
	assert.Error(t, err)
	assert.Nil(t, authLocal)
}

func TestUpdatePassword(t *testing.T) {
	const changePassUsername = "change-pass-uname"
	account, _ := sadb.CreateAccount("change-pass@asdf.com")
	assert.NotNil(t, account)
	authLocal, _ := sadb.CreateAuthLocal(account, changePassUsername, authSimplePassword)

	sadb.UpdateAuthLocalPassword(authLocal, "new-password")

	authLocal, _ = sadb.FindAuthLocal(account)

	passVerify := authLocal.VerifyPassword("new-password")
	assert.True(t, passVerify)
}

func TestUpdateTOTP(t *testing.T) {
	tfa, _ := totp.NewTOTP(12, "test", "test")

	account, _ := sadb.CreateAccount("totp-test@asdf.com")
	authLocal, _ := sadb.CreateAuthLocal(account, "test-totp", "test-totp")

	tStr := tfa.String()
	sadb.UpdateAuthLocalTOTP(authLocal, &tStr)

	authLocalUpdated, _ := sadb.FindAuthLocal(account)
	assert.True(t, authLocalUpdated.HasTOTP())
	assert.True(t, authLocal.VerifyTOTP(tfa.GetTOTP(), 1))
}