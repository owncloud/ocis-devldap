package config

type Log struct {
	Level string
}

type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

type LDAP struct {
	Addr string
	Root string
	Data string
}

type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

type Asset struct {
	Path string
}

type Config struct {
	File    string
	Log     Log
	Debug   Debug
	LDAP    LDAP
	Tracing Tracing
	Asset   Asset
}

func New() *Config {
	return &Config{}
}
