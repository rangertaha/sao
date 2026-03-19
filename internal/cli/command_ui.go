package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rangertaha/sao/internal/config"
	internalui "github.com/rangertaha/sao/internal/ui"
	"github.com/urfave/cli/v3"
)

// NewUICommand builds the "ui" command.
func NewUICommand() *cli.Command {
	return &cli.Command{
		Name:  "ui",
		Usage: "Serve embedded SAO React UI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Usage: "Path to HCL config file",
				Value: config.DefaultPath,
			},
			&cli.StringFlag{
				Name:  "addr",
				Usage: "HTTP listen address override",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg, err := config.EnsureAndLoad(resolveConfigPath(cmd.String("config")))
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			addr := cfg.UI.Address
			if override := cmd.String("addr"); override != "" {
				addr = override
			}

			runCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()

			if err := internalui.Serve(runCtx, addr); err != nil {
				return fmt.Errorf("serve embedded ui: %w", err)
			}

			return nil
		},
	}
}
