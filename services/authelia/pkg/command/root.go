package command

import (
	"os"

	"github.com/owncloud/ocis/v2/ocis-pkg/clihelper"
	"github.com/owncloud/ocis/v2/services/authelia/pkg/config"
	"github.com/urfave/cli/v2"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) cli.Commands {
	return []*cli.Command{
		// start this service
		Server(cfg),

		Version(cfg),
	}
}

// Execute is the entry point for the ocis-idm command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cli.App{
		Name:     "authelia",
		Usage:    "Embedded Authelia instance for oCIS",
		Commands: GetCommands(cfg),
	})

	return app.Run(os.Args)
}
