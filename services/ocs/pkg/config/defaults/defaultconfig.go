package defaults

import (
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9114",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9110",
			Root:      "/ocs",
			Namespace: "com.owncloud.web",
			CORS: config.CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		Service: config.Service{
			Name: "ocs",
		},
		AccountBackend: "cs3",
		Reva:           shared.DefaultRevaConfig(),
		IdentityManagement: config.IdentityManagement{
			Address: "https://localhost:9200",
		},
	}
}

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

	if cfg.CacheStore == nil && cfg.Commons != nil && cfg.Commons.CacheStore != nil {
		cfg.CacheStore = &config.CacheStore{
			Type:    cfg.Commons.CacheStore.Type,
			Address: cfg.Commons.CacheStore.Address,
			Size:    cfg.Commons.CacheStore.Size,
		}
	} else if cfg.CacheStore == nil {
		cfg.CacheStore = &config.CacheStore{}
	}

	if cfg.Reva == nil && cfg.Commons != nil && cfg.Commons.Reva != nil {
		cfg.Reva = &shared.Reva{
			Address:   cfg.Commons.Reva.Address,
			TLSMode:   cfg.Commons.Reva.TLSMode,
			TLSCACert: cfg.Commons.Reva.TLSCACert,
		}
	} else if cfg.Reva == nil {
		cfg.Reva = &shared.Reva{}
	}

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}

	if cfg.MachineAuthAPIKey == "" && cfg.Commons != nil && cfg.Commons.MachineAuthAPIKey != "" {
		cfg.MachineAuthAPIKey = cfg.Commons.MachineAuthAPIKey
	}

	if cfg.MicroGRPCClient == nil {
		cfg.MicroGRPCClient = &shared.MicroGRPCClient{}
		if cfg.Commons != nil && cfg.Commons.MicroGRPCClient != nil {
			cfg.MicroGRPCClient.TLSMode = cfg.Commons.MicroGRPCClient.TLSMode
			cfg.MicroGRPCClient.TLSCACert = cfg.Commons.MicroGRPCClient.TLSCACert
		}
	}
}

func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}
