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
	SiteKey string // Site key, as provided by google
	Secret  string // Secret for google v2 api (server side call)
	Theme   string // Themes (dark or light)
}

type ConfigJWT struct {
	Secret         string
	Issuer         string
	ExpiresMinutes int
}

type ConfigLogin struct {
	// SameDomain authentication uses a cookie set to a domain (and presumably shared with your site).  Easiest to implement in a full-trust environment
	Cookie struct {
		Enabled        bool
		Name           string // Name of the cookie
		SigningKey     string // Key used to sign cookie (and later for you to verify!)
		ExpiresMinutes int
		Path           string
		Domain         string
		SecureOnly     bool
		HTTPOnly       bool
	}
	// OIDC (OAuth 2) flow that allows an external site to securely login and verify the user
	OIDC struct {
		Enabled bool
		Clients map[string]struct { // ClientID Key
			Name        string
			Author      string
			ClientID    string
			RedirectURI []string
		}
	}
}

type ConfigWeb struct {
	Host         string
	Requirements ConfigWebRequirements
	RecaptchaV2  ConfigRecaptchaV2
	JWT          ConfigJWT
	Login        ConfigLogin
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
	Web            ConfigWeb            // Configure how the user interacts with the web
	Email          ConfigEmail          // SMTP/Email sending config
	Authenticators ConfigAuthencatorSet // Describes API Authenticators
	Production     bool                 // Production changes how logs are generated and tighter security checks
}
