package services

import (
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	"simple-auth/pkg/email/engine"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testOAuthService AuthOAuthService
var testOAuthAccount *db.Account

func init() {
	sadb := getDB()

	ctx := appcontext.NewContainer()
	ctx.Use(appcontext.WithSADB(sadb))

	localLoginService := NewLocalLoginService(
		email.New(engine.NewMockEngine(nil), "test@example.com"),
		&config.ConfigMetadata{},
		&config.ConfigLocalProvider{},
		"http://example.com",
	)

	testOAuthService = NewAuthOAuthService("test-client", &config.ConfigOAuth2Client{
		Secret:      "test-secret",
		RedirectURI: "http://example.com/redirect",
		Scopes:      []string{"email", "user"},
		OIDC: &config.OAuth2OIDCConfig{
			SigningMethod: "HS256",
			SigningKey:    "abcdef721yu4uih",
		},
	}, &config.ConfigOAuth2Settings{
		CodeExpiresSeconds:  config.IntPtr(10),
		TokenExpiresSeconds: config.IntPtr(20),
		CodeLength:          config.IntPtr(6),
		AllowCredentials:    config.TruePtr,
		IssueRefreshToken:   config.TruePtr,
		AllowAutoGrant:      config.TruePtr,
		ReuseToken:          config.FalsePtr,
		Issuer:              config.StrPtr("simple-auth"),
	}, localLoginService).WithContext(ctx)

	testOAuthAccount, _ = sadb.CreateAccount("test-oauth", "test-oauth@example.com")
	sadb.CreateAuthLocal(testOAuthAccount, "oauth-user", "oauth-pass")
}

func TestCodeGen(t *testing.T) {
	code, _ := genAccessCode(5)
	assert.Len(t, code, 5)
}

func TestCreateAccessCode(t *testing.T) {
	code, err := testOAuthService.CreateAccessCode(testOAuthAccount, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, code)
	assert.Len(t, code, 6)
}

func TestOAuthTradeAccessCode(t *testing.T) {
	code, _ := testOAuthService.CreateAccessCode(testOAuthAccount, nil)
	assert.NotEmpty(t, code)

	{
		token, err := testOAuthService.TradeCodeForToken("invalid", code)
		assert.Empty(t, token.AccessToken)
		assert.Empty(t, token.RefreshToken)
		assert.Error(t, err)
	}

	{
		token, err := testOAuthService.TradeCodeForToken("test-secret", "invalid")
		assert.Empty(t, token.AccessToken)
		assert.Empty(t, token.RefreshToken)
		assert.Error(t, err)
	}

	token, err := testOAuthService.TradeCodeForToken("test-secret", code)
	assert.NoError(t, err)
	assert.NotEmpty(t, token.AccessToken)
	assert.NotEmpty(t, token.RefreshToken)
	assert.NotEmpty(t, token.IDToken)
	assert.Greater(t, token.Expires, 0)
}

func TestTradeRefreshForToken(t *testing.T) {
	code, _ := testOAuthService.CreateAccessCode(testOAuthAccount, nil)
	token, _ := testOAuthService.TradeCodeForToken("test-secret", code)

	refreshed, err := testOAuthService.TradeRefreshTokenForAccessToken("test-secret", token.RefreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshed.AccessToken)
	assert.Empty(t, refreshed.RefreshToken)
}

func TestAutoRevokeTokenOnNew(t *testing.T) {
	code1, _ := testOAuthService.CreateAccessCode(testOAuthAccount, nil)
	code2, _ := testOAuthService.CreateAccessCode(testOAuthAccount, nil)
	assert.NotEmpty(t, code2)

	token1, err := testOAuthService.TradeCodeForToken("test-secret", code1)
	assert.NoError(t, err)
	assert.NotEmpty(t, token1.AccessToken)

	token2, err := testOAuthService.TradeCodeForToken("test-secret", code2)
	assert.Error(t, err)
	assert.Empty(t, token2.AccessToken)
}

func TestOAuthScopes(t *testing.T) {
	scope := db.NewOAuthScope("email user")
	code, _ := testOAuthService.CreateAccessCode(testOAuthAccount, scope)

	token, err := testOAuthService.TradeCodeForToken("test-secret", code)
	assert.NoError(t, err)
	assert.True(t, token.Scope.Matches(scope))
}

func TestOAuthScopesFail(t *testing.T) {
	scope := db.NewOAuthScope("email user admin")
	code, err := testOAuthService.CreateAccessCode(testOAuthAccount, scope)

	assert.Error(t, err)
	assert.Empty(t, code)
}

func TestOAuthTradeForCredentials(t *testing.T) {
	{
		token, err := testOAuthService.TradeCredentialsForToken("test-secret", "oauth-user", "bad-pass", nil, nil)
		assert.Error(t, err)
		assert.Empty(t, token.AccessToken)
		assert.Empty(t, token.RefreshToken)
	}

	{
		token, err := testOAuthService.TradeCredentialsForToken("bad-secret", "oauth-user", "oauth-pass", nil, nil)
		assert.Error(t, err)
		assert.Empty(t, token.AccessToken)
		assert.Empty(t, token.RefreshToken)
	}

	{
		token, err := testOAuthService.TradeCredentialsForToken("test-secret", "oauth-user", "oauth-pass", nil, db.NewOAuthScope("email"))
		assert.NoError(t, err)
		assert.NotEmpty(t, token.AccessToken)
		assert.NotEmpty(t, token.RefreshToken)
		assert.Greater(t, token.Expires, 0)
		assert.NotEmpty(t, token.IDToken)
	}

	{
		token, err := testOAuthService.TradeCredentialsForToken("test-secret", "oauth-user", "oauth-pass", nil, db.NewOAuthScope("email bad-scope"))
		assert.Error(t, err)
		assert.Empty(t, token.AccessToken)
		assert.Empty(t, token.RefreshToken)
	}
}

func TestFindToken(t *testing.T) {
	code, _ := testOAuthService.CreateAccessCode(testOAuthAccount, nil)

	{
		found, err := testOAuthService.FindExistingToken(testOAuthAccount, db.OAuthTypeCode, nil)
		assert.NoError(t, err)
		assert.Equal(t, code, found.AccessToken)
	}

	token, _ := testOAuthService.TradeCodeForToken("test-secret", code)
	{
		found, err := testOAuthService.FindExistingToken(testOAuthAccount, db.OAuthTypeCode, nil)
		assert.Error(t, err)
		assert.Empty(t, found.AccessToken)
	}
	{
		found, err := testOAuthService.FindExistingToken(testOAuthAccount, db.OAuthTypeAccessToken, nil)
		assert.NoError(t, err)
		assert.Equal(t, token.AccessToken, found.AccessToken)
	}
}

func TestVerifyRedirectURL(t *testing.T) {
	assert.True(t, testOAuthService.ValidateRedirectURI("http://example.com/redirect"))
	assert.False(t, testOAuthService.ValidateRedirectURI("http://example.com/redirect2"))
}

func TestValidateScopes(t *testing.T) {
	assert.True(t, testOAuthService.ValidateScopes(db.NewOAuthScope("email")))
	assert.True(t, testOAuthService.ValidateScopes(db.NewOAuthScope("user")))
	assert.True(t, testOAuthService.ValidateScopes(db.NewOAuthScope("user email")))

	assert.False(t, testOAuthService.ValidateScopes(db.NewOAuthScope("admin")))
	assert.False(t, testOAuthService.ValidateScopes(db.NewOAuthScope("admin email")))
}
