package middleware

import (
	"net/http"
	"strings"

	"github.com/owncloud/ocis-devldap/pkg/assets"
)

// Static is a middleware that serves static assets.
func Static(opts ...Option) func(http.Handler) http.Handler {
	options := newOptions(opts...)

	static := http.StripPrefix(
		"/",
		http.FileServer(
			assets.New(
				assets.Config(options.Config),
			),
		),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
			} else {
				if strings.HasSuffix(r.URL.Path, "/") {
					http.NotFound(w, r)
				} else {
					static.ServeHTTP(w, r)
				}
			}
		})
	}
}
