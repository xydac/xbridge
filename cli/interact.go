package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xydac/xbridge/internal/client"
)

// NewTapCmd creates the `xbridge tap` command.
func NewTapCmd() *cobra.Command {
	var (
		remote   string
		duration float64
	)

	cmd := &cobra.Command{
		Use:   "tap <x> <y>",
		Short: "Tap at coordinates on the simulator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"x": mustInt(args[0]),
				"y": mustInt(args[1]),
			}
			if duration > 0 {
				body["duration"] = duration
			}

			resp, err := c.Post("/simulators/BOOTED/tap", body)
			if err != nil {
				return err
			}

			var result map[string]string
			if err := client.ReadJSON(resp, &result); err != nil {
				return err
			}
			fmt.Println("Tapped.")
			return nil
		},
	}

	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")
	cmd.Flags().Float64Var(&duration, "duration", 0, "Long press duration in seconds")
	return cmd
}

// NewSwipeCmd creates the `xbridge swipe` command.
func NewSwipeCmd() *cobra.Command {
	var (
		remote   string
		duration float64
	)

	cmd := &cobra.Command{
		Use:   "swipe <x1> <y1> <x2> <y2>",
		Short: "Swipe on the simulator",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"x1": mustInt(args[0]),
				"y1": mustInt(args[1]),
				"x2": mustInt(args[2]),
				"y2": mustInt(args[3]),
			}
			if duration > 0 {
				body["duration"] = duration
			}

			resp, err := c.Post("/simulators/BOOTED/swipe", body)
			if err != nil {
				return err
			}

			var result map[string]string
			if err := client.ReadJSON(resp, &result); err != nil {
				return err
			}
			fmt.Println("Swiped.")
			return nil
		},
	}

	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")
	cmd.Flags().Float64Var(&duration, "duration", 0, "Swipe duration in seconds")
	return cmd
}

// NewTextCmd creates the `xbridge text` command.
func NewTextCmd() *cobra.Command {
	var remote string

	cmd := &cobra.Command{
		Use:   "text <string>",
		Short: "Type text into the simulator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			resp, err := c.Post("/simulators/BOOTED/text", map[string]string{"text": args[0]})
			if err != nil {
				return err
			}

			var result map[string]string
			if err := client.ReadJSON(resp, &result); err != nil {
				return err
			}
			fmt.Println("Typed.")
			return nil
		},
	}

	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")
	return cmd
}

// NewKeyCmd creates the `xbridge key` command.
func NewKeyCmd() *cobra.Command {
	var (
		remote   string
		duration float64
	)

	cmd := &cobra.Command{
		Use:   "key <keycode>",
		Short: "Send a key press to the simulator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := makeClient(remote)
			if err != nil {
				return err
			}

			body := map[string]interface{}{
				"key": args[0],
			}
			if duration > 0 {
				body["duration"] = duration
			}

			resp, err := c.Post("/simulators/BOOTED/key", body)
			if err != nil {
				return err
			}

			var result map[string]string
			if err := client.ReadJSON(resp, &result); err != nil {
				return err
			}
			fmt.Println("Pressed.")
			return nil
		},
	}

	cmd.Flags().StringVar(&remote, "remote", "", "Named remote to use")
	cmd.Flags().Float64Var(&duration, "duration", 0, "Key press duration in seconds")
	return cmd
}

func mustInt(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

