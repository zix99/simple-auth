package config

import (
	"regexp"

	"github.com/sirupsen/logrus"
)

func makeRegexBound(exp string) string {
	if exp[0] != '^' {
		exp = "^" + exp
	}
	if exp[len(exp)-1] != '$' {
		exp = exp + "$"
	}
	return exp
}

func (s *ConfigLoginSettings) cachedAllowedContinueUrls() []*regexp.Regexp {
	if s.allowedContinueUrlsRegexp == nil {
		s.allowedContinueUrlsRegexp = make([]*regexp.Regexp, 0, len(s.AllowedContinueUrls))
		for _, expStr := range s.AllowedContinueUrls {
			exp, err := regexp.Compile(makeRegexBound(expStr))
			if err == nil {
				s.allowedContinueUrlsRegexp = append(s.allowedContinueUrlsRegexp, exp)
			} else {
				logrus.Warnf("Error compiling continue URL Expression (%s): %v", s, err)
			}
		}
	}
	return s.allowedContinueUrlsRegexp
}

func matchesAny(val string, exps ...*regexp.Regexp) bool {
	for _, exp := range exps {
		if exp.MatchString(val) {
			return true
		}
	}
	return false
}

func (s *ConfigLoginSettings) ResolveContinueURL(asked string) (continueURL string) {
	continueURL = s.RouteOnLogin
	if asked != "" && matchesAny(asked, s.cachedAllowedContinueUrls()...) {
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
