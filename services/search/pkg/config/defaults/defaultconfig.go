package defaults

import (
	"path"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:  "127.0.0.1:9224",
			Token: "",
		},
		GRPC: config.GRPC{
			Addr:      "127.0.0.1:9220",
			Namespace: "com.owncloud.api",
		},
		Service: config.Service{
			Name: "search",
		},
		Datapath: path.Join(defaults.BaseDataPath(), "search"),
		Reva:     shared.DefaultRevaConfig(),
		Events: config.Events{
			Endpoint:      "127.0.0.1:9233",
			Cluster:       "ocis-cluster",
			ConsumerGroup: "search",
			EnableTLS:     false,
		},
		MachineAuthAPIKey: "",
	}
}

// EnsureDefaults adds default values to the configuration if they are not set yet
func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for "envdecode".
	if cfg.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Log = &config.Log{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Log == nil {
		cfg.Log = &config.Log{}
	}
	// provide with defaults for shared tracing, since we need a valid destination address for "envdecode".
	if cfg.Tracing == nil && cfg.Commons != nil && cfg.Commons.Tracing != nil {
		cfg.Tracing = &config.Tracing{
			Enabled:   cfg.Commons.Tracing.Enabled,
			Type:      cfg.Commons.Tracing.Type,
			Endpoint:  cfg.Commons.Tracing.Endpoint,
			Collector: cfg.Commons.Tracing.Collector,
		}
	} else if cfg.Tracing == nil {
		cfg.Tracing = &config.Tracing{}
	}

	if cfg.MachineAuthAPIKey == "" && cfg.Commons != nil && cfg.Commons.MachineAuthAPIKey != "" {
		cfg.MachineAuthAPIKey = cfg.Commons.MachineAuthAPIKey
	}

	if cfg.Reva == nil && cfg.Commons != nil && cfg.Commons.Reva != nil {
		cfg.Reva = &shared.Reva{
			Address:   cfg.Commons.Reva.Address,
			TLS: cfg.Commons.Reva.TLS,
		}
	} else if cfg.Reva == nil {
		cfg.Reva = &shared.Reva{}
	}
	if cfg.GRPCClientTLS == nil {
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}
		if cfg.Commons != nil && cfg.Commons.GRPCClientTLS != nil {
			cfg.GRPCClientTLS.Mode = cfg.Commons.GRPCClientTLS.Mode
			cfg.GRPCClientTLS.CACert = cfg.Commons.GRPCClientTLS.CACert
		}
	}
	if cfg.GRPCServiceTLS == nil {
		cfg.GRPCServiceTLS = &shared.GRPCServiceTLS{}
		if cfg.Commons != nil && cfg.Commons.GRPCServiceTLS != nil {
			cfg.GRPCServiceTLS.Enabled = cfg.Commons.GRPCServiceTLS.Enabled
			cfg.GRPCServiceTLS.Cert = cfg.Commons.GRPCServiceTLS.Cert
			cfg.GRPCServiceTLS.Key = cfg.Commons.GRPCServiceTLS.Key
		}
	}
}

// Sanitize sanitizes the configuration
func Sanitize(cfg *config.Config) {
	// no http endpoint to be sanitized
}
