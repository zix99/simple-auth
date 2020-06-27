package config

type ConfigDatabase struct {
	Driver string
	URL    string
}

type ConfigAuthencatorSet struct {
	Token struct {
		Enabled bool
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
