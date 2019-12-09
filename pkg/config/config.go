package config

// Log defines the available logging configuration.
type Log struct {
	Level  string
	Pretty bool
	Color  bool
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

// LDAP defines the available ldap configuration.
type LDAP struct {
	Addr string
	Data string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Config combines all available configuration parts.
type Config struct {
	File    string
	Log     Log
	Debug   Debug
	LDAP    LDAP
	Tracing Tracing
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
