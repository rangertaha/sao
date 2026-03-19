package cli

import "github.com/urfave/cli/v3"

// NewApp returns the root sao command.
func NewApp() *cli.Command {
	return &cli.Command{
		Name:  "sao",
		Usage: "Spatial Awareness Operator server",
		Commands: []*cli.Command{
			NewServerCommand(),
			NewUICommand(),
		},
	}
}
