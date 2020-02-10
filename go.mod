module github.com/owncloud/ocis-devldap

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.0
	contrib.go.opencensus.io/exporter/ocagent v0.6.0
	contrib.go.opencensus.io/exporter/zipkin v0.1.1
	github.com/Jeffail/gabs v1.4.0
	github.com/butonic/ldapserver v0.0.0-20191209092749-8fb2e7a4a628
	github.com/butonic/zerologr v0.0.0-20191210074216-d798ee237d84
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.1 // indirect
	github.com/lor00x/goldap v0.0.0-20180618054307-a546dffdd1a3
	github.com/micro/cli/v2 v2.1.1
	github.com/micro/go-micro/v2 v2.0.0
	github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2 v2.0.1
	github.com/micro/go-plugins/wrapper/trace/opencensus/v2 v2.0.1
	github.com/oklog/run v1.0.0
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/owncloud/ocis-pkg/v2 v2.0.1
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.5.0
	github.com/uber/jaeger-client-go v2.20.1+incompatible // indirect
	go.opencensus.io v0.22.2
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

go 1.13
