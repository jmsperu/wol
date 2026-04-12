package cmd

import (
	"fmt"
	"strings"

	"github.com/jmsperu/wol/internal/config"
	"github.com/spf13/cobra"
)

var wakeCmd = &cobra.Command{
	Use:   "wake [device...] | [MAC...]",
	Short: "Wake one or more saved devices or MAC addresses",
	Long: `Wake saved devices by name or MAC addresses.
Pass multiple names/MACs to batch wake.

Examples:
  wol wake myserver
  wol wake myserver mynas mydesktop
  wol wake AA:BB:CC:DD:EE:FF`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		var errors []string
		for _, target := range args {
			if macRegex.MatchString(target) {
				if err := sendWake(target, "", broadcast, password, port, wait); err != nil {
					errors = append(errors, fmt.Sprintf("%s: %v", target, err))
				}
				continue
			}

			dev, ok := cfg.Devices[target]
			if !ok {
				errors = append(errors, fmt.Sprintf("%s: device not found", target))
				continue
			}
			bcast := broadcast
			if bcast == "" && dev.Broadcast != "" {
				bcast = dev.Broadcast
			}
			pw := password
			if pw == "" && dev.Password != "" {
				pw = dev.Password
			}
			if err := sendWake(dev.MAC, dev.IP, bcast, pw, port, wait); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", target, err))
			}
		}

		if len(errors) > 0 {
			return fmt.Errorf("some devices failed:\n  %s", strings.Join(errors, "\n  "))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(wakeCmd)
}
