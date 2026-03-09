package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.xbridge.{{.Name}}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.Binary}}</string>
        <string>serve</string>
        <string>--port</string>
        <string>{{.Port}}</string>
    </array>
    <key>WorkingDirectory</key>
    <string>{{.Project}}</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{.LogPath}}</string>
    <key>StandardErrorPath</key>
    <string>{{.LogPath}}</string>
</dict>
</plist>`

// NewInstallServiceCmd creates the `xbridge install-service` command.
func NewInstallServiceCmd() *cobra.Command {
	var (
		project string
		port    int
	)

	cmd := &cobra.Command{
		Use:   "install-service",
		Short: "Install xbridge as a launchd service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if project == "" {
				var err error
				project, err = os.Getwd()
				if err != nil {
					return err
				}
			}

			name := filepath.Base(project)
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			binary, err := os.Executable()
			if err != nil {
				return err
			}

			plistPath := filepath.Join(home, "Library", "LaunchAgents", fmt.Sprintf("com.xbridge.%s.plist", name))
			logPath := filepath.Join(home, "Library", "Logs", fmt.Sprintf("xbridge-%s.log", name))

			tmpl, err := template.New("plist").Parse(plistTemplate)
			if err != nil {
				return err
			}

			f, err := os.Create(plistPath)
			if err != nil {
				return err
			}
			defer f.Close()

			data := map[string]interface{}{
				"Name":    name,
				"Binary":  binary,
				"Port":    fmt.Sprintf("%d", port),
				"Project": project,
				"LogPath": logPath,
			}

			if err := tmpl.Execute(f, data); err != nil {
				return err
			}

			fmt.Printf("Service installed: %s\n", plistPath)
			fmt.Printf("Start with: launchctl load %s\n", plistPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&project, "project", "", "Project directory")
	cmd.Flags().IntVar(&port, "port", 7900, "Server port")

	return cmd
}

// NewUninstallServiceCmd creates the `xbridge uninstall-service` command.
func NewUninstallServiceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall-service [name]",
		Short: "Uninstall xbridge launchd service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			plistPath := filepath.Join(home, "Library", "LaunchAgents", fmt.Sprintf("com.xbridge.%s.plist", name))

			if _, err := os.Stat(plistPath); os.IsNotExist(err) {
				return fmt.Errorf("service %q not found at %s", name, plistPath)
			}

			if err := os.Remove(plistPath); err != nil {
				return err
			}

			fmt.Printf("Service removed: %s\n", plistPath)
			_ = strings.TrimSpace(name) // suppress unused
			return nil
		},
	}

	return cmd
}
