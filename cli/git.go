package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xydac/xbridge/internal/client"
)

// NewPullCmd creates the `xbridge pull` command.
func NewPullCmd() *cobra.Command {
	var remote string
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Git pull on the remote Mac",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}
			resp, err := c.Post("/git/pull", nil)
			if err != nil {
				return err
			}
			var result map[string]string
			if err := client.ReadJSON(resp, &result); err != nil {
				return err
			}
			fmt.Println(result["output"])
			return nil
		},
	}
	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")
	return cmd
}

// NewCheckoutCmd creates the `xbridge checkout` command.
func NewCheckoutCmd() *cobra.Command {
	var remote string
	cmd := &cobra.Command{
		Use:   "checkout [branch]",
		Short: "Switch branch on the remote Mac",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}
			resp, err := c.Post("/git/checkout", map[string]string{"branch": args[0]})
			if err != nil {
				return err
			}
			var result map[string]string
			if err := client.ReadJSON(resp, &result); err != nil {
				return err
			}
			fmt.Println("Switched to", args[0])
			return nil
		},
	}
	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")
	return cmd
}

// NewStatusCmd creates the `xbridge status` command.
func NewStatusCmd() *cobra.Command {
	var remote string
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show git + build + sim status",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			// Get health
			healthResp, err := c.Get("/health")
			if err != nil {
				return err
			}
			var health map[string]interface{}
			client.ReadJSON(healthResp, &health)

			// Get git status
			gitResp, err := c.Get("/git/status")
			if err == nil {
				var gitStatus map[string]interface{}
				if err := client.ReadJSON(gitResp, &gitStatus); err == nil {
					data, _ := json.MarshalIndent(gitStatus, "", "  ")
					fmt.Println("Git:", string(data))
				}
			}

			data, _ := json.MarshalIndent(health, "", "  ")
			fmt.Println("Health:", string(data))

			return nil
		},
	}
	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")
	return cmd
}
