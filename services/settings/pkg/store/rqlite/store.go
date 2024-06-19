package rqlite

import (
	"database/sql"

	olog "github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config"
	"github.com/owncloud/ocis/v2/services/settings/pkg/settings"
	metadata "github.com/owncloud/ocis/v2/services/settings/pkg/store/metadata"
	_ "github.com/rqlite/gorqlite/stdlib"
)

var (
	managerName = "rqlite-metadata"
)

func init() {
	settings.Registry[managerName] = New
}

type Store struct {
	mdStore settings.Manager
	db      *sql.DB
	logger  olog.Logger
}

// New creates a new store
func New(cfg *config.Config) settings.Manager {
	mdStore := metadata.New(cfg)
	db, err := sql.Open("rqlite", "http://rqlite:4001/")
	if err != nil {
		panic(err)
	}
	err = BootstrapDB(cfg, db)
	if err != nil {
		panic(err)

	}
	s := Store{
		mdStore: mdStore,
		db:      db,
		logger: olog.NewLogger(
			olog.Color(cfg.Log.Color),
			olog.Pretty(cfg.Log.Pretty),
			olog.Level(cfg.Log.Level),
			olog.File(cfg.Log.File),
		),
	}
	return &s
}

func (s *Store) ListBundles(bundleType settingsmsg.Bundle_Type, bundleIDs []string) ([]*settingsmsg.Bundle, error) {
	return s.mdStore.ListBundles(bundleType, bundleIDs)
}

func (s *Store) ReadBundle(bundleID string) (*settingsmsg.Bundle, error) {
	return s.mdStore.ReadBundle(bundleID)
}

func (s *Store) WriteBundle(bundle *settingsmsg.Bundle) (*settingsmsg.Bundle, error) {
	return s.mdStore.WriteBundle(bundle)
}

func (s *Store) ReadSetting(settingID string) (*settingsmsg.Setting, error) {
	return s.mdStore.ReadSetting(settingID)
}

func (s *Store) AddSettingToBundle(bundleID string, setting *settingsmsg.Setting) (*settingsmsg.Setting, error) {
	return s.mdStore.AddSettingToBundle(bundleID, setting)
}

func (s *Store) RemoveSettingFromBundle(bundleID, settingID string) error {
	return s.mdStore.RemoveSettingFromBundle(bundleID, settingID)
}

func (s *Store) ListValues(bundleID, accountUUID string) ([]*settingsmsg.Value, error) {
	return s.mdStore.ListValues(bundleID, accountUUID)
}

func (s *Store) ReadValue(valueID string) (*settingsmsg.Value, error) {
	return s.mdStore.ReadValue(valueID)
}

func (s *Store) ReadValueByUniqueIdentifiers(accountUUID, settingID string) (*settingsmsg.Value, error) {
	return s.mdStore.ReadValueByUniqueIdentifiers(accountUUID, settingID)
}

func (s *Store) WriteValue(value *settingsmsg.Value) (*settingsmsg.Value, error) {
	return s.mdStore.WriteValue(value)
}

func (s *Store) ListPermissionsByResource(resource *settingsmsg.Resource, roleIDs []string) ([]*settingsmsg.Permission, error) {
	return s.mdStore.ListPermissionsByResource(resource, roleIDs)
}

func (s *Store) ReadPermissionByID(permissionID string, roleIDs []string) (*settingsmsg.Permission, error) {
	return s.mdStore.ReadPermissionByID(permissionID, roleIDs)
}

func (s *Store) ReadPermissionByName(name string, roleIDs []string) (*settingsmsg.Permission, error) {
	return s.mdStore.ReadPermissionByName(name, roleIDs)
}
