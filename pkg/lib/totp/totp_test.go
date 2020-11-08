package totp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestURLEncoding(t *testing.T) {
	otp, err := NewTOTP(10, "coolco", "george")
	assert.NoError(t, err)
	assert.NotNil(t, otp)

	url := otp.String()
	otp2, err2 := ParseTOTP(url)
	assert.NoError(t, err2)
	assert.Equal(t, otp.Issuer, otp2.Issuer)
	assert.Equal(t, otp.Subject, otp2.Subject)
	assert.Equal(t, otp.Secret(), otp2.Secret())
	assert.Equal(t, otp.GetHOTP(1), otp2.GetHOTP(1))
}

func TestSecretTransmission(t *testing.T) {
	otp, err := NewTOTP(10, "coolco", "george")
	assert.NoError(t, err)

	secret := otp.Secret()
	assert.NotEmpty(t, secret)

	otp2, err := FromSecret(secret, "", "")
	assert.NoError(t, err)
	assert.Equal(t, otp.Secret(), otp2.Secret())
	assert.Equal(t, otp.GetHOTP(2), otp2.GetHOTP(2))
}

func TestTOTPValidation(t *testing.T) {
	otp, err := NewTOTP(10, "coolco", "george")
	assert.NoError(t, err)

	code := otp.GetTOTP()
	assert.Len(t, code, 6)

	// Test validate
	assert.True(t, otp.Validate(code, 1))
}

func TestTOTPValidationDrift(t *testing.T) {
	otp, err := NewTOTP(10, "coolco", "george")
	assert.NoError(t, err)

	interval := time.Now().Unix()/30 - 1
	code := otp.GetHOTP(interval)
	assert.Len(t, code, 6)

	// Test validate
	assert.True(t, otp.Validate(code, 2))
	assert.False(t, otp.Validate(code, 0))
}
