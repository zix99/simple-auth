package db_test

import (
	"simple-auth/pkg/db"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var oauthTestAccount *db.Account

const oauthTestEmail = "oauth-test@asdf.com"
const oauthTestUsername = "oauth"
const oauthTestPassowrd = "oauthtest3"

const oauthTestClientID = "test-client"

func createOAuthMock() {
	oauthTestAccount, _ = sadb.CreateAccount("test-oauth", oauthTestEmail)
	sadb.CreateAuthLocal(oauthTestAccount, oauthTestUsername, oauthTestPassowrd)
}

func TestCreateOAuthToken(t *testing.T) {
	err := sadb.CreateOAuthToken(oauthTestAccount, oauthTestClientID, db.OAuthTypeAccessToken, uuid.New().String(), nil, 1*time.Hour)
	assert.NoError(t, err)
}

func TestGetToken(t *testing.T) {
	token := uuid.New().String()
	sadb.CreateOAuthToken(oauthTestAccount, oauthTestClientID, db.OAuthTypeAccessToken, token, nil, 1*time.Hour)

	got, err := sadb.GetValidOAuthToken(token)
	assert.NoError(t, err)
	assert.Equal(t, token, got.Token)
	assert.Equal(t, db.OAuthTypeAccessToken, got.Type)
	assert.Equal(t, oauthTestAccount.UUID, got.Account.UUID)
}

func TestRevokeGetToken(t *testing.T) {
	token := uuid.New().String()
	sadb.CreateOAuthToken(oauthTestAccount, oauthTestClientID, db.OAuthTypeAccessToken, token, nil, 1*time.Hour)
	sadb.InvalidateAllOAuth(oauthTestClientID, oauthTestAccount)

	got, err := sadb.GetValidOAuthToken(token)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func TestExpiredGetToken(t *testing.T) {
	token := uuid.New().String()
	sadb.CreateOAuthToken(oauthTestAccount, oauthTestClientID, db.OAuthTypeAccessToken, token, nil, -1*time.Hour)

	got, err := sadb.GetValidOAuthToken(token)
	assert.NoError(t, err)
	assert.Nil(t, got)
}
