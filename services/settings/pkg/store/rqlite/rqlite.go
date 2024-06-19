package rqlite

import (
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
)

func BootstrapDB(cfg *config.Config, db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS roleAssignments (
			id TEXT NOT NULL PRIMARY KEY,
			userid TEXT NOT NULL UNIQUE,
			roleid TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_roleAssignments_roleid ON roleAssignments(roleid);
		`)
	if err != nil {
		panic(err)
	}

	for _, assignment := range defaults.DefaultRoleAssignments(cfg) {
		assignment.Id = uuid.Must(uuid.NewV4()).String()
		_, err := db.Exec(`
			INSERT OR IGNORE INTO roleAssignments (id, userid, roleid) VALUES (?, ?, ?);
			`, assignment.Id, assignment.AccountUuid, assignment.RoleId)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
