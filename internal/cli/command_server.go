package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rangertaha/sao/internal/config"
	embeddednats "github.com/rangertaha/sao/internal/nats"
	"github.com/rangertaha/sao/internal/toc"
	"github.com/urfave/cli/v3"
)

// NewServerCommand builds the "server" command.
func NewServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Run the SAO TOC server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Usage: "Path to HCL config file",
				Value: config.DefaultPath,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			configPath := resolveConfigPath(cmd.String("config"))

			cfg, err := config.EnsureAndLoad(configPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			runCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()

			natsRuntime, err := embeddednats.Start(runCtx, cfg.NATS)
			if err != nil {
				return fmt.Errorf("start embedded nats: %w", err)
			}
			defer func() {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = natsRuntime.Shutdown(shutdownCtx)
			}()

			server := toc.NewServer(cfg, natsRuntime)
			if err := server.Run(runCtx); err != nil {
				return fmt.Errorf("run toc server: %w", err)
			}

			return nil
		},
	}
}

func resolveConfigPath(flagValue string) string {
	flagValue = strings.TrimSpace(flagValue)
	envValue := strings.TrimSpace(os.Getenv("SAO_CONFIG"))

	if envValue != "" && (flagValue == "" || flagValue == config.DefaultPath) {
		return envValue
	}
	if flagValue != "" {
		return flagValue
	}
	return config.DefaultPath
}
