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

type ConfigJWT struct {
	Issuer         string
	ExpiresMinutes int
}

type ConfigWeb struct {
	Host         string
	Secret       string
	Requirements ConfigWebRequirements
	RecaptchaV2  ConfigRecaptchaV2
	JWT          ConfigJWT
	Metadata     map[string]interface{}
}

type ConfigEmail struct {
	Enabled  bool
	Host     string
	Identity string
	Username string
	Password string
	From     string
}

type Config struct {
	Db             ConfigDatabase
	Web            ConfigWeb
	Email          ConfigEmail
	Authenticators ConfigAuthencatorSet
	Production     bool
}
