package defaults

import (
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/auth-basic/pkg/config"
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
			Addr:   "127.0.0.1:9147",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:      "127.0.0.1:9146",
			Namespace: "com.owncloud.api",
			Protocol:  "tcp",
		},
		Service: config.Service{
			Name: "auth-basic",
		},
		Reva:         shared.DefaultRevaConfig(),
		AuthProvider: "ldap",
		AuthProviders: config.AuthProviders{
			LDAP: config.LDAPProvider{
				URI:              "ldaps://localhost:9235",
				CACert:           filepath.Join(defaults.BaseDataPath(), "idm", "ldap.crt"),
				Insecure:         false,
				UserBaseDN:       "ou=users,o=libregraph-idm",
				GroupBaseDN:      "ou=groups,o=libregraph-idm",
				UserScope:        "sub",
				GroupScope:       "sub",
				LoginAttributes:  []string{"uid"},
				UserFilter:       "",
				GroupFilter:      "",
				UserObjectClass:  "inetOrgPerson",
				GroupObjectClass: "groupOfNames",
				BindDN:           "uid=reva,ou=sysusers,o=libregraph-idm",
				IDP:              "https://localhost:9200",
				UserSchema: config.LDAPUserSchema{
					ID:          "ownclouduuid",
					Mail:        "mail",
					DisplayName: "displayname",
					Username:    "uid",
				},
				GroupSchema: config.LDAPGroupSchema{
					ID:          "ownclouduuid",
					Mail:        "mail",
					DisplayName: "cn",
					Groupname:   "cn",
					Member:      "member",
				},
			},
			JSON: config.JSONProvider{},
			OwnCloudSQL: config.OwnCloudSQLProvider{
				DBUsername:       "owncloud",
				DBHost:           "mysql",
				DBPort:           3306,
				DBName:           "owncloud",
				IDP:              "https://localhost:9200",
				Nobody:           90,
				JoinUsername:     false,
				JoinOwnCloudUUID: false,
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

	if cfg.GRPC.TLS == nil {
		cfg.GRPC.TLS = &shared.GRPCServiceTLS{}
		if cfg.Commons != nil && cfg.Commons.GRPCServiceTLS != nil {
			cfg.GRPC.TLS.Enabled = cfg.Commons.GRPCServiceTLS.Enabled
			cfg.GRPC.TLS.Cert = cfg.Commons.GRPCServiceTLS.Cert
			cfg.GRPC.TLS.Key = cfg.Commons.GRPCServiceTLS.Key
		}
	}
}

func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
