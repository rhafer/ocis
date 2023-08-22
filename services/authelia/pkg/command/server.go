package command

import (
	"context"
	"fmt"

	autheliacmd "github.com/authelia/authelia/v4/pkg/commands"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/authelia/pkg/config"
	"github.com/owncloud/ocis/v2/services/authelia/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/authelia/pkg/logging"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			var (
				gr          = run.Group{}
				logger      = logging.Configure(cfg.Service.Name, cfg.Log)
				ctx, cancel = func() (context.Context, context.CancelFunc) {
					if cfg.Context == nil {
						return context.WithCancel(context.Background())
					}
					return context.WithCancel(cfg.Context)
				}()
			)

			defer cancel()

			{
				var err error
				cmd := autheliacmd.NewRootCmd()
				// To trigger persistant flags merge
				cmd.InitDefaultHelpFlag()

				cfgMap := AutheliaConfigFromStruct(cfg)

				cmd.SetArgs([]string{"config", "./configuration.yaml"})
				actx := autheliacmd.NewCmdCtx()
				if err = actx.ConfigEnsureExistsRunE(cmd, []string{}); err != nil {
					logger.Error().Err(err).Msg("ConfigEnsureExistsRunE")
					return err
				}

				if err = actx.ConfigSetDefaultsRunE(cfgMap)(cmd, []string{}); err != nil {
					logger.Error().Err(err).Msg("ConfigSetDefaultsRunE")
					return err
				}

				if err = actx.ConfigLoadRunE(cmd, []string{}); err != nil {
					logger.Error().Err(err).Msg("ConfigLoadRunE")
					return err
				}
				if err = actx.ConfigValidateKeysRunE(cmd, []string{}); err != nil {
					logger.Error().Err(err).Msg("ConfigValidateKeysRunE")
					return err
				}
				if err = actx.ConfigValidateRunE(cmd, []string{}); err != nil {
					logger.Error().Err(err).Msg("ConfigValidateRunE")
					return err
				}
				if err = actx.ConfigValidateLogRunE(cmd, []string{}); err != nil {
					logger.Error().Err(err).Msg("ConfigValidateLogRunE")
					return err
				}
				if err = actx.RootRunE(nil, []string{}); err != nil {
					logger.Error().Err(err).Msg("RootRunE")
					return err
				}

				//				cmd.Execute()
			}

			{
				server := debug.NewService(
					debug.Logger(logger),
					debug.Name(cfg.Service.Name),
					debug.Version(version.GetString()),
					debug.Address(cfg.Debug.Addr),
					debug.Token(cfg.Debug.Token),
					debug.Pprof(cfg.Debug.Pprof),
					debug.Zpages(cfg.Debug.Zpages),
					debug.Health(handlers.Health),
					debug.Ready(handlers.Ready),
				)

				gr.Add(server.ListenAndServe, func(_ error) {
					_ = server.Shutdown(ctx)
					cancel()
				})
			}

			return gr.Run()
			//return start(ctx, logger, cfg)
		},
	}
}

func AutheliaConfigFromStruct(cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"jwt_secret":                                       "a_very_important_secret",
		"server.address":                                   "tcp://0.0.0.0:9091/authelia",
		"authentication_backend.ldap.address":              cfg.Ldap.URI,
		"authentication_backend.ldap.user":                 cfg.Ldap.BindDN,
		"authentication_backend.ldap.password":             cfg.Ldap.BindPassword,
		"authentication_backend.ldap.tls.skip_verify":      true,
		"authentication_backend.ldap.base_dn":              cfg.Ldap.BaseDN,
		"authentication_backend.ldap.additional_users_dn":  "",
		"authentication_backend.ldap.users_filter":         "(&({username_attribute}={input})(objectClass=person))",
		"authentication_backend.ldap.additional_groups_dn": "",
		"authentication_backend.ldap.groups_filter":        "(&(member={dn})(objectClass=groupOfNames))",
		"access_control.default_policy":                    "one_factor",
		"session.domain":                                   "localhost.localdomain",
		"session.name":                                     "authelia_session",
		"storage.encryption_key":                           "you_must_generate_a_random_string_of_more_than_twenty_chars_and_configure_this",
		"storage.local.path":                               "db.sqlite3",
		"notifier.filesystem.filename":                     "./nofication.txt",
		"identity_providers.oidc.hmac_secret":              "this_is_a_secret_abc123abc123abc",
		"identity_providers.oidc.cors.allowed_origins":     "https://localhost.localdomain",
		"identity_providers.oidc.issuer_private_key": `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAzbWlhR1GzlFSru1RGk3mAprDUCimD7cMr5+qR4D43mCUtf/u
rQfN2P+DKkqec4c8RwK5GJ1b4/uqraz+iGTSytqlV6McphqsDpMqcQS2fGOqmrrq
3KzA3nh86e+1vDit91sxehvXZzr8Bp+vHXYOYODBdgV9kldQndRcszBsRIxOL3VA
/fQgJRHBCXKP8aBahPDpiUjTrZ1idaCOCybidiG7dVikCWMZ0VVaVVICbnpFI8Tg
cV4B04lgxc5/pXLPCDiS85nX2oaQtuder8UcyvlEypUUaI7M5gdjJwpbpbJDcK1s
3pQ7r/9c+S+sGTEeJ8Yfrrb9eO4BCrD3icupRwIDAQABAoIBAHd2wkVoFzLCb64/
DBamnqlsj3kB4k7GE4v6wbz83Yrns/VuSIIcQiN/YAEzjImzRAJJRo1Q9YiVIy3x
hXBYfUJpcBRRGdYtLzbXqJIyFnhuuSla3AKEIQ4SExkYqQZCNGWuhDFR8ep1it+d
5OjLItaIMFIGJkLSinLeXWeC+51igs4ZcsONm4qOWZdSPptp77Nf59KVGl+1S7ii
5hgFsPAuShKu0eAJV3uHn2x/BR7qKhxa/yvz374Q0QJ3EcjQ9PT5mHY6qJKSss/D
aaOco28zmQzQnfB+s7XQ0N99yO+REhkdQA+7A5+2sVlnujZQv86ukyF9RXvSvwLj
UapMqoECgYEA9y+4mGZFNGykbcKgjp+F5lv2Yf+Hj8SP3L9VJ9vBLivJjSvlNNOK
UW3g0xKZDbjAGpvpLT7FduOu1ih7XPaWdhfUXgDfU1WSfzVs/QMD67qeLm/BOGKo
iqjdDxklVPZGwNJNv6WFl2wrNz2wjveuH86+Q3UYP8AxgeXBDrCq6dECgYEA1QtU
zMXZfYs20KRlieD28fxD7m68eEMGZXePxeShSsikJSH5toQXa6sb4kCNvFunioKo
eeyYhJMHLtHleDADPW3itDWeu6OhPu57UFPmPwBvQ6jr1xCoNvZZH+/8ROFXiR1s
SYG7x2lKqAFVmiTooD0BtrgfyKogGzxGYM5gj5cCgYAbhqLlxa27MsX0uxGqEDWW
+3KqYwwzhE4I5P2UnLIcdB/TqqmxgkUK4FOC7bVBg+tQi0AiG7VdkeksTAHAzmze
5bRua2ZzHzpbFBX47tcG7xciUKuRndrq5fcH8WLo3Svv2Ptzdfk0bYU6d5IruYUY
YatqU6XJo5tfvbgL7Lx7oQKBgG3cvgorLRD0rXvCiyoi/LWlJVLbgA10YuQIV/fx
AswR07PiZWedjoZTYrm2GGE90pQ29LKLM8uKFnYqf28PM1yGQhY0YHra0tglGyo9
Wcq7aqU1gwkQ6e4N87/ofer3WbC3n5P4duFKhtlEduRajCu1yiBzqtBMCuqAlrpt
MpZZAoGAKkuA4Ge6/xmr0oYziheJfjOrrK657uo9/wZW0JtxRgjq0Wfr7dsfeVqd
Yk1oChMgoJWj4zhYcFgqWl+4/fSMvJBkc0E1ABiaIxm49ycIHEJvlJq9212Gn/XJ
PN8DxYw2Ck2cFHjiqBKQFwE++/t3nHvp8GYeegqI/1cZhWuR/sE=
-----END RSA PRIVATE KEY-----`,
		"identity_providers.oidc.clients": []interface{}{
			map[string]interface{}{
				"id":          "web",
				"description": "ownCloud web client",
				"public":      true,
				"redirect_uris": []string{
					"https://localhost.localdomain:9200/",
					"https://localhost.localdomain:9200/oidc-callback.html",
					"https://localhost.localdomain:9200/oidc-silent-redirect.html",
				},
				"authorization_policy": "one_factor",
			},
		},
	}
	return rcfg
}
