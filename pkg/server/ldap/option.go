package ldap

import (
	"context"
	"crypto/tls"

	"github.com/Jeffail/gabs"
	"github.com/owncloud/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger    log.Logger
	Context   context.Context
	Name      string
	Addr      string
	TLSConfig *tls.Config
	Data      *gabs.Container
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// Context provides a function to set the context option.
func Context(val context.Context) Option {
	return func(o *Options) {
		o.Context = val
	}
}

// Name provides a function to set the name option.
func Name(val string) Option {
	return func(o *Options) {
		o.Name = val
	}
}

// Addr provides a function to set the addr option.
func Addr(val string) Option {
	return func(o *Options) {
		o.Addr = val
	}
}

// TLSConfig provides a function to set the TLSConfig option.
func TLSConfig(val *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = val
	}
}

// Data provides a function to set the data option.
func Data(val *gabs.Container) Option {
	return func(o *Options) {
		o.Data = val
	}
}
