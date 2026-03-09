package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xydac/xbridge/internal/client"
)

// NewRemoteCmd creates the `xbridge remote` command group.
func NewRemoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage remote xbridge servers",
	}

	var key string

	setCmd := &cobra.Command{
		Use:   "set [host]",
		Short: "Set the default remote",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := client.LoadClientConfig()
			if err != nil {
				return err
			}
			cfg.DefaultRemote = args[0]
			if key != "" {
				cfg.Key = key
			}
			if err := client.SaveClientConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("Default remote set to %s\n", args[0])
			return nil
		},
	}
	setCmd.Flags().StringVar(&key, "key", "", "API key")

	addCmd := &cobra.Command{
		Use:   "add [name] [host]",
		Short: "Add a named remote",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := client.LoadClientConfig()
			if err != nil {
				return err
			}
			cfg.Remotes[args[0]] = client.Remote{
				Host: args[1],
				Key:  key,
			}
			if err := client.SaveClientConfig(cfg); err != nil {
				return err
			}
			fmt.Printf("Remote %q added (%s)\n", args[0], args[1])
			return nil
		},
	}
	addCmd.Flags().StringVar(&key, "key", "", "API key")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List configured remotes",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := client.LoadClientConfig()
			if err != nil {
				return err
			}
			if cfg.DefaultRemote != "" {
				fmt.Printf("default: %s\n", cfg.DefaultRemote)
			}
			for name, r := range cfg.Remotes {
				fmt.Printf("  %s: %s\n", name, r.Host)
			}
			return nil
		},
	}

	cmd.AddCommand(setCmd, addCmd, listCmd)
	return cmd
}
