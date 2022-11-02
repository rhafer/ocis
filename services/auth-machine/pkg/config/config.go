package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	GRPC GRPCConfig `yaml:"grpc"`

	TokenManager *TokenManager `yaml:"token_manager"`
	Reva         *shared.Reva  `yaml:"reva"`

	SkipUserGroupsInToken bool `yaml:"skip_user_groups_in_token" env:"AUTH_MACHINE_SKIP_USER_GROUPS_IN_TOKEN" desc:"Disables the encoding of the user's group memberships in the reva access token. This reduces the token size, especially when users are members of a large number of groups."`

	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;AUTH_MACHINE_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services."`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;AUTH_MACHINE_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;AUTH_MACHINE_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;AUTH_MACHINE_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;AUTH_MACHINE_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL;AUTH_MACHINE_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY;AUTH_MACHINE_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR;AUTH_MACHINE_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE;AUTH_MACHINE_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

type Service struct {
	Name string `yaml:"-"`
}

type Debug struct {
	Addr   string `yaml:"addr" env:"AUTH_MACHINE_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"AUTH_MACHINE_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"AUTH_MACHINE_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"AUTH_MACHINE_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

type GRPCConfig struct {
	Addr      string                 `yaml:"addr" env:"AUTH_MACHINE_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Namespace string                 `yaml:"-"`
	Protocol  string                 `yaml:"protocol" env:"AUTH_MACHINE_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service."`
}
