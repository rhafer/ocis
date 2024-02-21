package svc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ldapv3 "github.com/go-ldap/ldap/v3"
	"github.com/jellydator/ttlcache/v3"
	microstore "go-micro.dev/v4/store"

	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/store"

	ocisldap "github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity/ldap"
	graphm "github.com/owncloud/ocis/v2/services/graph/pkg/middleware"
)

const (
	// HeaderPurge defines the header name for the purge header.
	HeaderPurge     = "Purge"
	displayNameAttr = "displayName"
)

// Service defines the service handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)

	ListApplications(w http.ResponseWriter, r *http.Request)
	GetApplication(http.ResponseWriter, *http.Request)

	GetMe(http.ResponseWriter, *http.Request)
	GetUsers(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request)
	PostUser(http.ResponseWriter, *http.Request)
	DeleteUser(http.ResponseWriter, *http.Request)
	PatchUser(http.ResponseWriter, *http.Request)
	ChangeOwnPassword(http.ResponseWriter, *http.Request)

	ListAppRoleAssignments(http.ResponseWriter, *http.Request)
	CreateAppRoleAssignment(http.ResponseWriter, *http.Request)
	DeleteAppRoleAssignment(http.ResponseWriter, *http.Request)

	GetGroups(http.ResponseWriter, *http.Request)
	GetGroup(http.ResponseWriter, *http.Request)
	PostGroup(http.ResponseWriter, *http.Request)
	PatchGroup(http.ResponseWriter, *http.Request)
	DeleteGroup(http.ResponseWriter, *http.Request)
	GetGroupMembers(http.ResponseWriter, *http.Request)
	PostGroupMember(http.ResponseWriter, *http.Request)
	DeleteGroupMember(http.ResponseWriter, *http.Request)

	GetEducationSchools(http.ResponseWriter, *http.Request)
	GetEducationSchool(http.ResponseWriter, *http.Request)
	PostEducationSchool(http.ResponseWriter, *http.Request)
	PatchEducationSchool(http.ResponseWriter, *http.Request)
	DeleteEducationSchool(http.ResponseWriter, *http.Request)
	GetEducationSchoolUsers(http.ResponseWriter, *http.Request)
	PostEducationSchoolUser(http.ResponseWriter, *http.Request)
	DeleteEducationSchoolUser(http.ResponseWriter, *http.Request)
	GetEducationSchoolClasses(http.ResponseWriter, *http.Request)
	PostEducationSchoolClass(http.ResponseWriter, *http.Request)
	DeleteEducationSchoolClass(http.ResponseWriter, *http.Request)

	GetEducationClasses(http.ResponseWriter, *http.Request)
	GetEducationClass(http.ResponseWriter, *http.Request)
	PostEducationClass(http.ResponseWriter, *http.Request)
	PatchEducationClass(http.ResponseWriter, *http.Request)
	DeleteEducationClass(w http.ResponseWriter, r *http.Request)
	GetEducationClassMembers(w http.ResponseWriter, r *http.Request)
	PostEducationClassMember(w http.ResponseWriter, r *http.Request)

	GetEducationUsers(http.ResponseWriter, *http.Request)
	GetEducationUser(http.ResponseWriter, *http.Request)
	PostEducationUser(http.ResponseWriter, *http.Request)
	DeleteEducationUser(http.ResponseWriter, *http.Request)
	PatchEducationUser(http.ResponseWriter, *http.Request)
	DeleteEducationClassMember(w http.ResponseWriter, r *http.Request)

	GetEducationClassTeachers(w http.ResponseWriter, r *http.Request)
	PostEducationClassTeacher(w http.ResponseWriter, r *http.Request)
	DeleteEducationClassTeacher(w http.ResponseWriter, r *http.Request)

	GetDrivesV1(w http.ResponseWriter, r *http.Request)
	GetDrivesV1Beta1(w http.ResponseWriter, r *http.Request)
	GetSingleDrive(w http.ResponseWriter, r *http.Request)
	GetAllDrivesV1(w http.ResponseWriter, r *http.Request)
	GetAllDrivesV1Beta1(w http.ResponseWriter, r *http.Request)
	CreateDrive(w http.ResponseWriter, r *http.Request)
	UpdateDrive(w http.ResponseWriter, r *http.Request)
	DeleteDrive(w http.ResponseWriter, r *http.Request)

	GetSharedByMe(w http.ResponseWriter, r *http.Request)
	ListSharedWithMe(w http.ResponseWriter, r *http.Request)

	GetRootDriveChildren(w http.ResponseWriter, r *http.Request)
	GetDriveItem(w http.ResponseWriter, r *http.Request)
	GetDriveItemChildren(w http.ResponseWriter, r *http.Request)
	CreateLink(w http.ResponseWriter, r *http.Request)
	SetLinkPassword(writer http.ResponseWriter, request *http.Request)

	Invite(w http.ResponseWriter, r *http.Request)
	ListPermissions(w http.ResponseWriter, r *http.Request)
	UpdatePermission(w http.ResponseWriter, r *http.Request)
	DeletePermission(w http.ResponseWriter, r *http.Request)

	CreateUploadSession(w http.ResponseWriter, r *http.Request)

	GetTags(w http.ResponseWriter, r *http.Request)
	AssignTags(w http.ResponseWriter, r *http.Request)
	UnassignTags(w http.ResponseWriter, r *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) (Graph, error) {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	spacePropertiesCache := ttlcache.New(
		ttlcache.WithTTL[string, interface{}](
			time.Duration(options.Config.Spaces.ExtendedSpacePropertiesCacheTTL),
		),
		ttlcache.WithDisableTouchOnHit[string, interface{}](),
	)
	go spacePropertiesCache.Start()

	identityCache := identity.NewIdentityCache(
		identity.IdentityCacheWithGatewaySelector(options.GatewaySelector),
		identity.IdentityCacheWithUsersTTL(time.Duration(options.Config.Spaces.UsersCacheTTL)),
		identity.IdentityCacheWithGroupsTTL(time.Duration(options.Config.Spaces.GroupsCacheTTL)),
	)

	svc := Graph{
		config:                   options.Config,
		mux:                      m,
		logger:                   &options.Logger,
		specialDriveItemsCache:   spacePropertiesCache,
		identityCache:            identityCache,
		eventsPublisher:          options.EventsPublisher,
		gatewaySelector:          options.GatewaySelector,
		searchService:            options.SearchService,
		identityEducationBackend: options.IdentityEducationBackend,
		keycloakClient:           options.KeycloakClient,
		historyClient:            options.EventHistoryClient,
		traceProvider:            options.TraceProvider,
		valueService:             options.ValueService,
	}

	if err := setIdentityBackends(options, &svc); err != nil {
		return svc, err
	}

	if options.PermissionService == nil {
		grpcClient, err := grpc.NewClient(append(grpc.GetClientOptions(options.Config.GRPCClientTLS), grpc.WithTraceProvider(options.TraceProvider))...)
		if err != nil {
			return svc, err
		}
		svc.permissionsService = settingssvc.NewPermissionService("com.owncloud.api.settings", grpcClient)
	} else {
		svc.permissionsService = options.PermissionService
	}

	svc.roleService = options.RoleService

	roleManager := options.RoleManager
	if roleManager == nil {
		storeOptions := []microstore.Option{
			store.Store(options.Config.Cache.Store),
			store.TTL(options.Config.Cache.TTL),
			store.Size(options.Config.Cache.Size),
			microstore.Nodes(options.Config.Cache.Nodes...),
			microstore.Database(options.Config.Cache.Database),
			microstore.Table(options.Config.Cache.Table),
			store.DisablePersistence(options.Config.Cache.DisablePersistence),
			store.Authentication(options.Config.Cache.AuthUsername, options.Config.Cache.AuthPassword),
		}
		m := roles.NewManager(
			roles.StoreOptions(storeOptions),
			roles.Logger(options.Logger),
			roles.RoleService(options.RoleService),
		)
		roleManager = &m
	}

	var requireAdmin func(http.Handler) http.Handler
	if options.RequireAdminMiddleware == nil {
		requireAdmin = graphm.RequireAdmin(roleManager, options.Logger)
	} else {
		requireAdmin = options.RequireAdminMiddleware
	}

	drivesDriveItemService, err := NewDrivesDriveItemService(options.Logger, options.GatewaySelector)
	if err != nil {
		return svc, err
	}

	drivesDriveItemApi, err := NewDrivesDriveItemApi(drivesDriveItemService, options.Logger)
	if err != nil {
		return svc, err
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Use(middleware.StripSlashes)

		r.Route("/v1beta1", func(r chi.Router) {
			r.Route("/me", func(r chi.Router) {
				r.Get("/drives", svc.GetDrives(APIVersion_1_Beta_1))
				r.Route("/drive", func(r chi.Router) {
					r.Get("/sharedByMe", svc.GetSharedByMe)
					r.Get("/sharedWithMe", svc.ListSharedWithMe)
				})
			})
			r.Route("/drives", func(r chi.Router) {
				r.Get("/", svc.GetAllDrives(APIVersion_1_Beta_1))
				r.Route("/{driveID}/items/{itemID}", func(r chi.Router) {
					r.Delete("/", drivesDriveItemApi.DeleteDriveItem)
					r.Post("/children", drivesDriveItemApi.CreateDriveItem)
					r.Post("/invite", svc.Invite)
					r.Route("/permissions", func(r chi.Router) {
						r.Get("/", svc.ListPermissions)
						r.Route("/{permissionID}", func(r chi.Router) {
							r.Delete("/", svc.DeletePermission)
							r.Patch("/", svc.UpdatePermission)
							r.Post("/setPassword", svc.SetLinkPassword)
						})
					})
					r.Post("/createLink", svc.CreateLink)
				})
			})
			r.Route("/roleManagement/permissions/roleDefinitions", func(r chi.Router) {
				r.Get("/", svc.GetRoleDefinitions)
				r.Get("/{roleID}", svc.GetRoleDefinition)
			})
		})
		r.Route("/v1.0", func(r chi.Router) {
			r.Route("/extensions/org.libregraph", func(r chi.Router) {
				r.Get("/tags", svc.GetTags)
				r.Put("/tags", svc.AssignTags)
				r.Delete("/tags", svc.UnassignTags)
			})
			r.Route("/applications", func(r chi.Router) {
				r.Get("/", svc.ListApplications)
				r.Get("/{applicationID}", svc.GetApplication)
			})
			r.Route("/me", func(r chi.Router) {
				r.Get("/", svc.GetMe)
				r.Patch("/", svc.PatchMe)
				r.Route("/drive", func(r chi.Router) {
					r.Get("/", svc.GetUserDrive)
					r.Get("/root/children", svc.GetRootDriveChildren)
				})
				r.Get("/drives", svc.GetDrives(APIVersion_1))
				r.Post("/changePassword", svc.ChangeOwnPassword)
			})
			r.Route("/users", func(r chi.Router) {
				r.Get("/", svc.GetUsers)
				r.With(requireAdmin).Post("/", svc.PostUser)
				r.Route("/{userID}", func(r chi.Router) {
					r.Get("/", svc.GetUser)
					r.Get("/drive", svc.GetUserDrive)
					r.Post("/exportPersonalData", svc.ExportPersonalData)
					r.With(requireAdmin).Delete("/", svc.DeleteUser)
					r.With(requireAdmin).Patch("/", svc.PatchUser)
					if svc.roleService != nil {
						r.With(requireAdmin).Route("/appRoleAssignments", func(r chi.Router) {
							r.Get("/", svc.ListAppRoleAssignments)
							r.Post("/", svc.CreateAppRoleAssignment)
							r.Delete("/{appRoleAssignmentID}", svc.DeleteAppRoleAssignment)
						})
					}
				})
			})
			r.Route("/groups", func(r chi.Router) {
				r.Get("/", svc.GetGroups)
				r.With(requireAdmin).Post("/", svc.PostGroup)
				r.Route("/{groupID}", func(r chi.Router) {
					r.Get("/", svc.GetGroup)
					r.With(requireAdmin).Delete("/", svc.DeleteGroup)
					r.With(requireAdmin).Patch("/", svc.PatchGroup)
					r.Route("/members", func(r chi.Router) {
						r.With(requireAdmin).Get("/", svc.GetGroupMembers)
						r.With(requireAdmin).Post("/$ref", svc.PostGroupMember)
						r.With(requireAdmin).Delete("/{memberID}/$ref", svc.DeleteGroupMember)
					})
				})
			})
			r.Route("/drives", func(r chi.Router) {
				r.Get("/", svc.GetAllDrives(APIVersion_1))
				r.Post("/", svc.CreateDrive)
				r.Route("/{driveID}", func(r chi.Router) {
					r.Patch("/", svc.UpdateDrive)
					r.Get("/", svc.GetSingleDrive)
					r.Delete("/", svc.DeleteDrive)
					r.Route("/items/{driveItemID}", func(r chi.Router) {
						r.Get("/", svc.GetDriveItem)
						r.Get("/children", svc.GetDriveItemChildren)
						r.Post("/createUploadSession", svc.CreateUploadSession)
					})
				})
			})
			r.With(requireAdmin).Route("/education", func(r chi.Router) {
				r.Route("/schools", func(r chi.Router) {
					r.Get("/", svc.GetEducationSchools)
					r.Post("/", svc.PostEducationSchool)
					r.Route("/{schoolID}", func(r chi.Router) {
						r.Get("/", svc.GetEducationSchool)
						r.Delete("/", svc.DeleteEducationSchool)
						r.Patch("/", svc.PatchEducationSchool)
						r.Route("/users", func(r chi.Router) {
							r.Get("/", svc.GetEducationSchoolUsers)
							r.Post("/$ref", svc.PostEducationSchoolUser)
							r.Delete("/{userID}/$ref", svc.DeleteEducationSchoolUser)
						})
						r.Route("/classes", func(r chi.Router) {
							r.Get("/", svc.GetEducationSchoolClasses)
							r.Post("/$ref", svc.PostEducationSchoolClass)
							r.Delete("/{classID}/$ref", svc.DeleteEducationSchoolClass)
						})
					})
				})
				r.Route("/users", func(r chi.Router) {
					r.Get("/", svc.GetEducationUsers)
					r.Post("/", svc.PostEducationUser)
					r.Route("/{userID}", func(r chi.Router) {
						r.Get("/", svc.GetEducationUser)
						r.Delete("/", svc.DeleteEducationUser)
						r.Patch("/", svc.PatchEducationUser)
					})
				})
				r.Route("/classes", func(r chi.Router) {
					r.Get("/", svc.GetEducationClasses)
					r.Post("/", svc.PostEducationClass)
					r.Route("/{classID}", func(r chi.Router) {
						r.Get("/", svc.GetEducationClass)
						r.Delete("/", svc.DeleteEducationClass)
						r.Patch("/", svc.PatchEducationClass)
						r.Route("/members", func(r chi.Router) {
							r.Get("/", svc.GetEducationClassMembers)
							r.Post("/$ref", svc.PostEducationClassMember)
							r.Delete("/{memberID}/$ref", svc.DeleteEducationClassMember)
						})
						r.Route("/teachers", func(r chi.Router) {
							r.Get("/", svc.GetEducationClassTeachers)
							r.Post("/$ref", svc.PostEducationClassTeacher)
							r.Delete("/{teacherID}/$ref", svc.DeleteEducationClassTeacher)
						})
					})
				})
			})
		})
	})

	_ = chi.Walk(m, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		options.Logger.Debug().Str("method", method).Str("route", route).Int("middlewares", len(middlewares)).Msg("serving endpoint")
		return nil
	})

	return svc, nil
}

func setIdentityBackends(options Options, svc *Graph) error {
	if options.IdentityBackend == nil {
		switch options.Config.Identity.Backend {
		case "cs3":
			gatewaySelector, err := pool.GatewaySelector(
				options.Config.Reva.Address,
				append(
					options.Config.Reva.GetRevaOptions(),
					pool.WithRegistry(registry.GetRegistry()),
					pool.WithTracerProvider(options.TraceProvider),
				)...,
			)
			if err != nil {
				return err
			}

			svc.identityBackend = &identity.CS3{
				Config:          options.Config.Reva,
				Logger:          &options.Logger,
				GatewaySelector: gatewaySelector,
			}
		case "ldap":
			var err error

			var tlsConf *tls.Config
			if options.Config.Identity.LDAP.Insecure {

				// When insecure is set to true then we don't need a certificate.
				options.Config.Identity.LDAP.CACert = ""
				tlsConf = &tls.Config{
					MinVersion: tls.VersionTLS12,

					//nolint:gosec // We need the ability to run with "insecure" (dev/testing)
					InsecureSkipVerify: options.Config.Identity.LDAP.Insecure,
				}
			}

			if options.Config.Identity.LDAP.CACert != "" {
				if err := ocisldap.WaitForCA(options.Logger,
					options.Config.Identity.LDAP.Insecure,
					options.Config.Identity.LDAP.CACert); err != nil {
					options.Logger.Fatal().Err(err).Msg("The configured LDAP CA cert does not exist")
				}
				if tlsConf == nil {
					tlsConf = &tls.Config{
						MinVersion: tls.VersionTLS12,
					}
				}
				certs := x509.NewCertPool()
				pemData, err := os.ReadFile(options.Config.Identity.LDAP.CACert)
				if err != nil {
					options.Logger.Error().Err(err).Msg("Error initializing LDAP Backend")
					return err
				}
				if !certs.AppendCertsFromPEM(pemData) {
					options.Logger.Error().Msg("Error initializing LDAP Backend. Adding CA cert failed")
					return err
				}
				tlsConf.RootCAs = certs
			}

			conn := ldap.NewLDAPWithReconnect(&options.Logger,
				ldap.Config{
					URI:          options.Config.Identity.LDAP.URI,
					BindDN:       options.Config.Identity.LDAP.BindDN,
					BindPassword: options.Config.Identity.LDAP.BindPassword,
					TLSConfig:    tlsConf,
				},
			)
			lb, err := identity.NewLDAPBackend(conn, options.Config.Identity.LDAP, &options.Logger)
			if err != nil {
				options.Logger.Error().Err(err).Msg("Error initializing LDAP Backend")
				return err
			}
			svc.identityBackend = lb
			if options.IdentityEducationBackend == nil {
				if options.Config.Identity.LDAP.EducationResourcesEnabled {
					svc.identityEducationBackend = lb
				} else {
					errEduBackend := &identity.ErrEducationBackend{}
					svc.identityEducationBackend = errEduBackend
				}
			}

			disableMechanismType, err := identity.ParseDisableMechanismType(options.Config.Identity.LDAP.DisableUserMechanism)
			if err != nil {
				options.Logger.Error().Err(err).Msg("Error initializing LDAP Backend")
				return err
			}

			if disableMechanismType == identity.DisableMechanismGroup {
				options.Logger.Info().Msg("LocalUserDisable is true, will create group if not exists")
				err := lb.CreateLDAPGroupByDN(options.Config.Identity.LDAP.LdapDisabledUsersGroupDN)
				if err != nil {
					isAnError := false
					var lerr *ldapv3.Error
					if errors.As(err, &lerr) {
						if lerr.ResultCode != ldapv3.LDAPResultEntryAlreadyExists {
							isAnError = true
						}
					} else {
						isAnError = true
					}

					if isAnError {
						msg := "error adding group for disabling users"
						options.Logger.Error().Err(err).Str("local_user_disable", options.Config.Identity.LDAP.LdapDisabledUsersGroupDN).Msg(msg)
						return err
					}
				}
			}

		default:
			err := fmt.Errorf("unknown identity backend: '%s'", options.Config.Identity.Backend)
			options.Logger.Err(err)
			return err
		}
	} else {
		svc.identityBackend = options.IdentityBackend
	}

	return nil
}

// parseHeaderPurge parses the 'Purge' header.
// '1', 't', 'T', 'TRUE', 'true', 'True' are parsed as true
// all other values are false.
func parsePurgeHeader(h http.Header) bool {
	val := h.Get(HeaderPurge)

	if b, err := strconv.ParseBool(val); err == nil {
		return b
	}
	return false
}
