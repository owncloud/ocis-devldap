package ldap

import (
	"time"

	"github.com/butonic/ldapserver/pkg/constants"
	"github.com/butonic/ldapserver/pkg/ldap"
	"github.com/butonic/zerologr"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/micro/go-plugins/wrapper/trace/opencensus"
	"github.com/owncloud/ocis-devldap/pkg/version"
)

// Server initializes the ldap service and server.
func Server(opts ...Option) (*ldap.Server, error) {
	options := newOptions(opts...)
	options.Logger.Info().Str("addr", options.Addr).Msg("Server listening on")

	service := micro.NewService(
		micro.Name(options.Name),
		micro.Version(version.String),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Context(options.Context),
	)

	service.Init()

	zlog := zerologr.NewWithOptions(
		zerologr.Options{
			Name:   options.Name,
			Logger: &options.Logger.Logger,
		},
	)

	//Create a new LDAP Server
	server := ldap.NewServer(
		ldap.Addr(options.Addr),
		ldap.TLSConfig(options.TLSConfig),
		ldap.Logger(zlog),
	)
	h := &Handler{
		data: options.Data,
	}

	//Create routes bindings
	routes := ldap.NewRouteMux()
	routes.NotFound(h.NotFound)
	routes.Abandon(h.Abandon)
	routes.Bind(h.Bind)

	routes.Extended(h.WhoAmI).
		RequestName(constants.NoticeOfWhoAmI).Label("Ext - WhoAmI")

	routes.Extended(h.Extended).Label("Ext - Generic")

	routes.Search(h.Search).Label("Search - Generic")

	// Attach routes to server
	server.Handle(routes)

	return server, nil
}
