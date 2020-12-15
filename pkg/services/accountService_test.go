package services

import (
	"simple-auth/pkg/appcontext"
	"simple-auth/pkg/config"
	"simple-auth/pkg/email"
	"simple-auth/pkg/email/engine"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmailValidation(t *testing.T) {
	assert.NoError(t, validateEmail("a@b.co"))
	assert.NoError(t, validateEmail("a+b@c.com"))

	assert.Error(t, validateEmail("abcasdf"))
	assert.Error(t, validateEmail("asdf@asdf"))
}

func TestCreateAccount(t *testing.T) {
	mockEngine := engine.NewMockEngine(nil)
	emailService := email.New(mockEngine, "test@test.comm")
	ctx := appcontext.NewContainer()
	ctx.Use(appcontext.WithSADB(getDB()))
	acctSrv := NewAccountService(&config.ConfigMetadata{}, &config.ConfigWeb{}, emailService).WithContext(ctx)

	acct, err := acctSrv.CreateAccount("test create account", "create-acct-service@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, acct)
	assert.Eventually(t, func() bool {
		return mockEngine.SendCount() == 1
	}, 2*time.Second, 100*time.Millisecond)
}
