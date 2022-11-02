package revaconfig

import (
	"github.com/owncloud/ocis/v2/services/sharing/pkg/config"
)

// SharingConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func SharingConfigFromStruct(cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
			"grpc_client_options":       cfg.Reva.GetGRPCClientConfig(),
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			"tls_settings": map[string]interface{}{
				"enabled":     cfg.GRPC.TLS.Enabled,
				"certificate": cfg.GRPC.TLS.Cert,
				"key":         cfg.GRPC.TLS.Key,
			},
			// TODO build services dynamically
			"services": map[string]interface{}{
				"usershareprovider": map[string]interface{}{
					"driver": cfg.UserSharingDriver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file":         cfg.UserSharingDrivers.JSON.File,
							"gateway_addr": cfg.Reva.Address,
						},
						"sql": map[string]interface{}{ // cernbox sql
							"db_username":                   cfg.UserSharingDrivers.SQL.DBUsername,
							"db_password":                   cfg.UserSharingDrivers.SQL.DBPassword,
							"db_host":                       cfg.UserSharingDrivers.SQL.DBHost,
							"db_port":                       cfg.UserSharingDrivers.SQL.DBPort,
							"db_name":                       cfg.UserSharingDrivers.SQL.DBName,
							"password_hash_cost":            cfg.UserSharingDrivers.SQL.PasswordHashCost,
							"enable_expired_shares_cleanup": cfg.UserSharingDrivers.SQL.EnableExpiredSharesCleanup,
							"janitor_run_interval":          cfg.UserSharingDrivers.SQL.JanitorRunInterval,
						},
						"owncloudsql": map[string]interface{}{
							"gateway_addr":     cfg.Reva.Address,
							"storage_mount_id": cfg.UserSharingDrivers.OwnCloudSQL.UserStorageMountID,
							"db_username":      cfg.UserSharingDrivers.OwnCloudSQL.DBUsername,
							"db_password":      cfg.UserSharingDrivers.OwnCloudSQL.DBPassword,
							"db_host":          cfg.UserSharingDrivers.OwnCloudSQL.DBHost,
							"db_port":          cfg.UserSharingDrivers.OwnCloudSQL.DBPort,
							"db_name":          cfg.UserSharingDrivers.OwnCloudSQL.DBName,
						},
						"cs3": map[string]interface{}{
							"gateway_addr":        cfg.UserSharingDrivers.CS3.ProviderAddr,
							"provider_addr":       cfg.UserSharingDrivers.CS3.ProviderAddr,
							"service_user_id":     cfg.UserSharingDrivers.CS3.SystemUserID,
							"service_user_idp":    cfg.UserSharingDrivers.CS3.SystemUserIDP,
							"machine_auth_apikey": cfg.UserSharingDrivers.CS3.SystemUserAPIKey,
						},
						"jsoncs3": map[string]interface{}{
							"gateway_addr":        cfg.Reva.Address,
							"provider_addr":       cfg.UserSharingDrivers.JSONCS3.ProviderAddr,
							"service_user_id":     cfg.UserSharingDrivers.JSONCS3.SystemUserID,
							"service_user_idp":    cfg.UserSharingDrivers.JSONCS3.SystemUserIDP,
							"machine_auth_apikey": cfg.UserSharingDrivers.JSONCS3.SystemUserAPIKey,
							"ttl":                 cfg.UserSharingDrivers.JSONCS3.CacheTTL,
						},
					},
				},
				"publicshareprovider": map[string]interface{}{
					"driver": cfg.PublicSharingDriver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"file":         cfg.PublicSharingDrivers.JSON.File,
							"gateway_addr": cfg.Reva.Address,
						},
						"sql": map[string]interface{}{
							"db_username":                   cfg.PublicSharingDrivers.SQL.DBUsername,
							"db_password":                   cfg.PublicSharingDrivers.SQL.DBPassword,
							"db_host":                       cfg.PublicSharingDrivers.SQL.DBHost,
							"db_port":                       cfg.PublicSharingDrivers.SQL.DBPort,
							"db_name":                       cfg.PublicSharingDrivers.SQL.DBName,
							"password_hash_cost":            cfg.PublicSharingDrivers.SQL.PasswordHashCost,
							"enable_expired_shares_cleanup": cfg.PublicSharingDrivers.SQL.EnableExpiredSharesCleanup,
							"janitor_run_interval":          cfg.PublicSharingDrivers.SQL.JanitorRunInterval,
						},
						"cs3": map[string]interface{}{
							"gateway_addr":        cfg.PublicSharingDrivers.CS3.ProviderAddr,
							"provider_addr":       cfg.PublicSharingDrivers.CS3.ProviderAddr,
							"service_user_id":     cfg.PublicSharingDrivers.CS3.SystemUserID,
							"service_user_idp":    cfg.PublicSharingDrivers.CS3.SystemUserIDP,
							"machine_auth_apikey": cfg.PublicSharingDrivers.CS3.SystemUserAPIKey,
						},
						"jsoncs3": map[string]interface{}{
							"gateway_addr":        cfg.Reva.Address,
							"provider_addr":       cfg.PublicSharingDrivers.JSONCS3.ProviderAddr,
							"service_user_id":     cfg.PublicSharingDrivers.JSONCS3.SystemUserID,
							"service_user_idp":    cfg.PublicSharingDrivers.JSONCS3.SystemUserIDP,
							"machine_auth_apikey": cfg.PublicSharingDrivers.JSONCS3.SystemUserAPIKey,
						},
					},
				},
			},
			"interceptors": map[string]interface{}{
				"eventsmiddleware": map[string]interface{}{
					"group":            "sharing",
					"type":             "nats",
					"address":          cfg.Events.Addr,
					"clusterID":        cfg.Events.ClusterID,
					"tls-insecure":     cfg.Events.TLSInsecure,
					"tls-root-ca-cert": cfg.Events.TLSRootCaCertPath,
					"enable-tls":       cfg.Events.EnableTLS,
				},
				"prometheus": map[string]interface{}{
					"namespace": "ocis",
					"subsystem": "sharing",
				},
			},
		},
	}
	return rcfg
}
