package svc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"

	"github.com/CiscoM31/godata"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"

	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settings "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	ocissettingssvc "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
)

// GetMe implements the Service interface.
func (g Graph) GetMe(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Debug().Msg("calling get user in /me")
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")

	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		logger.Debug().Msg("could not get user: user not in context")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "user not in context")
		return
	}

	exp, err := identity.GetExpandValues(odataReq.Query)
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: $expand error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var me *libregraph.User
	// We can just return the user from context unless we need to expand the group memberships
	if !slices.Contains(exp, "memberOf") {
		me = identity.CreateUserModelFromCS3(u)
	} else {
		var err error
		logger.Debug().Msg("calling get user on backend")
		me, err = g.identityBackend.GetUser(r.Context(), u.GetId().GetOpaqueId(), odataReq)
		if err != nil {
			logger.Debug().Err(err).Interface("user", u).Msg("could not get user from backend")
			errorcode.RenderError(w, r, err)
			return
		}
		if me.MemberOf == nil {
			me.MemberOf = []libregraph.Group{}
		}
	}

	// expand appRoleAssignments if requested
	if slices.Contains(exp, "appRoleAssignments") {
		var err error
		me.AppRoleAssignments, err = g.fetchAppRoleAssignments(r.Context(), me.GetId())
		if err != nil {
			logger.Debug().Err(err).Str("userid", me.GetId()).Msg("could not get appRoleAssignments for self")
			errorcode.RenderError(w, r, err)
			return
		}
	}

	preferedLanguage, _, err := getUserLanguage(r.Context(), g.valueService, me.GetId())
	if err != nil {
		logger.Error().Err(err).Msg("could not get user language")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "could not get user language")
		return
	}

	me.PreferredLanguage = &preferedLanguage

	render.Status(r, http.StatusOK)
	render.JSON(w, r, me)
}

func (g Graph) fetchAppRoleAssignments(ctx context.Context, accountuuid string) ([]libregraph.AppRoleAssignment, error) {
	lrar, err := g.roleService.ListRoleAssignments(ctx, &settings.ListRoleAssignmentsRequest{
		AccountUuid: accountuuid,
	})
	if err != nil {
		return []libregraph.AppRoleAssignment{}, err
	}

	values := make([]libregraph.AppRoleAssignment, 0, len(lrar.Assignments))
	for _, assignment := range lrar.GetAssignments() {
		values = append(values, g.assignmentToAppRoleAssignment(assignment))
	}
	return values, nil
}

// GetUserDrive implements the Service interface.
func (g Graph) GetUserDrive(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Debug().Interface("query", r.URL.Query()).Msg("calling get user drive")

	userID, err := url.PathUnescape(chi.URLParam(r, "userID"))
	if err != nil {
		logger.Debug().Err(err).Str("userID", chi.URLParam(r, "userID")).Msg("could not get drive: unescaping drive id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping user id failed")
		return
	}

	if userID == "" {
		u, ok := revactx.ContextGetUser(r.Context())
		if !ok {
			logger.Debug().Msg("could not get user: user not in context")
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "user not in context")
			return
		}
		userID = u.GetId().GetOpaqueId()
	}

	logger.Debug().Str("userID", userID).Msg("calling list storage spaces with user and personal filter")
	ctx := r.Context()

	filters := []*storageprovider.ListStorageSpacesRequest_Filter{listStorageSpacesTypeFilter("personal"), listStorageSpacesUserFilter(userID)}
	res, err := g.ListStorageSpacesWithFilters(ctx, filters, true)
	switch {
	case err != nil:
		logger.Error().Err(err).Msg("could not get drive: transport error")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			// the client is doing a lookup for a specific space, therefore we need to return
			// not found to the caller
			logger.Debug().Str("userID", userID).Msg("could not get personal drive for user: not found")
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, "drive not found")
			return
		}
		logger.Debug().
			Str("userID", userID).
			Str("grpcmessage", res.GetStatus().GetMessage()).
			Msg("could not get personal drive for user: grpc error")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	webDavBaseURL, err := g.getWebDavBaseURL()
	if err != nil {
		logger.Error().Err(err).Str("url", webDavBaseURL.String()).Msg("could not get personal drive: error parsing webdav base url")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	spaces, err := g.formatDrives(ctx, webDavBaseURL, res.StorageSpaces)
	if err != nil {
		logger.Debug().Err(err).Msg("could not get personal drive: error parsing grpc response")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	switch num := len(spaces); {
	case num == 0:
		logger.Debug().Str("userID", userID).Msg("could not get personal drive: no drive returned from storage")
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, "no drive returned from storage")
		return
	case num == 1:
		render.Status(r, http.StatusOK)
		render.JSON(w, r, spaces[0])
	default:
		logger.Debug().Int("number", num).Msg("could not get personal drive: expected to find a single drive but fetched more")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "could not get personal drive: expected to find a single drive but fetched more")
		return
	}
}

func (g Graph) contextUserHasFullAccountPerms(reqctx context.Context) bool {
	// mostly copied from the canCreateSpace method
	pr, err := g.permissionsService.GetPermissionByID(reqctx, &settingssvc.GetPermissionByIDRequest{
		PermissionId: defaults.AccountManagementPermission(0).Id,
	})
	if err != nil || pr.Permission == nil {
		return false
	}

	if pr.Permission.Constraint != defaults.All {
		return false
	}
	return true
}

// GetUsers implements the Service interface.
func (g Graph) GetUsers(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Debug().Interface("query", r.URL.Query()).Msg("calling get users")
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctxHasFullPerms := g.contextUserHasFullAccountPerms(r.Context())
	if !ctxHasFullPerms && (odataReq.Query == nil || odataReq.Query.Search == nil || len(odataReq.Query.Search.RawValue) < 3) {
		// regular user must search with at least 3 chars
		logger.Debug().Interface("query", r.URL.Query()).Msg("search with less than 3 chars for a regular user")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "regular users must enter at least 3 characters to search")
		return
	}

	if !ctxHasFullPerms && (odataReq.Query.Filter != nil || odataReq.Query.Apply != nil || odataReq.Query.Expand != nil || odataReq.Query.Compute != nil) {
		// regular users can't use filter, apply, expand and compute
		logger.Debug().Interface("query", r.URL.Query()).Msg("forbidden query elements for a regular user")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "query has forbidden elements for regular users")
		return
	}

	logger.Debug().Interface("query", r.URL.Query()).Msg("calling get users on backend")

	var users []*libregraph.User

	if odataReq.Query.Filter != nil {
		users, err = g.applyUserFilter(r.Context(), odataReq, nil)
	} else {
		users, err = g.identityBackend.GetUsers(r.Context(), odataReq)
	}

	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users from backend")
		var errcode errorcode.Error
		var godataerr *godata.GoDataError
		switch {
		case errors.As(err, &errcode):
			errcode.Render(w, r)
		case errors.As(err, &godataerr):
			errorcode.GeneralException.Render(w, r, godataerr.ResponseCode, err.Error())
		default:
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// If the user isn't admin, we'll show just the minimum user attibutes
	if !ctxHasFullPerms {
		finalUsers := make([]*libregraph.User, len(users))
		for i, u := range users {
			finalUsers[i] = &libregraph.User{
				Id:          u.Id,
				DisplayName: u.DisplayName,
				Mail:        u.Mail,
				UserType:    u.UserType,
			}
		}
		users = finalUsers
	}

	exp, err := identity.GetExpandValues(odataReq.Query)
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: $expand error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	expandAppRoleAssignments := slices.Contains(exp, "appRoleAssignments")
	expandMemberOf := slices.Contains(exp, "memberOf")
	for _, u := range users {
		if expandAppRoleAssignments {
			u.AppRoleAssignments, err = g.fetchAppRoleAssignments(r.Context(), u.GetId())
			if err != nil {
				// TODO I think we should not continue here, see http://docs.oasis-open.org/odata/odata/v4.01/odata-v4.01-part1-protocol.html#sec_SystemQueryOptionexpand
				// > The $expand system query option indicates the related entities and stream values that MUST be represented inline. The service MUST return the specified content, and MAY choose to return additional information.
				logger.Debug().Err(err).Str("userid", u.GetId()).Msg("could not get appRoleAssignments when listing user")
			}
		}
		if expandMemberOf {
			if u.MemberOf == nil {
				u.MemberOf = []libregraph.Group{}
			}
		}
	}

	users, err = sortUsers(odataReq, users)
	if err != nil {
		logger.Debug().Interface("query", odataReq).Msg("error while sorting users according to query")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: users})
}

// PostUser implements the Service interface.
func (g Graph) PostUser(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Debug().Msg("calling create user")
	u := libregraph.NewUser()
	err := StrictJSONUnmarshal(r.Body, u)
	if err != nil {
		logger.Info().Err(err).Msg("could not create user: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err.Error()))
		return
	}

	if _, ok := u.GetDisplayNameOk(); !ok {
		logger.Info().Err(err).Interface("user", u).Msg("could not create user: missing required Attribute: 'displayName'")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing required Attribute: 'displayName'")
		return
	}
	if accountName, ok := u.GetOnPremisesSamAccountNameOk(); ok {
		if !g.isValidUsername(*accountName) {
			logger.Info().Str("username", *accountName).Msg("could not create user: invalid username")
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Invalid username")
			return
		}
	} else {
		logger.Info().Interface("user", u).Msg("could not create user: missing required Attribute: 'onPremisesSamAccountName'")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing required Attribute: 'onPremisesSamAccountName'")
		return
	}

	if mail, ok := u.GetMailOk(); ok {
		if !isValidEmail(*mail) {
			logger.Info().Str("mail", *u.Mail).Msg("could not create user: invalid email address")
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid email address")
			return
		}
	}

	// Disallow user-supplied IDs. It's supposed to be readonly. We're either
	// generating them in the backend ourselves or rely on the Backend's
	// storage (e.g. LDAP) to provide a unique ID.
	if _, ok := u.GetIdOk(); ok {
		logger.Info().Interface("user", u).Msg("could not create user: user id is a read-only attribute")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "user id is a read-only attribute")
		return
	}

	if u.HasUserType() {
		if !isValidUserType(*u.UserType) {
			logger.Info().Interface("user", u).Msg("invalid userType attribute")
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid userType attribute, valid options are 'Member' or 'Guest'")
			return
		}
	} else {
		u.SetUserType("Member")
	}

	logger.Debug().Interface("user", u).Msg("calling create user on backend")
	if u, err = g.identityBackend.CreateUser(r.Context(), *u); err != nil {
		logger.Error().Err(err).Msg("could not create user: backend error")
		errorcode.RenderError(w, r, err)
		return
	}

	// assign roles if possible
	if g.roleService != nil && g.config.API.AssignDefaultUserRole {
		// All users get the user role by default currently.
		// to all new users for now, as create Account request does not have any role field
		if _, err = g.roleService.AssignRoleToUser(r.Context(), &settings.AssignRoleToUserRequest{
			AccountUuid: *u.Id,
			RoleId:      ocissettingssvc.BundleUUIDRoleUser,
		}); err != nil {
			// log as error, admin eventually needs to do something
			logger.Error().Err(err).Str("id", *u.Id).Str("role", ocissettingssvc.BundleUUIDRoleUser).Msg("could not create user: role assignment failed")
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "role assignment failed")
			return
		}
	}

	e := events.UserCreated{UserID: *u.Id}
	if currentUser, ok := revactx.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(r.Context(), e)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, u)
}

// GetUser implements the Service interface.
func (g Graph) GetUser(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Debug().Msg("calling get user")

	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")

	userID := chi.URLParam(r, "userID")
	userID, err := url.PathUnescape(userID)
	if err != nil {
		logger.Debug().Err(err).Str("id", userID).Msg("could not get user: unescaping user id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping user id failed")
		return
	}

	if userID == "" {
		logger.Debug().Msg("could not get user: missing user id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing user id")
		return
	}

	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	exp, err := identity.GetExpandValues(odataReq.Query)
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: $expand error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debug().Str("id", userID).Msg("calling get user from backend")
	user, err := g.identityBackend.GetUser(r.Context(), userID, odataReq)
	if err != nil {
		logger.Debug().Err(err).Msg("could not get user: error fetching user from backend")
		errorcode.RenderError(w, r, err)
		return
	}

	listDrives := slices.Contains(exp, "drives")
	listDrive := slices.Contains(exp, "drive")

	// do we need to list all or only the personal drive
	filters := []*storageprovider.ListStorageSpacesRequest_Filter{}
	if listDrives || listDrive {
		filters = append(filters, listStorageSpacesUserFilter(user.GetId()))
		if !listDrives {
			// TODO filter by owner when decomposedfs supports the OWNER filter
			filters = append(filters, listStorageSpacesTypeFilter("personal"))
		}
	}

	if len(filters) > 0 {
		wdu, err := g.getWebDavBaseURL()
		if err != nil {
			// log error, wrong configuration
			logger.Error().
				Err(err).
				Str("webdav_base", g.config.Spaces.WebDavBase).
				Str("webdav_path", g.config.Spaces.WebDavPath).
				Msg("error parsing webdav URL")
			render.Status(r, http.StatusInternalServerError)
			return
		}
		logger.Debug().Str("id", user.GetId()).Msg("calling list storage spaces with filter")

		// use the unrestricted flag to get all possible spaces
		// users with the canListAllSpaces permission should see all spaces

		client, err := g.gatewaySelector.Next()
		if err != nil {
			logger.Error().Err(err).Msg("error selecting next gateway client")
			render.Status(r, http.StatusInternalServerError)
			return
		}
		opaque := utils.AppendPlainToOpaque(nil, "unrestricted", "T")
		lspr, err := client.ListStorageSpaces(r.Context(), &storageprovider.ListStorageSpacesRequest{
			Opaque:  opaque,
			Filters: filters,
		})
		if err != nil {
			// transport error, needs to be fixed by admin
			logger.Error().Err(err).Interface("query", r.URL.Query()).Msg("error getting storages: transport error")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, user)
			return
		}
		if lspr.GetStatus().GetCode() != cs3rpc.Code_CODE_OK {
			logger.Debug().Str("grpc", lspr.GetStatus().GetMessage()).Msg("could not list drives for user")
			// in case of NOT_OK, we can just return the user object with empty drives
			render.Status(r, status.HTTPStatusFromCode(http.StatusOK))
			render.JSON(w, r, user)
			return
		}
		if listDrives {
			user.Drives = make([]libregraph.Drive, 0, len(lspr.GetStorageSpaces()))
		}
		if listDrive {
			user.Drive = &libregraph.Drive{}
		}
		for _, sp := range lspr.GetStorageSpaces() {
			d, err := g.cs3StorageSpaceToDrive(r.Context(), wdu, sp)
			if err != nil {
				logger.Debug().Err(err).Interface("id", sp.Id).Msg("error converting space to drive")
				continue
			}
			quota, err := g.getDriveQuota(r.Context(), sp)
			if err != nil {
				logger.Debug().Err(err).Interface("id", sp.Id).Msg("error calling get quota on drive")
			}
			d.Quota = &quota
			if d.GetDriveType() == "personal" && sp.GetOwner().GetId().GetOpaqueId() == user.GetId() {
				if listDrive {
					user.Drive = d
				}
			} else {
				if listDrives {
					user.Drives = append(user.Drives, *d)
				}
			}
		}
	}

	// expand appRoleAssignments if requested
	if slices.Contains(exp, "appRoleAssignments") {
		user.AppRoleAssignments, err = g.fetchAppRoleAssignments(r.Context(), user.GetId())
		if err != nil {
			logger.Debug().Err(err).Str("userid", user.GetId()).Msg("could not get appRoleAssignments for user")
			errorcode.RenderError(w, r, err)
			return
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

// getUserLanguage returns the language of the user in the context.
func getUserLanguage(ctx context.Context, valueService settingssvc.ValueService, userID string) (string, string, error) {
	gvr, err := valueService.GetValueByUniqueIdentifiers(ctx, &settingssvc.GetValueByUniqueIdentifiersRequest{
		AccountUuid: userID,
		SettingId:   defaults.SettingUUIDProfileLanguage,
	})
	if err != nil {
		return "", "", err
	}

	langVal := gvr.GetValue().GetValue().GetListValue().GetValues()
	if len(langVal) > 0 && langVal[0] != nil {
		return langVal[0].GetStringValue(), gvr.GetValue().GetValue().GetId(), nil
	}
	return "", "", errors.New("no language value found")
}

// DeleteUser implements the Service interface.
func (g Graph) DeleteUser(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Debug().Msg("calling delete user")
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
	userID := chi.URLParam(r, "userID")
	userID, err := url.PathUnescape(userID)
	if err != nil {
		logger.Debug().Err(err).Str("id", userID).Msg("could not delete user: unescaping user id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping user id failed")
		return
	}

	if userID == "" {
		logger.Debug().Msg("could not delete user: missing user id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing user id")
		return
	}

	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debug().Str("id", userID).Msg("calling get user on user backend")
	user, err := g.identityBackend.GetUser(r.Context(), userID, odataReq)
	if err != nil {
		logger.Debug().Err(err).Str("userID", userID).Msg("failed to get user from backend")
		errorcode.RenderError(w, r, err)
		return
	}

	e := events.UserDeleted{UserID: user.GetId()}
	if currentUser, ok := revactx.ContextGetUser(r.Context()); ok {
		if currentUser.GetId().GetOpaqueId() == user.GetId() {
			logger.Debug().Msg("could not delete user: self deletion forbidden")
			errorcode.NotAllowed.Render(w, r, http.StatusForbidden, "self deletion forbidden")
			return
		}
		e.Executant = currentUser.GetId()
	}

	if g.gatewaySelector != nil {
		logger.Debug().
			Str("user", user.GetId()).
			Msg("calling list spaces with user filter to fetch the personal space for deletion")
		opaque := utils.AppendPlainToOpaque(nil, "unrestricted", "T")
		f := listStorageSpacesUserFilter(user.GetId())
		client, err := g.gatewaySelector.Next()
		if err != nil {
			logger.Error().Err(err).Msg("error selecting next gateway client")
			errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, "error selecting next gateway client, aborting")
			return
		}
		lspr, err := client.ListStorageSpaces(r.Context(), &storageprovider.ListStorageSpacesRequest{
			Opaque:  opaque,
			Filters: []*storageprovider.ListStorageSpacesRequest_Filter{f},
		})
		if err != nil {
			// transport error, log as error
			logger.Error().Err(err).Msg("could not fetch spaces: transport error")
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "could not fetch spaces for deletion, aborting")
			return
		}
		for _, sp := range lspr.GetStorageSpaces() {
			if !(sp.SpaceType == "personal" && sp.Owner.Id.OpaqueId == user.GetId()) {
				continue
			}
			// TODO: check if request contains a homespace and if, check if requesting user has the privilege to
			// delete it and make sure it is not deleting its own homespace
			// needs modification of the cs3api

			// Deleting a space a two step process (1. disabling/trashing, 2. purging)
			// Do the "disable/trash" step only if the space is not marked as trashed yet:
			if _, ok := sp.Opaque.Map["trashed"]; !ok {
				_, err := client.DeleteStorageSpace(r.Context(), &storageprovider.DeleteStorageSpaceRequest{
					Id: &storageprovider.StorageSpaceId{
						OpaqueId: sp.Id.OpaqueId,
					},
				})
				if err != nil {
					logger.Error().Err(err).Msg("could not disable homespace: transport error")
					errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "could not disable homespace, aborting")
					return
				}
			}
			purgeFlag := utils.AppendPlainToOpaque(nil, "purge", "")
			_, err := client.DeleteStorageSpace(r.Context(), &storageprovider.DeleteStorageSpaceRequest{
				Opaque: purgeFlag,
				Id: &storageprovider.StorageSpaceId{
					OpaqueId: sp.Id.OpaqueId,
				},
			})
			if err != nil {
				// transport error, log as error
				logger.Error().Err(err).Msg("could not delete homespace: transport error")
				errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "could not delete homespace, aborting")
				return
			}
			break
		}
	}

	logger.Debug().Str("id", user.GetId()).Msg("calling delete user on backend")
	err = g.identityBackend.DeleteUser(r.Context(), user.GetId())

	if err != nil {
		logger.Debug().Err(err).Msg("could not delete user: backend error")
		errorcode.RenderError(w, r, err)
		return
	}

	g.publishEvent(r.Context(), e)

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// PatchMe implements the Service Interface. Updates the specified attributes of the current user
func (g Graph) PatchMe(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Debug().Msg("calling patch me")
	userID := revactx.ContextMustGetUser(r.Context()).GetId().GetOpaqueId()
	if userID == "" {
		logger.Debug().Msg("could not update user: missing user id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing user id")
		return
	}
	changes := libregraph.NewUser()
	err := StrictJSONUnmarshal(r.Body, changes)
	if err != nil {
		logger.Debug().Err(err).Interface("body", r.Body).Msg("could not update user: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest,
			fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}
	if _, ok := changes.GetDisplayNameOk(); ok {
		logger.Info().Interface("user", changes).Msg("could not update user: user is not allowed to change own displayname")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "user is not allowed to change own displayname")
		return
	}
	if _, ok := changes.GetMailOk(); ok {
		logger.Info().Interface("user", changes).Msg("could not update user: user is not allowed to change own mail")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "user is not allowed to change own mail")
		return
	}
	g.patchUser(w, r, userID, changes)
}

// PatchUser implements the Service Interface. Updates the specified attributes of an
// ExistingUser
func (g Graph) PatchUser(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	nameOrID := chi.URLParam(r, "userID")
	nameOrID, err := url.PathUnescape(nameOrID)
	if err != nil {
		logger.Debug().Err(err).Str("id", nameOrID).Msg("could not update user: unescaping user id failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping user id failed")
		return
	}
	changes := libregraph.NewUser()
	err = StrictJSONUnmarshal(r.Body, changes)
	if err != nil {
		logger.Debug().Err(err).Interface("body", r.Body).Msg("could not update user: invalid request body")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest,
			fmt.Sprintf("invalid request body: %s", err.Error()))
		return
	}
	if _, ok := changes.GetPreferredLanguageOk(); ok {
		logger.Info().Interface("user", changes).Msg("could not update user: user is not allowed to change other users language")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "user is not allowed to change other users language")
		return
	}

	g.patchUser(w, r, nameOrID, changes)
}

func (g Graph) patchUser(w http.ResponseWriter, r *http.Request, nameOrID string, changes *libregraph.User) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Debug().Msg("calling patch user")

	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")

	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	oldUserValues, err := g.identityBackend.GetUser(r.Context(), nameOrID, odataReq)
	if err != nil {
		logger.Debug().Err(err).Str("id", nameOrID).Msg("could not get user: backend error")
		errorcode.RenderError(w, r, err)
		return
	}

	if nameOrID == "" {
		logger.Debug().Msg("could not update user: missing user id")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing user id")
		return
	}

	if reflect.ValueOf(*changes).IsZero() {
		logger.Debug().Interface("body", r.Body).Msg("ignoring empty request body")
		render.Status(r, http.StatusNoContent)
		render.NoContent(w, r)
		return
	}

	if accountName, ok := changes.GetOnPremisesSamAccountNameOk(); ok {
		if !g.isValidUsername(*accountName) {
			logger.Info().Str("username", *accountName).Msg("could not update user: invalid username")
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "Invalid username")
			return
		}
	}

	preferredLanguage, ok := changes.GetPreferredLanguageOk()
	if ok {
		_, vID, err := getUserLanguage(r.Context(), g.valueService, oldUserValues.GetId())
		if err != nil {
			logger.Error().Err(err).Msg("could not get user language")
			tvID, err := uuid.NewUUID()
			if err != nil {
				logger.Error().Err(err).Msg("could not create user: error generating uuid")
				errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "error generating uuid")
				return
			}
			vID = tvID.String()
		}
		_, err = g.valueService.SaveValue(r.Context(), &settings.SaveValueRequest{
			Value: &settingsmsg.Value{
				Id:          vID,
				BundleId:    defaults.BundleUUIDProfile,
				SettingId:   defaults.SettingUUIDProfileLanguage,
				AccountUuid: oldUserValues.GetId(),
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
				},
				Value: &settingsmsg.Value_ListValue{
					ListValue: &settingsmsg.ListValue{Values: []*settingsmsg.ListOptionValue{
						{
							Option: &settingsmsg.ListOptionValue_StringValue{
								StringValue: *preferredLanguage,
							},
						},
					}},
				},
			},
		})
		if err != nil {
			logger.Error().Err(err).Msg("could not update user: error saving language setting")
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "error saving language setting")
			return
		}

	}

	var features []events.UserFeature
	if mail, ok := changes.GetMailOk(); ok {
		if !isValidEmail(*mail) {
			logger.Debug().Str("mail", *mail).Msg("could not update user: email is not a valid email address")
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest,
				fmt.Sprintf("'%s' is not a valid email address", *mail))
			return
		}
		features = append(features, events.UserFeature{
			Name:     "email",
			Value:    *mail,
			OldValue: oldUserValues.Mail,
		})
	}

	if name, ok := changes.GetDisplayNameOk(); ok {
		features = append(features, events.UserFeature{
			Name:     "displayname",
			Value:    *name,
			OldValue: oldUserValues.DisplayName,
		})
	}

	if userType, ok := changes.GetUserTypeOk(); ok {
		if !isValidUserType(*changes.UserType) {
			logger.Debug().Interface("user", changes).Msg("invalid userType attribute")
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid userType attribute, valid options are 'Member' or 'Guest'")
			return
		}
		features = append(features, events.UserFeature{
			Name:     "changeUserType",
			Value:    *userType,
			OldValue: oldUserValues.UserType,
		})
	}

	if accEnabled, ok := changes.GetAccountEnabledOk(); ok {
		oldAccVal := strconv.FormatBool(oldUserValues.GetAccountEnabled())
		features = append(features, events.UserFeature{
			Name:     "accountEnabled",
			Value:    strconv.FormatBool(*accEnabled),
			OldValue: &oldAccVal,
		})
	}

	if changes.HasPasswordProfile() {
		features = append(features, events.UserFeature{
			Name: "passwordChanged",
		})
	}

	logger.Debug().Str("nameid", nameOrID).Interface("changes", *changes).Msg("calling update user on backend")
	u, err := g.identityBackend.UpdateUser(r.Context(), nameOrID, *changes)
	if err != nil {
		logger.Debug().Err(err).Str("id", nameOrID).Msg("could not update user: backend error")
		errorcode.RenderError(w, r, err)
		return
	}
	u.PreferredLanguage = preferredLanguage

	e := events.UserFeatureChanged{
		UserID:    nameOrID,
		Features:  features,
		Timestamp: utils.TSNow(),
	}
	if currentUser, ok := revactx.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(r.Context(), e)

	render.Status(r, http.StatusOK) // TODO StatusNoContent when prefer=minimal is used
	render.JSON(w, r, u)
}

const (
	usernameMatchDefault = "default"
	usernameMatchNone    = "none"
)

var usernameRegexes = map[string]*regexp.Regexp{
	// We want to allow email addresses as usernames so they show up when using them in ACLs on storages that allow integration with our glauth LDAP service
	// so we are adding a few restrictions from https://stackoverflow.com/questions/6949667/what-are-the-real-rules-for-linux-usernames-on-centos-6-and-rhel-6
	// names should not start with numbers
	usernameMatchDefault: regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]*(@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)*$"),

	// In some cases users will be provisioned from an existing system, which may or may not have strange usernames. Because of this we want to "trust" the
	// upstream system and allow weird usernames, so relying on the used identity provider to complain if a username is violating its restrictions.
	usernameMatchNone: regexp.MustCompile(".*"),
}

func (g Graph) isValidUsername(e string) bool {
	if len(e) < 1 && len(e) > 254 {
		return false
	}

	return usernameRegexes[g.config.API.UsernameMatch].MatchString(e)
}

// regex from https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#valid-e-mail-address
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValidEmail(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func sortUsers(req *godata.GoDataRequest, users []*libregraph.User) ([]*libregraph.User, error) {
	if req.Query.OrderBy == nil || len(req.Query.OrderBy.OrderByItems) != 1 {
		return users, nil
	}
	var less func(i, j int) bool

	switch req.Query.OrderBy.OrderByItems[0].Field.Value {
	case displayNameAttr:
		less = func(i, j int) bool {
			return strings.ToLower(users[i].GetDisplayName()) < strings.ToLower(users[j].GetDisplayName())
		}
	case "mail":
		less = func(i, j int) bool {
			return strings.ToLower(users[i].GetMail()) < strings.ToLower(users[j].GetMail())
		}
	case "onPremisesSamAccountName":
		less = func(i, j int) bool {
			return strings.ToLower(users[i].GetOnPremisesSamAccountName()) < strings.ToLower(users[j].GetOnPremisesSamAccountName())
		}
	default:
		return nil, fmt.Errorf("we do not support <%s> as a order parameter", req.Query.OrderBy.OrderByItems[0].Field.Value)
	}

	if req.Query.OrderBy.OrderByItems[0].Order == _sortDescending {
		sort.Slice(users, reverse(less))
	} else {
		sort.Slice(users, less)
	}
	return users, nil
}

func isValidUserType(userType string) bool {
	userType = strings.ToLower(userType)

	for _, value := range []string{"member", "guest"} {
		if userType == value {
			return true
		}
	}

	return false
}
