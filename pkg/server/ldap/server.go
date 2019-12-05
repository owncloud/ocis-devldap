package ldap

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/Jeffail/gabs"
	ls "github.com/butonic/ldapserver"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/micro/go-plugins/wrapper/trace/opencensus"
	"github.com/owncloud/ocis-devldap/pkg/assets"
	"github.com/owncloud/ocis-devldap/pkg/version"
)

func loadData(file io.Reader) (*gabs.Container, error) {
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	data, err := gabs.ParseJSON(raw)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Server(opts ...Option) (*ls.Server, error) {
	options := newOptions(opts...)
	log.Infof("Server [ldap] listening on [%s]", options.Config.LDAP.Addr)

	// &cli.StringFlag{
	// 	Name:        "ldap-addr",
	// 	Value:       "0.0.0.0:10389",
	// 	Usage:       "Address to bind ldap server",
	// 	EnvVar:      "DEVLDAP_LDAP_ADDR",
	// 	Destination: &cfg.LDAP.Addr,
	// },

	service := micro.NewService(
		micro.Name("com.owncloud.ocis.devldap"),
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

	a := assets.New(assets.Config(options.Config))
	d, err := a.Open(options.Config.LDAP.Data)
	if err != nil {
		return nil, err
	}

	data, err := loadData(d)
	if err != nil {
		return nil, err
	}
	//Create a new LDAP Server
	server := ls.NewServer()
	handler := &Handler{
		data: data,
	}

	//Create routes bindings
	routes := ls.NewRouteMux()
	routes.NotFound(handler.NotFound)
	routes.Abandon(handler.Abandon)
	routes.Bind(handler.Bind)

	routes.Extended(handler.WhoAmI).
		RequestName(ls.NoticeOfWhoAmI).Label("Ext - WhoAmI")

	routes.Extended(handler.Extended).Label("Ext - Generic")

	routes.Search(handler.Search).Label("Search - Generic")

	// Attach routes to server
	server.Handle(routes)

	return server, nil
}
