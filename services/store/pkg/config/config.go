package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	GRPC GRPC `yaml:"grpc"`

	MicroGRPCClient  *shared.MicroGRPCClient  `yaml:"micro_grpc_client"`
	MicroGRPCService *shared.MicroGRPCService `yaml:"micro_grpc_service"`

	Datapath string `yaml:"data_path" env:"STORE_DATA_PATH" desc:"The directory where the filesystem storage will store ocis settings. If not definied, the root directory derives from $OCIS_BASE_DATA_PATH:/store."`

	Context context.Context `yaml:"-"`
}
