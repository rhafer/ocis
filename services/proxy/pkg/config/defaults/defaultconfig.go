package defaults

import (
	"path"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
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
			Addr:  "127.0.0.1:9205",
			Token: "",
		},
		HTTP: config.HTTP{
			Addr:      "0.0.0.0:9200",
			Root:      "/",
			Namespace: "com.owncloud.web",
			TLSCert:   path.Join(defaults.BaseDataPath(), "proxy", "server.crt"),
			TLSKey:    path.Join(defaults.BaseDataPath(), "proxy", "server.key"),
			TLS:       true,
		},
		Service: config.Service{
			Name: "proxy",
		},
		OIDC: config.OIDC{
			Issuer: "https://localhost:9200",

			AccessTokenVerifyMethod: config.AccessTokenVerificationJWT,
			UserinfoCache: config.UserinfoCache{
				Size: 1024,
				TTL:  10,
			},
			JWKS: config.JWKS{
				RefreshInterval:   60, // minutes
				RefreshRateLimit:  60, // seconds
				RefreshTimeout:    10, // seconds
				RefreshUnknownKID: true,
			},
		},
		PolicySelector: nil,
		Reva:           shared.DefaultRevaConfig(),
		PreSignedURL: config.PreSignedURL{
			AllowedHTTPMethods: []string{"GET"},
			Enabled:            true,
		},
		AccountBackend:        "cs3",
		UserOIDCClaim:         "preferred_username",
		UserCS3Claim:          "username",
		AutoprovisionAccounts: false,
		EnableBasicAuth:       false,
		InsecureBackends:      false,
	}
}

func DefaultPolicies() []config.Policy {
	return []config.Policy{
		{
			Name: "ocis",
			Routes: []config.Route{
				{
					Endpoint:    "/",
					Backend:     "http://localhost:9100",
					Unprotected: true,
				},
				{
					Endpoint:    "/.well-known/",
					Backend:     "http://localhost:9130",
					Unprotected: true,
				},
				{
					Endpoint:    "/konnect/",
					Backend:     "http://localhost:9130",
					Unprotected: true,
				},
				{
					Endpoint:    "/signin/",
					Backend:     "http://localhost:9130",
					Unprotected: true,
				},
				{
					Endpoint: "/archiver",
					Backend:  "http://localhost:9140",
				},
				{
					Type:     config.RegexRoute,
					Endpoint: "/ocs/v[12].php/cloud/user/signing-key", // only `user/signing-key` is left in ocis-ocs
					Backend:  "http://localhost:9110",
				},
				{
					Type:        config.RegexRoute,
					Endpoint:    "/ocs/v[12].php/config",
					Backend:     "http://localhost:9140",
					Unprotected: true,
				},
				{
					Endpoint: "/ocs/",
					Backend:  "http://localhost:9140",
				},
				{
					Type:     config.QueryRoute,
					Endpoint: "/remote.php/?preview=1",
					Backend:  "http://localhost:9115",
				},
				// TODO the actual REPORT goes to /dav/files/{username}, which is user specific ... how would this work in a spaces world?
				// TODO what paths are returned? the href contains the full path so it should be possible to return urls from other spaces?
				// TODO or we allow a REPORT on /dav/spaces to search all spaces and /dav/space/{spaceid} to search a specific space
				// send webdav REPORT requests to search service
				{
					Method:   "REPORT",
					Endpoint: "/remote.php/dav/",
					Backend:  "http://localhost:9115", // TODO use registry?
				},
				{
					Method:   "REPORT",
					Endpoint: "/remote.php/webdav",
					Backend:  "http://localhost:9115", // TODO use registry?
				},
				{
					Method:   "REPORT",
					Endpoint: "/dav/spaces",
					Backend:  "http://localhost:9115",
				},
				{
					Type:     config.QueryRoute,
					Endpoint: "/dav/?preview=1",
					Backend:  "http://localhost:9115",
				},
				{
					Type:     config.QueryRoute,
					Endpoint: "/webdav/?preview=1",
					Backend:  "http://localhost:9115",
				},
				{
					Endpoint: "/remote.php/",
					Service:  "com.owncloud.web.ocdav",
				},
				{
					Endpoint: "/dav/",
					Service:  "com.owncloud.web.ocdav",
				},
				{
					Endpoint: "/webdav/",
					Service:  "com.owncloud.web.ocdav",
				},
				{
					Endpoint:    "/status",
					Service:     "com.owncloud.web.ocdav",
					Unprotected: true,
				},
				{
					Endpoint:    "/status.php",
					Service:     "com.owncloud.web.ocdav",
					Unprotected: true,
				},
				{
					Endpoint: "/index.php/",
					Service:  "com.owncloud.web.ocdav",
				},
				{
					Endpoint: "/apps/",
					Service:  "com.owncloud.web.ocdav",
				},
				{
					Endpoint:    "/data",
					Backend:     "http://localhost:9140",
					Unprotected: true,
				},
				{
					Endpoint:    "/app/list",
					Backend:     "http://localhost:9140",
					Unprotected: true,
				},
				{
					Endpoint: "/app/", // /app or /apps? ocdav only handles /apps
					Backend:  "http://localhost:9140",
				},
				{
					Endpoint: "/graph/",
					Backend:  "http://localhost:9120",
				},
				{
					Endpoint: "/api/v0/settings",
					Backend:  "http://localhost:9190",
				},
				{
					Endpoint:    "/settings.js",
					Backend:     "http://localhost:9190",
					Unprotected: true,
				},
			},
		},
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

	if cfg.Reva == nil && cfg.Commons != nil && cfg.Commons.Reva != nil {
		cfg.Reva = &shared.Reva{
			Address:   cfg.Commons.Reva.Address,
			TLSMode:   cfg.Commons.Reva.TLSMode,
			TLSCACert: cfg.Commons.Reva.TLSCACert,
		}
	} else if cfg.Reva == nil {
		cfg.Reva = &shared.Reva{}
	}

	if cfg.MicroGRPCClient == nil {
		cfg.MicroGRPCClient = &shared.MicroGRPCClient{}
		if cfg.Commons != nil && cfg.Commons.MicroGRPCClient != nil {
			cfg.MicroGRPCClient.TLSMode = cfg.Commons.MicroGRPCClient.TLSMode
			cfg.MicroGRPCClient.TLSCACert = cfg.Commons.MicroGRPCClient.TLSCACert
		}
	}
}

// Sanitize sanitizes the configuration
func Sanitize(cfg *config.Config) {
	if cfg.Policies == nil {
		cfg.Policies = DefaultPolicies()
	}

	if cfg.PolicySelector == nil {
		cfg.PolicySelector = &config.PolicySelector{
			Static: &config.StaticSelectorConf{
				Policy: "ocis",
			},
		}
	}

	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}
