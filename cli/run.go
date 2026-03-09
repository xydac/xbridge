package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/xydac/xbridge/internal/client"
)

// NewRunCmd creates the `xbridge run` combo command.
func NewRunCmd() *cobra.Command {
	var (
		profile string
		watch   bool
		remote  string
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Pull, build, install, and launch in one command",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			// 1. Pull
			fmt.Println("Pulling latest...")
			if resp, err := c.Post("/git/pull", nil); err == nil {
				var result map[string]string
				client.ReadJSON(resp, &result)
			}

			// 2. Build
			fmt.Println("Building...")
			body := map[string]string{}
			if profile != "" {
				body["profile"] = profile
			}

			resp, err := c.Post("/build", body)
			if err != nil {
				return err
			}

			var buildResult map[string]string
			if err := client.ReadJSON(resp, &buildResult); err != nil {
				return err
			}

			buildID := buildResult["id"]
			fmt.Printf("Build started: %s\n", buildID)

			// 3. Wait for build to complete
			for {
				resp, err := c.Get(fmt.Sprintf("/build/%s", buildID))
				if err != nil {
					return err
				}
				var status map[string]interface{}
				if err := client.ReadJSON(resp, &status); err != nil {
					return err
				}

				s := status["status"].(string)
				if s == "success" {
					fmt.Println("Build succeeded!")
					break
				}
				if s == "failed" {
					return fmt.Errorf("build failed: %v", status["error"])
				}

				time.Sleep(2 * time.Second)
			}

			fmt.Println("Done!")

			if watch {
				fmt.Println("Streaming logs...")
				resp, err := c.Get(fmt.Sprintf("/build/%s/logs", buildID))
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
	cmd.Flags().BoolVar(&watch, "watch", false, "Stream logs after build")
	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")

	return cmd
}
