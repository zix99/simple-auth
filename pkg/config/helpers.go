package config

func (s *ConfigWeb) GetBaseURL() string {
	if s.BaseURL != "" {
		return s.BaseURL
	}
	if s.TLS.Enabled {
		return "https://" + s.Host
	}
	return "http://" + s.Host
}

func (s *ConfigOAuth2Settings) Coalesce(other *ConfigOAuth2Settings) *ConfigOAuth2Settings {
	return &ConfigOAuth2Settings{
		CoalesceBool(s.IssueRefreshToken, other.IssueRefreshToken),
		CoalesceInt(s.CodeExpiresSeconds, other.CodeExpiresSeconds),
		CoalesceInt(s.TokenExpiresSeconds, other.TokenExpiresSeconds),
		CoalesceInt(s.CodeLength, other.CodeLength),
		CoalesceBool(s.AllowAutoGrant, other.AllowAutoGrant),
		CoalesceBool(s.AllowCredentials, other.AllowCredentials),
		CoalesceBool(s.ReuseToken, other.ReuseToken),
		CoalesceString(s.Issuer, other.Issuer),
	}
}
