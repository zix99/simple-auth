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
