package config

func stringInList(val string, lst []string) bool {
	for _, item := range lst {
		if item == val {
			return true
		}
	}
	return false
}

func (s *ConfigLoginSettings) ResolveContinueURL(asked string) (continueURL string) {
	continueURL = s.RouteOnLogin
	if asked != "" && stringInList(asked, s.AllowedContinueUrls) {
		continueURL = asked
	}
	return
}

func (s *ConfigWeb) GetBaseURL() string {
	if s.BaseURL != "" {
		return s.BaseURL
	}
	if s.TLS.Enabled {
		return "https://" + s.Host
	}
	return "http://" + s.Host
}
