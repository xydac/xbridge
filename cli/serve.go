package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/xydac/xbridge/internal/config"
	"github.com/xydac/xbridge/internal/engine"
	"github.com/xydac/xbridge/internal/server"
)

// NewServeCmd creates the `xbridge serve` command.
func NewServeCmd() *cobra.Command {
	var (
		port      int
		project   string
		scheme    string
		simulator string
		key       string
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the xbridge server",
		RunE: func(cmd *cobra.Command, args []string) error {
			workDir, err := os.Getwd()
			if err != nil {
				return err
			}

			cfg, err := config.Load(workDir)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			// CLI flags override config file
			if port != 0 {
				cfg.Server.Port = port
			}
			if project != "" {
				cfg.Project = project
			}
			if scheme != "" {
				cfg.Scheme = scheme
			}
			if simulator != "" {
				cfg.Simulator.Device = simulator
			}
			if key != "" {
				cfg.Server.Key = key
			}

			// Auto-detect project if not configured
			if cfg.Project == "" {
				detected, err := engine.DetectProject(workDir)
				if err != nil {
					return fmt.Errorf("detecting project: %w", err)
				}
				if detected != "" {
					cfg.Project = detected
					log.Printf("Detected %s", detected)
				}
			}

			// Auto-detect scheme if not configured
			if cfg.Scheme == "" && cfg.Project != "" {
				schemes, err := engine.ListSchemes(cfg.Project)
				if err == nil && len(schemes) > 0 {
					cfg.Scheme = schemes[0]
					log.Printf("Using scheme: %s", cfg.Scheme)
				}
			}

			srv := server.New(workDir, cfg)

			if cfg.Project != "" {
				log.Printf("Project: %s (scheme: %s)", cfg.Project, cfg.Scheme)
			}
			log.Println("Ready.")

			return srv.ListenAndServe()
		},
	}

	cmd.Flags().IntVar(&port, "port", 0, "Server port (default 7900)")
	cmd.Flags().StringVar(&project, "project", "", "Path to .xcworkspace or .xcodeproj")
	cmd.Flags().StringVar(&scheme, "scheme", "", "Xcode scheme to build")
	cmd.Flags().StringVar(&simulator, "simulator", "", "Simulator device name")
	cmd.Flags().StringVar(&key, "key", "", "API key for authentication")

	return cmd
}
