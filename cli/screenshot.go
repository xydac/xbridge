package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

// NewScreenshotCmd creates the `xbridge screenshot` command.
func NewScreenshotCmd() *cobra.Command {
	var (
		remote   string
		openFile bool
		output   string
	)

	cmd := &cobra.Command{
		Use:   "screenshot",
		Short: "Take a screenshot of the simulator",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			resp, err := c.Get("/simulators/BOOTED/screenshot")
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if output == "" {
				output = fmt.Sprintf("screenshot-%s.png", time.Now().Format("20060102-150405"))
			}

			f, err := os.Create(output)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(f, resp.Body); err != nil {
				return err
			}

			fmt.Printf("Screenshot saved to %s\n", output)

			if openFile {
				switch runtime.GOOS {
				case "darwin":
					exec.Command("open", output).Run()
				case "linux":
					exec.Command("xdg-open", output).Run()
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")
	cmd.Flags().BoolVar(&openFile, "open", false, "Open screenshot in image viewer")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path")

	return cmd
}
