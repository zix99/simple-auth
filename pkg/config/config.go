package config

type ConfigDatabase struct {
	Driver string
	URL    string
}

type ConfigAuthenticator struct {
	Enabled bool
}

type ConfigAuthencatorSet struct {
	Exchange struct {
		ConfigAuthenticator
	}
}

type ConfigWebRequirements struct {
	PasswordMinLength int
	PasswordMaxLength int
	UsernameMinLength int
	UsernameMaxLength int
}

type ConfigWeb struct {
	Host         string
	Requirements ConfigWebRequirements
	Metadata     map[string]interface{}
}

type Config struct {
	Db             ConfigDatabase
	Web            ConfigWeb
	Authenticators ConfigAuthencatorSet
	Production     bool
}
