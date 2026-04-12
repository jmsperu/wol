package cmd

import (
	"fmt"

	"github.com/jmsperu/wol/internal/config"
	"github.com/spf13/cobra"
)

var (
	addIP        string
	addBroadcast string
	addPassword  string
)

var addCmd = &cobra.Command{
	Use:   "add <name> <MAC>",
	Short: "Save a device for quick access",
	Long: `Save a device with a friendly name.

Examples:
  wol add myserver AA:BB:CC:DD:EE:FF
  wol add myserver AA:BB:CC:DD:EE:FF -i 192.168.1.100
  wol add myserver AA:BB:CC:DD:EE:FF -i 192.168.1.100 -b 192.168.1.255`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		mac := args[1]

		if !macRegex.MatchString(mac) {
			return fmt.Errorf("invalid MAC address: %s", mac)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		cfg.Devices[name] = config.Device{
			MAC:       mac,
			IP:        addIP,
			Broadcast: addBroadcast,
			Password:  addPassword,
		}

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Saved device %q (MAC: %s", name, mac)
		if addIP != "" {
			fmt.Printf(", IP: %s", addIP)
		}
		if addBroadcast != "" {
			fmt.Printf(", broadcast: %s", addBroadcast)
		}
		fmt.Println(")")
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addIP, "ip", "i", "", "IP address for status checks")
	addCmd.Flags().StringVarP(&addBroadcast, "broadcast", "b", "", "broadcast address for this device")
	addCmd.Flags().StringVarP(&addPassword, "password", "p", "", "SecureOn password for this device")
	rootCmd.AddCommand(addCmd)
}
