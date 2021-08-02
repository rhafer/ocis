package flagset

import (
	"path"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	pkgos "github.com/owncloud/ocis/ocis-pkg/os"
	"github.com/owncloud/ocis/proxy/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"PROXY_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"PROXY_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"PROXY_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode",
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       cfg.Debug.Addr,
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"PROXY_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-file",
			Usage:       "Enable log to file",
			EnvVars:     []string{"PROXY_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.StringFlag{
			Name:        "config-file",
			Value:       cfg.File,
			Usage:       "Path to config file",
			EnvVars:     []string{"PROXY_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.StringFlag{
			Name:        "policies-file",
			Value:       cfg.PoliciesFile,
			Usage:       "Path to policies file",
			EnvVars:     []string{"PROXY_POLICIES_FILE"},
			Destination: &cfg.PoliciesFile,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"PROXY_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       cfg.Tracing.Type,
			Usage:       "Tracing backend type",
			EnvVars:     []string{"PROXY_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       cfg.Tracing.Endpoint,
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"PROXY_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       cfg.Tracing.Collector,
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"PROXY_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       cfg.Tracing.Service,
			Usage:       "Service name for tracing",
			EnvVars:     []string{"PROXY_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       cfg.Debug.Addr,
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"PROXY_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       cfg.Debug.Token,
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"PROXY_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"PROXY_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"PROXY_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       cfg.HTTP.Addr,
			Usage:       "Address to bind http server",
			EnvVars:     []string{"PROXY_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       cfg.HTTP.Root,
			Usage:       "Root path of http server",
			EnvVars:     []string{"PROXY_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       cfg.Asset.Path,
			Usage:       "Path to custom assets",
			EnvVars:     []string{"PROXY_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "service-namespace",
			Value:       cfg.Service.Namespace,
			Usage:       "Set the base namespace for the service namespace",
			EnvVars:     []string{"PROXY_SERVICE_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "service-name",
			Value:       cfg.Service.Name,
			Usage:       "Service name",
			EnvVars:     []string{"PROXY_SERVICE_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "transport-tls-cert",
			Value:       flags.OverrideDefaultString(cfg.HTTP.TLSCert, path.Join(pkgos.MustUserConfigDir("ocis", "proxy"), "server.crt")),
			Usage:       "Certificate file for transport encryption",
			EnvVars:     []string{"PROXY_TRANSPORT_TLS_CERT"},
			Destination: &cfg.HTTP.TLSCert,
		},
		&cli.StringFlag{
			Name:        "transport-tls-key",
			Value:       flags.OverrideDefaultString(cfg.HTTP.TLSKey, path.Join(pkgos.MustUserConfigDir("ocis", "proxy"), "server.key")),
			Usage:       "Secret file for transport encryption",
			EnvVars:     []string{"PROXY_TRANSPORT_TLS_KEY"},
			Destination: &cfg.HTTP.TLSKey,
		},
		&cli.BoolFlag{
			Name:        "tls",
			Value:       cfg.HTTP.TLS,
			Usage:       "Use TLS (disable only if proxy is behind a TLS-terminating reverse-proxy).",
			EnvVars:     []string{"PROXY_TLS"},
			Destination: &cfg.HTTP.TLS,
		},
		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       cfg.TokenManager.JWTSecret,
			Usage:       "Used to create JWT to talk to reva, should equal reva's jwt-secret",
			EnvVars:     []string{"PROXY_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       cfg.Reva.Address,
			Usage:       "REVA Gateway Endpoint",
			EnvVars:     []string{"PROXY_REVA_GATEWAY_ADDR"},
			Destination: &cfg.Reva.Address,
		},
		&cli.BoolFlag{
			Name:        "insecure",
			Value:       cfg.InsecureBackends,
			Usage:       "allow insecure communication to upstream servers",
			EnvVars:     []string{"PROXY_INSECURE_BACKENDS"},
			Destination: &cfg.InsecureBackends,
		},

		// OIDC

		&cli.StringFlag{
			Name:        "oidc-issuer",
			Value:       cfg.OIDC.Issuer,
			Usage:       "OIDC issuer",
			EnvVars:     []string{"PROXY_OIDC_ISSUER", "OCIS_URL"}, // PROXY_OIDC_ISSUER takes precedence over OCIS_URL
			Destination: &cfg.OIDC.Issuer,
		},
		&cli.BoolFlag{
			Name:        "oidc-insecure",
			Value:       cfg.OIDC.Insecure,
			Usage:       "OIDC allow insecure communication",
			EnvVars:     []string{"PROXY_OIDC_INSECURE"},
			Destination: &cfg.OIDC.Insecure,
		},
		&cli.IntFlag{
			Name:        "oidc-userinfo-cache-tll",
			Value:       cfg.OIDC.UserinfoCache.TTL,
			Usage:       "Fallback TTL in seconds for caching userinfo, when no token lifetime can be identified",
			EnvVars:     []string{"PROXY_OIDC_USERINFO_CACHE_TTL"},
			Destination: &cfg.OIDC.UserinfoCache.TTL,
		},
		&cli.IntFlag{
			Name:        "oidc-userinfo-cache-size",
			Value:       cfg.OIDC.UserinfoCache.Size,
			Usage:       "Max entries for caching userinfo",
			EnvVars:     []string{"PROXY_OIDC_USERINFO_CACHE_SIZE"},
			Destination: &cfg.OIDC.UserinfoCache.Size,
		},

		// account related config

		&cli.BoolFlag{
			Name:        "autoprovision-accounts",
			Value:       cfg.AutoprovisionAccounts,
			Usage:       "create accounts from OIDC access tokens to learn new users",
			EnvVars:     []string{"PROXY_AUTOPROVISION_ACCOUNTS"},
			Destination: &cfg.AutoprovisionAccounts,
		},
		&cli.StringFlag{
			Name:        "user-oidc-claim",
			Value:       flags.OverrideDefaultString(cfg.UserOIDCClaim, "email"),
			Usage:       "The OIDC claim that is used to identify users, eg. 'ownclouduuid', 'uid', 'cn' or 'email'",
			EnvVars:     []string{"PROXY_USER_OIDC_CLAIM"},
			Destination: &cfg.UserOIDCClaim,
		},
		&cli.StringFlag{
			Name:        "user-cs3-claim",
			Value:       flags.OverrideDefaultString(cfg.UserCS3Claim, "mail"),
			Usage:       "The CS3 claim to use when looking up a user in the CS3 users API, eg. 'userid', 'username' or 'mail'",
			EnvVars:     []string{"PROXY_USER_CS3_CLAIM"},
			Destination: &cfg.UserCS3Claim,
		},

		// Pre Signed URLs
		&cli.StringSliceFlag{
			Name:    "presignedurl-allow-method",
			Value:   cli.NewStringSlice(cfg.PreSignedURL.AllowedHTTPMethods...),
			Usage:   "--presignedurl-allow-method GET [--presignedurl-allow-method POST]",
			EnvVars: []string{"PRESIGNEDURL_ALLOWED_METHODS"},
		},
		&cli.BoolFlag{
			Name:        "enable-presignedurls",
			Value:       cfg.PreSignedURL.Enabled,
			Usage:       "Enable or disable handling the presigned urls in the proxy",
			EnvVars:     []string{"PROXY_ENABLE_PRESIGNEDURLS"},
			Destination: &cfg.PreSignedURL.Enabled,
		},

		// Basic auth
		&cli.BoolFlag{
			Name:        "enable-basic-auth",
			Value:       cfg.EnableBasicAuth,
			Usage:       "enable basic authentication",
			EnvVars:     []string{"PROXY_ENABLE_BASIC_AUTH"},
			Destination: &cfg.EnableBasicAuth,
		},

		&cli.StringFlag{
			Name:        "account-backend-type",
			Value:       cfg.AccountBackend,
			Usage:       "account-backend-type",
			EnvVars:     []string{"PROXY_ACCOUNT_BACKEND_TYPE"},
			Destination: &cfg.AccountBackend,
		},

		// Reva Middlewares Config
		&cli.StringSliceFlag{
			Name:    "proxy-user-agent-lock-in",
			Usage:   "--user-agent-whitelist-lock-in=mirall:basic,foo:bearer Given a tuple of [UserAgent:challenge] it locks a given user agent to the authentication challenge. Particularly useful for old clients whose USer-Agent is known and only support one authentication challenge. When this flag is set in the proxy it configures the authentication middlewares.",
			EnvVars: []string{"PROXY_MIDDLEWARE_AUTH_CREDENTIALS_BY_USER_AGENT"},
		},
	}
}

// ListProxyWithConfig applies the config to the list commands flags.
func ListProxyWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "service-namespace",
			Value:       cfg.Service.Namespace,
			Usage:       "Set the base namespace for the service namespace",
			EnvVars:     []string{"PROXY_SERVICE_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "service-name",
			Value:       cfg.Service.Name,
			Usage:       "Service name",
			EnvVars:     []string{"PROXY_SERVICE_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
