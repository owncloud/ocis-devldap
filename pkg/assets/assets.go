package assets

import (
	"io/ioutil"
	"os"

	"github.com/owncloud/ocis-phoenix/pkg/config"
	"github.com/owncloud/ocis-pkg/log"
)

//go:generate go run github.com/UnnoTed/fileb0x embed.yml

// assets gets initialized by New and provides the handler.
type assets struct {
	logger log.Logger
	config *config.Config
}

// Data simply provides ldap server data file.
func (a assets) Data() ([]byte, error) {
	if a.config.LDAP.Data != "" {
		if _, err := os.Stat(a.config.LDAP.Data); os.IsNotExist(err) {
			a.logger.Error().
				Str("path", a.config.LDAP.Data).
				Msg("Data file doesn't exist")

			return []byte(), err
		}

		content, err := ioutil.ReadFile(a.config.LDAP.Data)

		if err != nil {
			a.logger.Error().
				Str("path", a.config.LDAP.Data).
				Msg("Failed to read data file")

			return []byte(), err
		}

		return content, nil
	}

	return ReadFile("data.json")
}

// New returns a new handler to serve assets.
func New(opts ...Option) assets {
	options := newOptions(opts...)

	return assets{
		logger: options.Logger,
		config: options.Config,
	}
}
