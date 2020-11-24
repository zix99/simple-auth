package config

// ConfigDatabase holds database-specific configuration
type ConfigDatabase struct {
	Driver string // Driver, eg "sqlite3", "postgres", or "mysql"
	URL    string // Connection string for the driver. See http://gorm.io/docs/connecting_to_the_database.html
	Debug  bool   // Enable debug logging
}

type ConfigTokenAuthenticator struct {
	Enabled                    bool
	SessionExpiresMinutes      int
	VerificationExpiresSeconds int
}

type ConfigSimpleAuthenticator struct {
	Enabled      bool
	SharedSecret string // If non-empty, will be required as a Bearer token in the Authorization header. If empty, anyone can use this endpoint (if enabled)
}

type ConfigVouchAuthenticator struct {
	Enabled bool
}

// Authenticators are how someone external to SA can authenticate with it
type ConfigAuthencatorSet struct {
	Token  ConfigTokenAuthenticator
	Simple ConfigSimpleAuthenticator
	Vouch  ConfigVouchAuthenticator
}

type ConfigWebRequirements struct {
	UsernameRegex     string // Regex match for allowed username characters (server & client enforced)
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
	SigningMethod  string // As defined by go-jwt. Commonly HS256, HS512, RS256, RS512
	SigningKey     string // Key used to sign cookie (and later for you to verify!). If RS based, will be parsed as PEM
	ExpiresMinutes int
	Issuer         string
}

type ConfigLoginGateway struct {
	Enabled    bool
	Targets    []string
	Host       string            // Override the host header
	LogoutPath string            // Path for the logout url (to override & skip proxying)
	Rewrite    map[string]string // Rewrite URLs upon proxying eg "/old"->"/new" or "/api/*"->"/$1"
	Headers    map[string]string // Override additional headers (excluding host header)
	NoCache    bool              // If true, will attempt to disable caching to gateway target
}

type ConfigLoginCookie struct {
	Name       string // Name of the cookie
	JWT        ConfigJWT
	Path       string
	Domain     string
	SecureOnly bool
	HTTPOnly   bool
}

type ConfigOIDCProvider struct {
	ID           string
	Name         string // Display name
	Icon         string
	ClientID     string
	ClientSecret string
	AuthURL      string // URL to redirect the user to for auth
	TokenURL     string // URL to trade code for token
}

type ConfigLoginSettings struct {
	CreateAccountEnabled    bool
	EmailValidationRequired bool
	RouteOnLogin            string
	AllowedContinueUrls     []string
	ThrottleDuration        string // Parsed as Duration, represents a delay from any major action (Helps mitigate brute-force attacks)
}

type OneTimeConfig struct {
	Enabled             bool
	AllowForgotPassword bool   // Separate from enabled, will allow issuing one-time via email
	TokenDuration       string // Parsed as duration
}

type TwoFactorConfig struct {
	Enabled   bool
	KeyLength int
	Issuer    string
	Drift     int
}

type ConfigLogin struct {
	Settings ConfigLoginSettings
	// SameDomain authentication uses a cookie set to a domain (and presumably shared with your site).  Easiest to implement in a full-trust environment
	Cookie ConfigLoginCookie
	// OIDC (OAuth 2) flow that allows an external site to securely login and verify the user
	OIDC []*ConfigOIDCProvider
	// Configuration for one-time password (eg. forgotten password)
	OneTime OneTimeConfig
	// 2FA/TOTP Configuration
	TwoFactor TwoFactorConfig
}

type ConfigMetadata struct {
	Company string
	Footer  string
	Bucket  map[string]interface{}
}

type ConfigWebTLS struct {
	Enabled  bool
	Auto     bool   // Auto get certificate via Let's Encrypt
	Cache    string // If auto TLS, directory certs to be stored
	CertFile string // If not auto, cert
	KeyFile  string // If not auto, key
}

type ConfigWeb struct {
	Host         string
	BaseURL      string // If empty, will attempt to be derivied
	TLS          ConfigWebTLS
	Requirements ConfigWebRequirements
	RecaptchaV2  ConfigRecaptchaV2
	Login        ConfigLogin
	Gateway      ConfigLoginGateway
	Prometheus   bool // If true, will enable /metrics endpoint
}

type ConfigEmailSMTP struct {
	Host     string
	Identity string
	Username string
	Password string
	From     string
}

type ConfigEmail struct {
	Enabled bool
	SMTP    ConfigEmailSMTP
}

type ConfigAPI struct {
	External     bool // If true, allows external API calls (outside of session API)
	SharedSecret string
}

// Config represents the root configuration
type Config struct {
	Include        []string
	Metadata       ConfigMetadata
	Db             ConfigDatabase
	Web            ConfigWeb            // Configure how the user interacts with the web
	Email          ConfigEmail          // SMTP/Email sending config
	Authenticators ConfigAuthencatorSet // Describes API Authenticators
	API            ConfigAPI            // API configuration
	Production     bool                 // Production changes how logs are generated and tighter security checks
	Verbose        bool                 // Turns on additional logging
	StaticFromDisk bool                 // Checks the disk for static files
}
