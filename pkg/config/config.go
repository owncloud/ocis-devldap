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
	Addr    string
	TLSAddr string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string
	Data string
	Crt  string
	Key  string
}

// Config combines all available configuration parts.
type Config struct {
	File    string
	Log     Log
	Debug   Debug
	LDAP    LDAP
	Asset   Asset
	Tracing Tracing
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}
