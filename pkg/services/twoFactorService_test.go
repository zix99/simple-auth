package services

import (
	"simple-auth/pkg/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testTwoFator TwoFactorService

func init() {
	testTwoFator = NewTwoFactorService(&config.TwoFactorConfig{
		Enabled:   true,
		Drift:     1,
		Issuer:    "test-issuer",
		KeyLength: 12,
	})
}

func TestCreateSecret(t *testing.T) {
	s, err := testTwoFator.CreateSecret()
	assert.NotEmpty(t, s)
	assert.NoError(t, err)

	fs, err := testTwoFator.CreateFullSpecFromSecret(s, testAuthLocalAccount)
	assert.NotNil(t, fs)
	assert.NoError(t, err)
	assert.Equal(t, "test-issuer", fs.Issuer)
}
