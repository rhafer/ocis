package init

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/gofrs/uuid"
	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"gopkg.in/yaml.v2"
)

const (
	configFilename = "ocis.yaml" // TODO: use also a constant for reading this file
	passwordLength = 32
)

var (
	_insecureService = InsecureService{Insecure: true}
	_insecureEvents  = Events{TLSInsecure: true}
)

type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret"`
}

type InsecureService struct {
	Insecure bool
}

type InsecureProxyService struct {
	OIDC             InsecureProxyOIDC `yaml:"oidc"`
	InsecureBackends bool              `yaml:"insecure_backends"`
}

type InsecureProxyOIDC struct {
	Insecure bool `yaml:"insecure"`
}

type LdapSettings struct {
	BindPassword string `yaml:"bind_password"`
}
type LdapBasedService struct {
	Ldap LdapSettings
}

type Events struct {
	TLSInsecure bool `yaml:"tls_insecure"`
}
type GraphApplication struct {
	ID string `yaml:"id"`
}

type GraphService struct {
	Application GraphApplication
	Events      Events
	Spaces      InsecureService
	Identity    LdapBasedService
}

type ServiceUserPasswordsSettings struct {
	AdminPassword string `yaml:"admin_password"`
	IdmPassword   string `yaml:"idm_password"`
	RevaPassword  string `yaml:"reva_password"`
	IdpPassword   string `yaml:"idp_password"`
}
type IdmService struct {
	ServiceUserPasswords ServiceUserPasswordsSettings `yaml:"service_user_passwords"`
}

type FrontendService struct {
	Archiver InsecureService
}

type AuthbasicService struct {
	AuthProviders LdapBasedService `yaml:"auth_providers"`
}

type AuthProviderSettings struct {
	Oidc InsecureService
}
type AuthbearerService struct {
	AuthProviders AuthProviderSettings `yaml:"auth_providers"`
}

type UsersAndGroupsService struct {
	Drivers LdapBasedService
}

type ThumbnailSettings struct {
	TransferSecret      string `yaml:"transfer_secret"`
	WebdavAllowInsecure bool   `yaml:"webdav_allow_insecure"`
	Cs3AllowInsecure    bool   `yaml:"cs3_allow_insecure"`
}

type ThumbnailService struct {
	Thumbnail ThumbnailSettings
}

type Search struct {
	Events Events
}

type Audit struct {
	Events Events
}

type Sharing struct {
	Events Events
}

type StorageUsers struct {
	Events  Events
	MountID string `yaml:"mount_id"`
}

type Gateway struct {
	StorageRegistry StorageRegistry `yaml:"storage_registry"`
}

type StorageRegistry struct {
	StorageUsersMountID string `yaml:"storage_users_mount_id"`
}

type Notifications struct {
	Notifications struct{ Events Events } // The notifications config has a field called notifications
}

type Nats struct {
	// The nats config has a field called nats
	Nats struct {
		TLSSkipVerifyClientCert bool `yaml:"tls_skip_verify_client_cert"`
	}
}

// TODO: use the oCIS config struct instead of this custom struct
// We can't use it right now, because it would need  "omitempty" on
// all elements, in order to produce a slim config file with `ocis init`.
// We can't just add these "omitempty" tags, since we want to generate
// full example configuration files with that struct, too.
// Proposed solution to  get rid of this temporary solution:
// - use the oCIS config struct
// - set the needed values like below
// - marshal it to yaml
// - unmarshal it into yaml.Node
// - recurse through the nodes and delete empty / default ones
// - marshal it to yaml
type OcisConfig struct {
	TokenManager      TokenManager `yaml:"token_manager"`
	MachineAuthAPIKey string       `yaml:"machine_auth_api_key"`
	SystemUserAPIKey  string       `yaml:"system_user_api_key"`
	TransferSecret    string       `yaml:"transfer_secret"`
	SystemUserID      string       `yaml:"system_user_id"`
	AdminUserID       string       `yaml:"admin_user_id"`
	Graph             GraphService
	Idp               LdapBasedService
	Authelia          LdapBasedService
	Idm               IdmService
	Proxy             InsecureProxyService
	Frontend          FrontendService
	AuthBasic         AuthbasicService  `yaml:"auth_basic"`
	AuthBearer        AuthbearerService `yaml:"auth_bearer"`
	Users             UsersAndGroupsService
	Groups            UsersAndGroupsService
	Ocdav             InsecureService
	Thumbnails        ThumbnailService
	Search            Search
	Audit             Audit
	Sharing           Sharing
	StorageUsers      StorageUsers `yaml:"storage_users"`
	Notifications     Notifications
	Nats              Nats
	Gateway           Gateway
}

func checkConfigPath(configPath string) error {
	targetPath := path.Join(configPath, configFilename)
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("config in %s already exists", targetPath)
	}
	return nil
}

func backupOcisConfigFile(configPath string) (string, error) {
	sourceConfig := path.Join(configPath, configFilename)
	targetBackupConfig := path.Join(configPath, configFilename+"."+time.Now().Format("2006-01-02-15-04-05")+".backup")
	source, err := os.Open(sourceConfig)
	if err != nil {
		log.Fatalf("Could not read %s (%s)", sourceConfig, err)
	}
	defer source.Close()
	target, err := os.Create(targetBackupConfig)
	if err != nil {
		log.Fatalf("Could not generate backup %s (%s)", targetBackupConfig, err)
	}
	defer target.Close()
	_, err = io.Copy(target, source)
	if err != nil {
		log.Fatalf("Could not write backup %s (%s)", targetBackupConfig, err)
	}
	return targetBackupConfig, nil
}

// CreateConfig creates a config file with random passwords at configPath
func CreateConfig(insecure, forceOverwrite bool, configPath, adminPassword string) error {
	err := checkConfigPath(configPath)
	if err != nil && !forceOverwrite {
		return err
	}
	targetBackupConfig := ""
	if err != nil {
		targetBackupConfig, err = backupOcisConfigFile(configPath)
		if err != nil {
			return err
		}
	}
	err = os.MkdirAll(configPath, 0700)
	if err != nil {
		return err
	}

	systemUserID := uuid.Must(uuid.NewV4()).String()
	adminUserID := uuid.Must(uuid.NewV4()).String()
	graphApplicationID := uuid.Must(uuid.NewV4()).String()
	storageUsersMountID := uuid.Must(uuid.NewV4()).String()

	idmServicePassword, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for idm: %s", err)
	}
	idpServicePassword, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for idp: %s", err)
	}
	ocisAdminServicePassword := adminPassword
	if ocisAdminServicePassword == "" {
		ocisAdminServicePassword, err = generators.GenerateRandomPassword(passwordLength)
		if err != nil {
			return fmt.Errorf("could not generate random password for ocis admin: %s", err)
		}
	}

	revaServicePassword, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for reva: %s", err)
	}
	tokenManagerJwtSecret, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for tokenmanager: %s", err)
	}
	machineAuthAPIKey, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for machineauthsecret: %s", err)
	}
	systemUserAPIKey, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random system user API key: %s", err)
	}
	revaTransferSecret, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for revaTransferSecret: %s", err)
	}
	thumbnailsTransferSecret, err := generators.GenerateRandomPassword(passwordLength)
	if err != nil {
		return fmt.Errorf("could not generate random password for thumbnailsTransferSecret: %s", err)
	}

	cfg := OcisConfig{
		TokenManager: TokenManager{
			JWTSecret: tokenManagerJwtSecret,
		},
		MachineAuthAPIKey: machineAuthAPIKey,
		SystemUserAPIKey:  systemUserAPIKey,
		TransferSecret:    revaTransferSecret,
		SystemUserID:      systemUserID,
		AdminUserID:       adminUserID,
		Idm: IdmService{
			ServiceUserPasswords: ServiceUserPasswordsSettings{
				AdminPassword: ocisAdminServicePassword,
				IdpPassword:   idpServicePassword,
				RevaPassword:  revaServicePassword,
				IdmPassword:   idmServicePassword,
			},
		},
		Idp: LdapBasedService{
			Ldap: LdapSettings{
				BindPassword: idpServicePassword,
			},
		},
		Authelia: LdapBasedService{
			Ldap: LdapSettings{
				BindPassword: idpServicePassword,
			},
		},
		AuthBasic: AuthbasicService{
			AuthProviders: LdapBasedService{
				Ldap: LdapSettings{
					BindPassword: revaServicePassword,
				},
			},
		},
		Groups: UsersAndGroupsService{
			Drivers: LdapBasedService{
				Ldap: LdapSettings{
					BindPassword: revaServicePassword,
				},
			},
		},
		Users: UsersAndGroupsService{
			Drivers: LdapBasedService{
				Ldap: LdapSettings{
					BindPassword: revaServicePassword,
				},
			},
		},
		Graph: GraphService{
			Application: GraphApplication{
				ID: graphApplicationID,
			},
			Identity: LdapBasedService{
				Ldap: LdapSettings{
					BindPassword: idmServicePassword,
				},
			},
		},
		Thumbnails: ThumbnailService{
			Thumbnail: ThumbnailSettings{
				TransferSecret: thumbnailsTransferSecret,
			},
		},
		Gateway: Gateway{
			StorageRegistry: StorageRegistry{
				StorageUsersMountID: storageUsersMountID,
			},
		},
		StorageUsers: StorageUsers{
			MountID: storageUsersMountID,
		},
	}

	if insecure {

		cfg.AuthBearer = AuthbearerService{
			AuthProviders: AuthProviderSettings{Oidc: _insecureService},
		}
		cfg.Frontend = FrontendService{Archiver: _insecureService}
		cfg.Graph.Spaces = _insecureService
		cfg.Graph.Events = _insecureEvents
		cfg.Notifications.Notifications.Events = _insecureEvents
		cfg.Search.Events = _insecureEvents
		cfg.Audit.Events = _insecureEvents
		cfg.Sharing.Events = _insecureEvents
		cfg.StorageUsers.Events = _insecureEvents
		cfg.Nats.Nats.TLSSkipVerifyClientCert = true
		cfg.Ocdav = _insecureService
		cfg.Proxy = InsecureProxyService{
			InsecureBackends: true,
			OIDC: InsecureProxyOIDC{
				Insecure: true,
			},
		}

		cfg.Thumbnails.Thumbnail.WebdavAllowInsecure = true
		cfg.Thumbnails.Thumbnail.Cs3AllowInsecure = true
	}

	yamlOutput, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not marshall config into yaml: %s", err)
	}
	targetPath := path.Join(configPath, configFilename)
	err = os.WriteFile(targetPath, yamlOutput, 0600)
	if err != nil {
		return err
	}
	fmt.Printf(
		"\n=========================================\n"+
			" generated OCIS Config\n"+
			"=========================================\n"+
			" configpath : %s\n"+
			" user       : admin\n"+
			" password   : %s\n\n",
		targetPath, ocisAdminServicePassword)
	if targetBackupConfig != "" {
		fmt.Printf("\n=========================================\n"+
			"An older config file has been backuped to\n %s\n\n",
			targetBackupConfig)
	}
	return nil
}
