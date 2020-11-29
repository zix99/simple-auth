package db_test

import (
	"simple-auth/pkg/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenStipulation(t *testing.T) {
	account, _ := sadb.CreateAccount("test", "stip@asdf.com")

	assert.False(t, sadb.AccountHasUnsatisfiedStipulations(account))

	err := sadb.AddStipulation(account, &db.TokenStipulation{
		Code: "abc",
	})
	assert.NoError(t, err)

	assert.True(t, sadb.AccountHasUnsatisfiedStipulations(account))

	{
		err := sadb.SatisfyStipulation(account, &db.TokenStipulation{
			Code: "qrf",
		})
		assert.Error(t, err)
	}

	{
		err := sadb.SatisfyStipulation(account, &db.TokenStipulation{
			Code: "abc",
		})
		assert.NoError(t, err)
		assert.False(t, sadb.AccountHasUnsatisfiedStipulations(account))
	}
}

func TestUpdateAllStipulations(t *testing.T) {
	account, _ := sadb.CreateAccount("test", "stip2@asdf.com")

	assert.NoError(t, sadb.AddStipulation(account, db.NewTokenStipulation()))
	assert.NoError(t, sadb.AddStipulation(account, &db.ManualStipulation{}))

	assert.True(t, sadb.AccountHasUnsatisfiedStipulations(account))

	assert.NoError(t, sadb.ForceSatisfyStipulations(account))

	assert.False(t, sadb.AccountHasUnsatisfiedStipulations(account))
}
