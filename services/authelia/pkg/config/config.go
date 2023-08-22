package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	Ldap Ldap `yaml:"ldap"`

	Context context.Context `yaml:"-"`
}

// Ldap defines the available LDAP configuration.
type Ldap struct {
	URI          string `yaml:"uri" env:"OCIS_LDAP_URI;LDAP_URI;IDP_LDAP_URI" desc:"Url of the LDAP service to use as IDP." deprecationVersion:"3.0" removalVersion:"4.0.0" deprecationInfo:"LDAP_URI changing name for consistency" deprecationReplacement:"OCIS_LDAP_URI"`
	BindDN       string `yaml:"bind_dn" env:"OCIS_LDAP_BIND_DN;LDAP_BIND_DN;IDP_LDAP_BIND_DN" desc:"LDAP DN to use for simple bind authentication with the target LDAP server." deprecationVersion:"3.0" removalVersion:"4.0.0" deprecationInfo:"LDAP_BIND_DN changing name for consistency" deprecationReplacement:"OCIS_LDAP_BIND_DN"`
	BindPassword string `yaml:"bind_password" env:"LDAP_BIND_PASSWORD;IDP_LDAP_BIND_PASSWORD" desc:"Password to use for authenticating the 'bind_dn'."`

	BaseDN string `yaml:"base_dn" env:"OCIS_LDAP_USER_BASE_DN;LDAP_USER_BASE_DN;IDP_LDAP_BASE_DN" desc:"Search base DN for looking up LDAP users." deprecationVersion:"3.0" removalVersion:"4.0.0" deprecationInfo:"LDAP_USER_BASE_DN changing name for consistency" deprecationReplacement:"OCIS_LDAP_USER_BASE_DN"`
}
