package services

import (
	"simple-auth/pkg/config"
	"simple-auth/pkg/db"
	"simple-auth/pkg/email"
	"simple-auth/pkg/email/engine"
	"simple-auth/pkg/lib/totp"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testLocalLogin LocalLoginService
var testLocalLoginAccount *db.Account
var testAuthLocalAccount *db.AuthLocal
var (
	testLocalEmail    = "test@locallogin.com"
	testLocalUsername = "test-user"
	testLocalPassword = "test-pass"
)

func init() {
	sadb := getDB()
	mockEmailService := email.New(logrus.StandardLogger(), engine.NewMockEngine(nil), "test@example.com")

	testLocalLogin = NewLocalLoginService(sadb, mockEmailService, &config.ConfigMetadata{
		Company: "test-corp",
	}, &config.TwoFactorConfig{
		Enabled:   true,
		Drift:     1,
		Issuer:    "test",
		KeyLength: 12,
	}, &config.ConfigWebRequirements{
		EmailValidationRequired: true,
	}, "http://example.com")

	testLocalLoginAccount, _ = sadb.CreateAccount("test", testLocalEmail)
	testAuthLocalAccount, _ = sadb.CreateAuthLocal(testLocalLoginAccount, testLocalUsername, testLocalPassword)
}

func TestAssertLoginSuccess(t *testing.T) {
	authLocal, err := testLocalLogin.AssertLogin(testLocalUsername, testLocalPassword, nil)
	assert.NotNil(t, authLocal)
	assert.NoError(t, err)
	assert.Equal(t, testLocalEmail, authLocal.Account().Email)
	assert.Equal(t, testLocalLoginAccount.UUID, authLocal.Account().UUID)
}

func TestAssertLoginFail(t *testing.T) {
	authLocal, err := testLocalLogin.AssertLogin(testLocalUsername, "made-up", nil)
	assert.Nil(t, authLocal)
	assert.Error(t, err)
}

func TestFindLoginByAccount(t *testing.T) {
	authLocal, err := testLocalLogin.FindAuthLocal(testLocalLoginAccount.UUID)
	assert.NotNil(t, authLocal)
	assert.NoError(t, err)
}

func TestFindLoginByAccountMissing(t *testing.T) {
	authLocal, err := testLocalLogin.FindAuthLocal("made-up")
	assert.Nil(t, authLocal)
	assert.Error(t, err)
}

func TestSimpleAuthTOTP(t *testing.T) {
	sadb := getDB()
	account, _ := sadb.CreateAccount("test", "totp-account@asdf.com")
	assert.NotNil(t, account)

	// Simple setup
	authLocal, err := sadb.CreateAuthLocal(account, "totp", "totp-pass")
	assert.NoError(t, err)
	{
		authLocal, err := testLocalLogin.AssertLogin("totp", "totp-pass", nil)
		assert.NotNil(t, authLocal)
		assert.NoError(t, err)
	}
	{
		totpCode := "123"
		authLocal, err := testLocalLogin.AssertLogin("totp", "totp-pass", &totpCode)
		assert.NotNil(t, authLocal)
		assert.NoError(t, err)
	}

	// Set up totp
	otp, err := totp.NewTOTP(8, "test", "totp")
	assert.NoError(t, err)

	otpURL := otp.String()
	assert.NoError(t, sadb.UpdateAuthLocalTOTP(authLocal, &otpURL))

	{
		authLocal, err := testLocalLogin.AssertLogin("totp", "totp-pass", nil)
		assert.Error(t, err)
		assert.Nil(t, authLocal)
	}
	{
		totpCode := "abcdef" // Will never be letters
		authLocal, err := testLocalLogin.AssertLogin("totp", "totp-pass", &totpCode)
		assert.Error(t, err)
		assert.Nil(t, authLocal)
	}
	{
		totpCode := otp.GetTOTP()
		authLocal, err := testLocalLogin.AssertLogin("totp", "totp-pass", &totpCode)
		assert.NotNil(t, authLocal)
		assert.NoError(t, err)
	}
}

func TestAuthLocalModifyPassword(t *testing.T) {
	sadb := getDB()
	account, _ := sadb.CreateAccount("test", "passwd-account@asdf.com")
	assert.NotNil(t, account)

	// Simple setup
	authLocal, err := sadb.CreateAuthLocal(account, "passchange", "passchange-test")
	assert.NoError(t, err)
	assert.NotNil(t, authLocal)

	assert.Error(t, testLocalLogin.UpdatePassword(authLocal, "passchange-WRONG", "bla"))
	assert.NoError(t, testLocalLogin.UpdatePassword(authLocal, "passchange-test", "bla"))
}
