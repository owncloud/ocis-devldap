package assets

import (
	"net/http"
	"os"
	"path"

	"github.com/micro/go-micro/util/log"
	"github.com/owncloud/ocis-devldap/pkg/config"
)

//go:generate go run github.com/UnnoTed/fileb0x embed.yml

// assets gets initialized by New and provides the handler.
type assets struct {
	config *config.Config
}

// Open just implements the HTTP filesystem interface.
func (a assets) Open(original string) (http.File, error) {
	if a.config.Asset.Path != "" {
		if stat, err := os.Stat(a.config.Asset.Path); err == nil && stat.IsDir() {
			custom := path.Join(
				a.config.Asset.Path,
				original,
			)

			if _, err := os.Stat(custom); !os.IsNotExist(err) {
				f, err := os.Open(custom)

				if err != nil {
					return nil, err
				}

				return f, nil
			}
		} else {
			log.Warnf("Asset directory [%s] doesn't exist", a.config.Asset.Path)
		}
	}

	return FS.OpenFile(
		CTX,
		original,
		os.O_RDONLY,
		0644,
	)
}

// New returns a new http filesystem to serve assets.
func New(opts ...Option) http.FileSystem {
	options := newOptions(opts...)

	return assets{
		config: options.Config,
	}
}
