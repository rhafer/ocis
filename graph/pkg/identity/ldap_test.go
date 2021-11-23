package identity

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/owncloud/ocis/graph/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// ldapMock implements the ldap.Client interfac
type ldapMock struct {
	SearchFunc func(*ldap.SearchRequest) (*ldap.SearchResult, error)
}

var lconfig = config.LDAP{
	UserBaseDN:               "dc=test",
	UserSearchScope:          "sub",
	UserFilter:               "filter",
	UserDisplayNameAttribute: "displayname",
	UserIDAttribute:          "entryUUID",
	UserEmailAttribute:       "mail",
	UserNameAttribute:        "uid",

	GroupBaseDN:        "dc=test",
	GroupSearchScope:   "sub",
	GroupFilter:        "filter",
	GroupNameAttribute: "cn",
	GroupIDAttribute:   "entryUUID",
}

var logger = log.NewLogger(log.Level("debug"))

func TestNewLDAPBackend(t *testing.T) {

	l := ldapMock{}

	tc := lconfig
	tc.UserDisplayNameAttribute = ""
	if _, err := NewLDAPBackend(l, tc, &logger); err == nil {
		t.Error("Should fail with incomplete user attr config")
	}

	tc = lconfig
	tc.GroupIDAttribute = ""
	if _, err := NewLDAPBackend(l, tc, &logger); err == nil {
		t.Errorf("Should fail with incomplete group	config")
	}

	tc = lconfig
	tc.UserSearchScope = ""
	if _, err := NewLDAPBackend(l, tc, &logger); err == nil {
		t.Errorf("Should fail with invalid user search scope")
	}

	tc = lconfig
	tc.GroupSearchScope = ""
	if _, err := NewLDAPBackend(l, tc, &logger); err == nil {
		t.Errorf("Should fail with invalid group search scope")
	}

	if _, err := NewLDAPBackend(l, lconfig, &logger); err != nil {
		t.Errorf("Should fail with invalid group search scope")
	}

}

func TestCreateUserModelFromLDAP(t *testing.T) {
	l := ldapMock{}
	logger := log.NewLogger(log.Level("debug"))

	b, _ := NewLDAPBackend(l, lconfig, &logger)
	if user := b.createUserModelFromLDAP(nil); user != nil {
		t.Errorf("createUserModelFromLDAP should return on nil Entry")
	}
	userAttrs := map[string][]string{
		"uid":         {"userid"},
		"displayname": {"DisplayName"},
		"mail":        {"user@example"},
		"entryuuid":   {"abcd-defg"},
	}

	e := ldap.NewEntry("uid=user", userAttrs)
	user := b.createUserModelFromLDAP(e)
	if user == nil {
		t.Error("Converting a valid LDAP Entry should succeed")
	}
	if *user.OnPremisesSamAccountName != e.GetEqualFoldAttributeValue(b.userAttributeMap.userName) {
		t.Errorf("Error creating msGraph User from LDAP Entry: %s != %s", *user.OnPremisesSamAccountName, e.GetEqualFoldAttributeValue(b.userAttributeMap.userName))
	}
	if *user.Mail != e.GetEqualFoldAttributeValue(b.userAttributeMap.mail) {
		t.Errorf("Error creating msGraph User from LDAP Entry: %s != %s", *user.Mail, e.GetEqualFoldAttributeValue(b.userAttributeMap.mail))
	}
	if *user.DisplayName != e.GetEqualFoldAttributeValue(b.userAttributeMap.displayName) {
		t.Errorf("Error creating msGraph User from LDAP Entry: %s != %s", *user.DisplayName, e.GetEqualFoldAttributeValue(b.userAttributeMap.displayName))
	}
	if *user.ID != e.GetEqualFoldAttributeValue(b.userAttributeMap.id) {
		t.Errorf("Error creating msGraph User from LDAP Entry: %s != %s", *user.ID, e.GetEqualFoldAttributeValue(b.userAttributeMap.id))
	}
}

func TestGetUser(t *testing.T) {
	// Mock a Sizelimit Error correctly
	lm := ldapMock{
		SearchFunc: func(*ldap.SearchRequest) (*ldap.SearchResult, error) {
			return nil, ldap.NewError(ldap.LDAPResultSizeLimitExceeded, errors.New("mock"))
		},
	}
	b, _ := NewLDAPBackend(lm, lconfig, &logger)
	_, err := b.GetUser(context.Background(), "fred")
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}

	// Mock an empty Search Result
	lm = ldapMock{
		SearchFunc: func(*ldap.SearchRequest) (*ldap.SearchResult, error) {
			return &ldap.SearchResult{}, nil
		},
	}
	b, _ = NewLDAPBackend(lm, lconfig, &logger)
	_, err = b.GetUser(context.Background(), "fred")
	if err == nil || err.Error() != "itemNotFound" {
		t.Errorf("Expected 'itemNotFound' got '%s'", err.Error())
	}
}

// below here ldap.Client interface method for ldapMock

func (c ldapMock) Start() {}

func (c ldapMock) StartTLS(*tls.Config) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) Close() {}

func (c ldapMock) IsClosing() bool {
	return false
}

func (c ldapMock) SetTimeout(time.Duration) {}

func (c ldapMock) Bind(username, password string) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) UnauthenticatedBind(username string) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) SimpleBind(*ldap.SimpleBindRequest) (*ldap.SimpleBindResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) ExternalBind() error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) Add(*ldap.AddRequest) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) Del(*ldap.DelRequest) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) Modify(*ldap.ModifyRequest) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) ModifyDN(*ldap.ModifyDNRequest) error {
	return ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) ModifyWithResult(*ldap.ModifyRequest) (*ldap.ModifyResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) Compare(dn, attribute, value string) (bool, error) {
	return false, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) PasswordModify(*ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}

func (c ldapMock) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	if c.SearchFunc != nil {
		return c.SearchFunc(searchRequest)
	}

	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}
func (c ldapMock) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	return nil, ldap.NewError(ldap.LDAPResultNotSupported, fmt.Errorf("not implemented"))
}
