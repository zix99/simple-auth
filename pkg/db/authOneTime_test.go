package db_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIssueOneTime(t *testing.T) {
	account, _ := sadb.CreateAccount("onetime1@asdf.com")
	assert.NotNil(t, account)
	ott, err := sadb.CreateAccountOneTimeToken(account, 5*time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, ott)
}

func TestIssueResolveOneTime(t *testing.T) {
	account, _ := sadb.CreateAccount("onetime2@asdf.com")
	assert.NotNil(t, account)
	ott, _ := sadb.CreateAccountOneTimeToken(account, 5*time.Minute)

	ret, err := sadb.AssertOneTimeToken(ott)
	assert.NoError(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, account.UUID, ret.UUID)
}

func TestAssertEmptyOneTime(t *testing.T) {
	ret, err := sadb.AssertOneTimeToken("")
	assert.Error(t, err)
	assert.Nil(t, ret)
}

func TestAssertInvalidOneTime(t *testing.T) {
	ret, err := sadb.AssertOneTimeToken(uuid.New().String())
	assert.Error(t, err)
	assert.Nil(t, ret)
}

func TestCantDoubleConsumeOneTime(t *testing.T) {
	account, _ := sadb.CreateAccount("onetime3@asdf.com")
	assert.NotNil(t, account)
	ott, _ := sadb.CreateAccountOneTimeToken(account, 5*time.Minute)

	ret, err := sadb.AssertOneTimeToken(ott)
	assert.NoError(t, err)
	assert.NotNil(t, ret)
	assert.Equal(t, account.UUID, ret.UUID)

	ret2, err2 := sadb.AssertOneTimeToken(ott)
	assert.Error(t, err2)
	assert.Nil(t, ret2)
}

func TestCantConsumeExpiredToken(t *testing.T) {
	account, _ := sadb.CreateAccount("onetime4@asdf.com")
	assert.NotNil(t, account)
	ott, _ := sadb.CreateAccountOneTimeToken(account, 0*time.Minute)

	ret, err := sadb.AssertOneTimeToken(ott)
	assert.Error(t, err)
	assert.Nil(t, ret)
}
