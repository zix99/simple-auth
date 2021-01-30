package db

import "strings"

// OAuthScope is a unified scope format with helper functions
// Normalized to lower-case
type OAuthScope []string

func NewOAuthScope(s string) OAuthScope {
	ret := strings.Split(strings.ToLower(s), " ")
	if ret[0] == "" {
		return ret[0:0]
	}
	return ret
}

func (s OAuthScope) String() string {
	if s == nil {
		return ""
	}
	return strings.Join(s, " ")
}

func (s OAuthScope) Contains(scope string) bool {
	if s == nil {
		return false
	}
	scope = strings.ToLower(scope)

	for _, ele := range s {
		if ele == scope {
			return true
		}
	}

	return false
}

func (s OAuthScope) ContainsAll(scopes ...string) bool {
	if len(scopes) == 0 {
		return true
	}
	if s == nil {
		return false
	}

	for _, ele := range scopes {
		if !s.Contains(ele) {
			return false
		}
	}

	return true
}

func (s OAuthScope) ContainsScopes(other OAuthScope) bool {
	return s.ContainsAll(other...)
}

func (s OAuthScope) Matches(other OAuthScope) bool {
	return s.ContainsScopes(other) && other.ContainsScopes(s)
}
