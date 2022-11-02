package defaults

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/config"
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
			Addr:   "127.0.0.1:9141",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTPConfig{
			Addr:      "127.0.0.1:9140",
			Namespace: "com.owncloud.web",
			Protocol:  "tcp",
			Prefix:    "",
		},
		Service: config.Service{
			Name: "frontend",
		},
		Reva:                     shared.DefaultRevaConfig(),
		PublicURL:                "https://localhost:9200",
		EnableFavorites:          false,
		EnableProjectSpaces:      true,
		EnableShareJail:          true,
		UploadMaxChunkSize:       1e+8,
		UploadHTTPMethodOverride: "",
		DefaultUploadProtocol:    "tus",
		EnableResharing:          true,
		SearchMinLength:          3,
		Checksums: config.Checksums{
			SupportedTypes:      []string{"sha1", "md5", "adler32"},
			PreferredUploadType: "",
		},
		AppHandler: config.AppHandler{
			Prefix: "app",
		},
		Archiver: config.Archiver{
			Insecure:    false,
			Prefix:      "archiver",
			MaxNumFiles: 10000,
			MaxSize:     1073741824,
		},
		DataGateway: config.DataGateway{
			Prefix: "data",
		},
		OCS: config.OCS{
			Prefix:                  "ocs",
			SharePrefix:             "/Shares",
			HomeNamespace:           "/users/{{.Id.OpaqueId}}",
			AdditionalInfoAttribute: "{{.Mail}}",
			ResourceInfoCacheType:   "memory",
			ResourceInfoCacheTTL:    0,
		},
		Middleware: config.Middleware{
			Auth: config.Auth{
				CredentialsByUserAgent: map[string]string{},
			},
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

	if cfg.Reva == nil && cfg.Commons != nil && cfg.Commons.Reva != nil {
		cfg.Reva = &shared.Reva{
			Address:   cfg.Commons.Reva.Address,
			TLS: cfg.Commons.Reva.TLS,
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

	if cfg.TransferSecret == "" && cfg.Commons != nil && cfg.Commons.TransferSecret != "" {
		cfg.TransferSecret = cfg.Commons.TransferSecret
	}

	if cfg.MachineAuthAPIKey == "" && cfg.Commons != nil && cfg.Commons.MachineAuthAPIKey != "" {
		cfg.MachineAuthAPIKey = cfg.Commons.MachineAuthAPIKey
	}

}

func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
