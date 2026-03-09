package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xydac/xbridge/internal/client"
)

// NewSimCmd creates the `xbridge sim` command group.
func NewSimCmd() *cobra.Command {
	var remote string

	cmd := &cobra.Command{
		Use:   "sim",
		Short: "Manage simulators on the remote Mac",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available simulators",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			resp, err := c.Get("/simulators")
			if err != nil {
				return err
			}

			var devices []map[string]interface{}
			if err := client.ReadJSON(resp, &devices); err != nil {
				return err
			}

			for _, d := range devices {
				state := d["state"]
				name := d["name"]
				udid := d["udid"]
				runtime := d["runtime"]
				fmt.Printf("%-30s %-15s %-10s %s\n", name, runtime, state, udid)
			}
			return nil
		},
	}

	bootCmd := &cobra.Command{
		Use:   "boot [device]",
		Short: "Boot a simulator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			resp, err := c.Post("/simulators/boot", map[string]string{"device": args[0]})
			if err != nil {
				return err
			}

			var result map[string]string
			if err := client.ReadJSON(resp, &result); err != nil {
				return err
			}

			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return nil
		},
	}

	shutdownCmd := &cobra.Command{
		Use:   "shutdown",
		Short: "Shutdown the booted simulator",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			resp, err := c.Post("/simulators/shutdown", map[string]string{})
			if err != nil {
				return err
			}

			var result map[string]string
			if err := client.ReadJSON(resp, &result); err != nil {
				return err
			}
			fmt.Println("Simulator shut down.")
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&remote, "remote", "", "Named remote to use")
	cmd.AddCommand(listCmd, bootCmd, shutdownCmd)
	return cmd
}

func makeClient(remote string) (*client.Client, error) {
	cfg, err := client.LoadClientConfig()
	if err != nil {
		return nil, err
	}
	host, key := cfg.ResolveRemote(remote)
	if host == "" {
		return nil, fmt.Errorf("no remote configured. Run: xbridge remote set <host>")
	}
	return client.New(host, key), nil
}
