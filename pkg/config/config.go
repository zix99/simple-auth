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

type ConfigAuthencatorSet struct {
	Token  ConfigTokenAuthenticator
	Simple ConfigSimpleAuthenticator
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
}

type ConfigLoginCookie struct {
	Name       string // Name of the cookie
	JWT        ConfigJWT
	Path       string
	Domain     string
	SecureOnly bool
	HTTPOnly   bool
}

type ConfigLoginOIDC struct {
	Enabled bool
	Clients map[string]struct { // ClientID Key
		Name        string
		Author      string
		ClientID    string
		RedirectURI []string
	}
}

type ConfigLogin struct {
	// SameDomain authentication uses a cookie set to a domain (and presumably shared with your site).  Easiest to implement in a full-trust environment
	Cookie ConfigLoginCookie
	// Gateway
	Gateway ConfigLoginGateway
	// OIDC (OAuth 2) flow that allows an external site to securely login and verify the user
	OIDC ConfigLoginOIDC
	// Login Rules
	CreateAccountEnabled bool
	RouteOnLogin         string
	AllowedContinueUrls  []string
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
	TLS          ConfigWebTLS
	Requirements ConfigWebRequirements
	RecaptchaV2  ConfigRecaptchaV2
	Login        ConfigLogin
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
	Metadata       ConfigMetadata
	Db             ConfigDatabase
	Web            ConfigWeb            // Configure how the user interacts with the web
	Email          ConfigEmail          // SMTP/Email sending config
	Authenticators ConfigAuthencatorSet // Describes API Authenticators
	Production     bool                 // Production changes how logs are generated and tighter security checks
}
