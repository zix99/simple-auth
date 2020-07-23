package db_test

import (
	"simple-auth/pkg/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var authTokenAccount *db.Account

const authTokenEmail = "auth-token-test@asdf.com"
const authTokenUsername = "authtoken"
const authTokenPassword = "howdy"

func createAuthTokenMock() {
	authTokenAccount, _ = sadb.CreateAccount(authTokenEmail)
	sadb.CreateAccountAuthSimple(authTokenAccount, authTokenUsername, authTokenPassword)
}

func TestTokenHappyPath(t *testing.T) {
	session, err := sadb.AssertCreateSessionToken(authTokenUsername, authTokenPassword, 10*time.Second)
	assert.NoError(t, err)

	verification, err := sadb.CreateVerificationToken(authTokenUsername, session)
	assert.NoError(t, err)

	account, err := sadb.AssertVerificationToken(authTokenUsername, verification)
	assert.NoError(t, err)
	assert.Equal(t, authTokenAccount.UUID, account.UUID)

	// Should fail 2nd time
	account, err = sadb.AssertVerificationToken(authTokenUsername, verification)
	assert.Error(t, err)
	assert.Nil(t, account)
}

func TestTokenExpiration(t *testing.T) {
	session, err := sadb.AssertCreateSessionToken(authTokenUsername, authTokenPassword, 0*time.Second)
	assert.NoError(t, err)

	verification, err := sadb.CreateVerificationToken(authTokenUsername, session)
	assert.Error(t, err)
	assert.Empty(t, verification)
}

func TestNoDoubleSession(t *testing.T) {
	session1, err := sadb.AssertCreateSessionToken(authTokenUsername, authTokenPassword, 10*time.Second)
	assert.NoError(t, err)

	session2, err := sadb.AssertCreateSessionToken(authTokenUsername, authTokenPassword, 10*time.Second)
	assert.NoError(t, err)

	// Use first to create/verify
	verify1, err := sadb.CreateVerificationToken(authTokenUsername, session1)
	assert.Empty(t, verify1)
	assert.Error(t, err)

	account, err := sadb.AssertVerificationToken(authTokenUsername, verify1)
	assert.Nil(t, account)
	assert.Error(t, err)

	// Use 2nd to create/verify
	verify2, err := sadb.CreateVerificationToken(authTokenUsername, session2)
	assert.NotEmpty(t, verify2)
	assert.NoError(t, err)

	account, err = sadb.AssertVerificationToken(authTokenUsername, verify2)
	assert.NoError(t, err)
	assert.NotNil(t, account)
}

func TestInvalidateSession(t *testing.T) {
	session, err := sadb.AssertCreateSessionToken(authTokenUsername, authTokenPassword, 10*time.Second)
	assert.NoError(t, err)

	assert.NoError(t, sadb.InvalidateSession(session))

	verify, err := sadb.CreateVerificationToken(authTokenUsername, session)
	assert.Empty(t, verify)
	assert.Error(t, err)
}
