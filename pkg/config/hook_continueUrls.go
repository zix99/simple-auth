package config

import (
	"regexp"

	"github.com/sirupsen/logrus"
)

func hookAddContinueUrls(config *Config) {
	config.Web.Login.Settings._allowedContinueUrlsRegexp = buildAllowedContinueUrlsRegexp(config.Web.Login.Settings.AllowedContinueUrls...)
}

func makeRegexBound(exp string) string {
	if exp[0] != '^' {
		exp = "^" + exp
	}
	if exp[len(exp)-1] != '$' {
		exp = exp + "$"
	}
	return exp
}

func buildAllowedContinueUrlsRegexp(urls ...string) []*regexp.Regexp {
	out := make([]*regexp.Regexp, 0, len(urls))
	for _, expStr := range urls {
		exp, err := regexp.Compile(makeRegexBound(expStr))
		if err == nil {
			out = append(out, exp)
		} else {
			logrus.Fatalf("Error compiling continue URL Expression (%s): %v", expStr, err)
		}
	}
	return out
}

func matchesAny(val string, exps ...*regexp.Regexp) bool {
	for _, exp := range exps {
		if exp.MatchString(val) {
			return true
		}
	}
	return false
}

var allowedInternalUrls = buildAllowedContinueUrlsRegexp(
	"/#/oauth2.*",
)

func (s *ConfigLoginSettings) ResolveContinueURL(asked string) (continueURL string) {
	continueURL = s.RouteOnLogin
	if asked != "" && (matchesAny(asked, allowedInternalUrls...) || matchesAny(asked, s._allowedContinueUrlsRegexp...)) {
		continueURL = asked
	}
	return
}
