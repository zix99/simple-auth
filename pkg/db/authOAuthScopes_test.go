package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseScopes(t *testing.T) {
	s := NewOAuthScope("")
	assert.Len(t, s, 0)

	s = NewOAuthScope("a")
	assert.Len(t, s, 1)
	assert.Contains(t, s, "a")

	s = NewOAuthScope("a b cde")
	assert.Len(t, s, 3)
	assert.Contains(t, s, "a")
	assert.Contains(t, s, "b")
	assert.Contains(t, s, "cde")
}

func TestStringifyScopes(t *testing.T) {
	s := NewOAuthScope("")
	assert.Equal(t, s.String(), "")

	s = NewOAuthScope("a")
	assert.Equal(t, s.String(), "a")

	s = NewOAuthScope("a b cef")
	assert.Equal(t, s.String(), "a b cef")
}

func TestContainsScope(t *testing.T) {
	assert.False(t, NewOAuthScope("").Contains("a"))
	assert.False(t, NewOAuthScope("b").Contains("a"))
	assert.False(t, NewOAuthScope("b c").Contains("a"))

	assert.True(t, NewOAuthScope("a").Contains("a"))
	assert.True(t, NewOAuthScope("a b cef").Contains("a"))
	assert.True(t, NewOAuthScope("a b cef").Contains("cef"))
}

func TestContainsAllScopes(t *testing.T) {
	assert.False(t, NewOAuthScope("").ContainsAll("a"))
	assert.False(t, NewOAuthScope("b").ContainsAll("a", "b"))
	assert.False(t, NewOAuthScope("b c").ContainsAll("a"))

	assert.True(t, NewOAuthScope("a").ContainsAll())
	assert.True(t, NewOAuthScope("a").ContainsAll("a"))
	assert.True(t, NewOAuthScope("a b cef").ContainsAll("a", "b"))
	assert.True(t, NewOAuthScope("a b cef").ContainsAll("cef", "a"))
}

func TestMatchesScope(t *testing.T) {
	assert.False(t, NewOAuthScope("a").Matches([]string{"a", "b"}))
	assert.False(t, NewOAuthScope("a").Matches([]string{"b", "b"}))
	assert.False(t, NewOAuthScope("a q").Matches([]string{"b", "b"}))
	assert.False(t, NewOAuthScope("a").Matches([]string{}))

	assert.True(t, NewOAuthScope("a").Matches([]string{"a"}))
	assert.True(t, NewOAuthScope("a").Matches([]string{"a", "a"}))
	assert.True(t, NewOAuthScope("a b").Matches([]string{"b", "a"}))
	assert.True(t, NewOAuthScope("a b").Matches([]string{"b", "a", "b"}))
	assert.True(t, NewOAuthScope("b b").Matches([]string{"b"}))
}
