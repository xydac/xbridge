package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xydac/xbridge/cli"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "xbridge",
		Short:   "A zero-config build proxy for your Mac mini",
		Long:    "Single binary. No dependencies. Point it at your Xcode project. Build and test from any machine on your network.",
		Version: version,
	}

	rootCmd.AddCommand(
		cli.NewServeCmd(),
		cli.NewBuildCmd(),
		cli.NewSimCmd(),
		cli.NewRunCmd(),
		cli.NewPullCmd(),
		cli.NewCheckoutCmd(),
		cli.NewStatusCmd(),
		cli.NewRemoteCmd(),
		cli.NewScreenshotCmd(),
		cli.NewInstallServiceCmd(),
		cli.NewUninstallServiceCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
