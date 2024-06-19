package rqlite

import (
	"github.com/gofrs/uuid"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
)

func (s *Store) ListRoleAssignments(accountUUID string) ([]*settingsmsg.UserRoleAssignment, error) {
	query := `
		SELECT id, userid, roleid FROM roleAssignments WHERE userid = ?;
	`
	rows, err := s.db.Query(query, accountUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []*settingsmsg.UserRoleAssignment{}
	for rows.Next() {
		assignment := settingsmsg.UserRoleAssignment{}
		rows.Scan(&assignment.Id, &assignment.AccountUuid, &assignment.RoleId)
		s.logger.Debug().Str("id", assignment.Id).Str("accountUuid", assignment.AccountUuid).Str("roleId", assignment.RoleId).Msg("role assignment")
		assignments = append(assignments, &assignment)
	}
	return assignments, nil
}

func (s *Store) ListRoleAssignmentsByRole(roleID string) ([]*settingsmsg.UserRoleAssignment, error) {
	query := `
		SELECT id, userid, roleid FROM roleAssignments WHERE roleid = ?;
	`
	rows, err := s.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignments := []*settingsmsg.UserRoleAssignment{}
	for rows.Next() {
		assignment := settingsmsg.UserRoleAssignment{}
		rows.Scan(&assignment.Id, &assignment.AccountUuid, &assignment.RoleId)
		s.logger.Debug().Str("id", assignment.Id).Str("accountUuid", assignment.AccountUuid).Str("roleId", assignment.RoleId).Msg("role assignment")
		assignments = append(assignments, &assignment)
	}
	return assignments, nil
}

func (s *Store) WriteRoleAssignment(accountUUID, roleID string) (*settingsmsg.UserRoleAssignment, error) {
	query := `
		INSERT OR REPLACE INTO roleAssignments (id, userid, roleid) VALUES (?, ?, ?);
	`
	id := uuid.Must(uuid.NewV4()).String()

	_, err := s.db.Exec(query, id, accountUUID, roleID)
	if err != nil {
		return nil, err
	}
	assignment := &settingsmsg.UserRoleAssignment{
		Id:          id,
		AccountUuid: accountUUID,
		RoleId:      roleID,
	}
	return assignment, nil
}

func (s *Store) RemoveRoleAssignment(assignmentID string) error {
	return s.mdStore.RemoveRoleAssignment(assignmentID)
}
