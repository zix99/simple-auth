package config

type ConfigDatabase struct {
	Driver string
	URL    string
}

type ConfigAuthencatorSet struct {
	Token struct {
		Enabled                    bool
		SessionExpiresMinutes      int
		VerificationExpiresSeconds int
	}
	Simple struct {
		Enabled bool
	}
}

type ConfigWebRequirements struct {
	PasswordMinLength int
	PasswordMaxLength int
	UsernameMinLength int
	UsernameMaxLength int
}

type ConfigRecaptchaV2 struct {
	Enabled bool
	SiteKey string
	Secret  string
	Theme   string
}

type ConfigWeb struct {
	Host         string
	Requirements ConfigWebRequirements
	RecaptchaV2  ConfigRecaptchaV2
	Metadata     map[string]interface{}
}

type Config struct {
	Db             ConfigDatabase
	Web            ConfigWeb
	Authenticators ConfigAuthencatorSet
	Production     bool
}
