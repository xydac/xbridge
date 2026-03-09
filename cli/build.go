package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xydac/xbridge/internal/client"
)

// NewBuildCmd creates the `xbridge build` command.
func NewBuildCmd() *cobra.Command {
	var (
		profile string
		clean   bool
		watch   bool
		remote  string
	)

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Trigger a build on the remote Mac",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := client.LoadClientConfig()
			if err != nil {
				return err
			}

			host, key := cfg.ResolveRemote(remote)
			if host == "" {
				return fmt.Errorf("no remote configured. Run: xbridge remote set <host>")
			}

			c := client.New(host, key)

			// Clean first if requested
			if clean {
				fmt.Println("Cleaning...")
				resp, err := c.Post("/build/clean", nil)
				if err != nil {
					return err
				}
				var result map[string]string
				if err := client.ReadJSON(resp, &result); err != nil {
					return err
				}
				fmt.Println("Clean complete.")
			}

			// Trigger build
			body := map[string]string{}
			if profile != "" {
				body["profile"] = profile
			}

			resp, err := c.Post("/build", body)
			if err != nil {
				return err
			}

			var result map[string]string
			if err := client.ReadJSON(resp, &result); err != nil {
				return err
			}

			fmt.Printf("Build started: %s\n", result["id"])

			if watch {
				fmt.Println("Streaming logs...")
				resp, err := c.Get(fmt.Sprintf("/build/%s/logs", result["id"]))
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				buf := make([]byte, 4096)
				for {
					n, err := resp.Body.Read(buf)
					if n > 0 {
						fmt.Print(string(buf[:n]))
					}
					if err != nil {
						break
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&profile, "profile", "", "Build profile to use")
	cmd.Flags().BoolVar(&clean, "clean", false, "Clean build folder first")
	cmd.Flags().BoolVar(&watch, "watch", false, "Stream build logs")
	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")

	return cmd
}
