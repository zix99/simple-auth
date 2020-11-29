package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailValidation(t *testing.T) {
	assert.NoError(t, validateEmail("a@b.co"))
	assert.NoError(t, validateEmail("a+b@c.com"))

	assert.Error(t, validateEmail("abcasdf"))
	assert.Error(t, validateEmail("asdf@asdf"))
}
