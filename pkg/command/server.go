package command

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/Jeffail/gabs"
	"github.com/micro/cli"
	"github.com/oklog/run"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/owncloud/ocis-devldap/pkg/assets"
	"github.com/owncloud/ocis-devldap/pkg/config"
	"github.com/owncloud/ocis-devldap/pkg/flagset"
	"github.com/owncloud/ocis-devldap/pkg/server/debug"
	"github.com/owncloud/ocis-devldap/pkg/server/ldap"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start integrated server",
		Flags: flagset.ServerWithConfig(cfg),
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					exporter, err := ocagent.NewExporter(
						ocagent.WithReconnectionPeriod(5*time.Second),
						ocagent.WithAddress(cfg.Tracing.Endpoint),
						ocagent.WithServiceName(cfg.Tracing.Service),
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create agent tracing")

						return err
					}

					trace.RegisterExporter(exporter)
					view.RegisterExporter(exporter)

				case "jaeger":
					exporter, err := jaeger.NewExporter(
						jaeger.Options{
							AgentEndpoint:     cfg.Tracing.Endpoint,
							CollectorEndpoint: cfg.Tracing.Collector,
							ServiceName:       cfg.Tracing.Service,
						},
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create jaeger tracing")

						return err
					}

					trace.RegisterExporter(exporter)

				case "zipkin":
					endpoint, err := openzipkin.NewEndpoint(
						cfg.Tracing.Service,
						cfg.Tracing.Endpoint,
					)

					if err != nil {
						logger.Error().
							Err(err).
							Str("endpoint", cfg.Tracing.Endpoint).
							Str("collector", cfg.Tracing.Collector).
							Msg("Failed to create zipkin tracing")

						return err
					}

					exporter := zipkin.NewExporter(
						zipkinhttp.NewReporter(
							cfg.Tracing.Collector,
						),
						endpoint,
					)

					trace.RegisterExporter(exporter)

				default:
					logger.Warn().
						Str("type", t).
						Msg("Unknown tracing backend")
				}

				trace.ApplyConfig(
					trace.Config{
						DefaultSampler: trace.AlwaysSample(),
					},
				)
			} else {
				logger.Debug().
					Msg("Tracing is not enabled")
			}

			var (
				gr          = run.Group{}
				ctx, cancel = context.WithCancel(context.Background())
				//metrics     = metrics.New()
			)

			defer cancel()

			a := assets.New(assets.Config(cfg))

			// load the data from the assets

			d, err := a.Open(cfg.Asset.Data)
			if err != nil {
				return err
			}
			defer d.Close()

			data, err := gabs.ParseJSONBuffer(d)
			if err != nil {
				return err
			}

			{
				server, err := ldap.Server(
					ldap.Logger(logger),
					ldap.Context(ctx),
					ldap.Addr(cfg.LDAP.Addr),
					ldap.Data(data),
					ldap.Name("com.owncloud.ocis.devldap"),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("transport", "ldap").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)
					defer timeout()
					defer cancel()

					logger.Info().
						Str("transport", "ldap").
						Msg("Shutting down server")

					server.Shutdown(ctx)
				})
			}

			// load certificate
			crt, err := a.Open(cfg.Asset.Crt)
			if err != nil {
				return err
			}
			defer crt.Close()

			certPem, err := ioutil.ReadAll(crt)
			if err != nil {
				return err
			}

			key, err := a.Open(cfg.Asset.Key)
			if err != nil {
				return err
			}
			defer key.Close()

			keyPem, err := ioutil.ReadAll(key)
			if err != nil {
				return err
			}
			cert, err := tls.X509KeyPair(certPem, keyPem)
			if err != nil {
				return err
			}

			tlsConfig := tls.Config{
				Certificates: []tls.Certificate{cert},
			}

			{
				server, err := ldap.Server(
					ldap.Logger(logger),
					ldap.Context(ctx),
					ldap.Addr(cfg.LDAP.TLSAddr),
					ldap.Data(data),
					ldap.TLSConfig(&tlsConfig),
					ldap.Name("com.owncloud.ocis.devldaps"),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("transport", "ldaps").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServeTLS()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)
					defer timeout()
					defer cancel()

					logger.Info().
						Str("transport", "ldaps").
						Msg("Shutting down server")

					server.Shutdown(ctx)
				})
			}

			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("transport", "debug").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)

					defer timeout()
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						logger.Info().
							Err(err).
							Str("transport", "debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("transport", "debug").
							Msg("Shutting down server")
					}
				})
			}

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}
