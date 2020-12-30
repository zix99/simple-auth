package config

import "regexp"

// ConfigDatabase holds database-specific configuration
type ConfigDatabase struct {
	Driver string
	URL    string
	Debug  bool
}

type ConfigMetadata struct {
	Company string
	Footer  string
	TagLine string
	Bucket  map[string]interface{}
}

// Authenticators
type (
	ConfigTokenAuthenticator struct {
		Enabled                    bool
		SessionExpiresMinutes      int
		VerificationExpiresSeconds int
	}

	ConfigSimpleAuthenticator struct {
		Enabled      bool
		SharedSecret string // If non-empty, will be required as a Bearer token in the Authorization header. If empty, anyone can use this endpoint (if enabled)
	}

	ConfigVouchAuthenticator struct {
		Enabled bool
	}

	// Authenticators are how someone external to SA can authenticate with it
	ConfigAuthenticatorSet struct {
		Token  ConfigTokenAuthenticator
		Simple ConfigSimpleAuthenticator
		Vouch  ConfigVouchAuthenticator
	}
)

// Web
type (
	ConfigRecaptchaV2 struct {
		Enabled bool
		SiteKey string // Site key, as provided by google
		Secret  string // Secret for google v2 api (server side call)
		Theme   string // Themes (dark or light)
	}

	ConfigJWT struct {
		SigningMethod  string // As defined by go-jwt. Commonly HS256, HS512, RS256, RS512
		SigningKey     string // Key used to sign cookie (and later for you to verify!). If RS based, will be parsed as PEM
		ExpiresMinutes int
		Issuer         string
	}

	ConfigLoginGateway struct {
		Enabled    bool
		Targets    []string
		Host       string
		LogoutPath string
		Rewrite    map[string]string
		Headers    map[string]string
		NoCache    bool
	}

	ConfigLoginCookie struct {
		Name       string // Name of the cookie
		JWT        ConfigJWT
		Path       string
		Domain     string
		SecureOnly bool
		HTTPOnly   bool
	}

	ConfigLoginSettings struct {
		RouteOnLogin               string
		AllowedContinueUrls        []string // List of allowed regex's
		_allowedContinueUrlsRegexp []*regexp.Regexp
		ThrottleDuration           string // Parsed as Duration, represents a delay from any major action (Helps mitigate brute-force attacks)
	}

	OneTimeConfig struct {
		Enabled             bool
		AllowForgotPassword bool   // Separate from enabled, will allow issuing one-time via email
		TokenDuration       string // Parsed as duration
	}

	ConfigLogin struct {
		Settings ConfigLoginSettings
		// SameDomain authentication uses a cookie set to a domain (and presumably shared with your site).  Easiest to implement in a full-trust environment
		Cookie ConfigLoginCookie
		// Configuration for one-time password (eg. forgotten password)
		OneTime OneTimeConfig
	}

	ConfigWebTLS struct {
		Enabled  bool
		Auto     bool   // Auto get certificate via Let's Encrypt
		Cache    string // If auto TLS, directory certs to be stored
		CertFile string // If not auto, cert
		KeyFile  string // If not auto, key
	}

	ConfigWeb struct {
		Host        string
		BaseURL     string // If empty, will attempt to be derivied
		TLS         ConfigWebTLS
		RecaptchaV2 ConfigRecaptchaV2
		Login       ConfigLogin
		Gateway     ConfigLoginGateway
		Prometheus  bool
		Swagger     bool
	}
)

// Providers
type (
	ConfigTwoFactor struct {
		Enabled   bool
		KeyLength int
		Issuer    string
		Drift     int
	}

	ConfigOIDCProvider struct {
		ID           string
		Name         string // Display name
		Icon         string
		ClientID     string
		ClientSecret string
		AuthURL      string // URL to redirect the user to for auth
		TokenURL     string // URL to trade code for token
	}

	ConfigLocalLoginRequirements struct {
		UsernameRegex     string // Regex match for allowed username characters (server & client enforced)
		PasswordMinLength int
		PasswordMaxLength int
		UsernameMinLength int
		UsernameMaxLength int
	}

	ConfigLocalProvider struct {
		EmailValidationRequired bool
		Requirements            ConfigLocalLoginRequirements
		TwoFactor               ConfigTwoFactor
	}

	ConfigProviderSettings struct {
		CreateAccountEnabled bool
	}

	ConfigProviders struct {
		Settings ConfigProviderSettings
		Local    ConfigLocalProvider
		OIDC     []*ConfigOIDCProvider // OIDC (OAuth 2) flow that allows an external site to securely login and verify the user
	}
)

// Email
type (
	ConfigEmailSMTP struct {
		Host     string
		Identity string
		Username string
		Password string
	}

	ConfigEmail struct {
		Engine string
		From   string
		SMTP   ConfigEmailSMTP
	}
)

type ConfigAPI struct {
	External         bool // If true, allows external API calls (outside of session API)
	SharedSecret     string
	ThrottleDuration string // Parsed as Duration, represents a delay from any major action (Helps mitigate brute-force attacks)
}

// Config represents the root configuration
type Config struct {
	Include        []string
	Metadata       ConfigMetadata
	Db             ConfigDatabase
	Web            ConfigWeb              // Configure how the user interacts with the web
	Email          ConfigEmail            // SMTP/Email sending config
	Providers      ConfigProviders        // Login providers
	Authenticators ConfigAuthenticatorSet // Describes API Authenticators
	API            ConfigAPI              // API configuration
	Production     bool                   // Production changes how logs are generated and tighter security checks
	Verbose        bool                   // Turns on additional logging
	StaticFromDisk bool                   // Checks the disk for static files

	// Meta config
	Version bool // Show version
	Help    bool // Show help
}
