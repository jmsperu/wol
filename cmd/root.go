package cmd

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/jmsperu/wol/internal/config"
	"github.com/jmsperu/wol/internal/wol"
	"github.com/spf13/cobra"
)

var (
	broadcast string
	password  string
	port      int
	wait      bool
)

var macRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2}$`)

var rootCmd = &cobra.Command{
	Use:   "wol [MAC or device name]",
	Short: "Wake-on-LAN tool — send magic packets to wake devices",
	Long: `Wake-on-LAN CLI tool.

Send magic packets to wake devices on your network.
Save devices by name for quick access, check their status, and batch wake.

Examples:
  wol AA:BB:CC:DD:EE:FF              # wake by MAC address
  wol wake myserver                  # wake a saved device
  wol add myserver AA:BB:CC:DD:EE:FF # save a device
  wol list                           # list saved devices
  wol status                         # check which devices are online`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		target := args[0]

		// Check if it's a MAC address
		if macRegex.MatchString(target) {
			return sendWake(target, "", broadcast, password, port, wait)
		}

		// Otherwise treat as device name
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		dev, ok := cfg.Devices[target]
		if !ok {
			return fmt.Errorf("device %q not found — use 'wol add' to save it, or pass a MAC address", target)
		}
		bcast := broadcast
		if bcast == "" && dev.Broadcast != "" {
			bcast = dev.Broadcast
		}
		pw := password
		if pw == "" && dev.Password != "" {
			pw = dev.Password
		}
		return sendWake(dev.MAC, dev.IP, bcast, pw, port, wait)
	},
}

func sendWake(mac, ip, bcast, pw string, port int, waitPing bool) error {
	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("invalid MAC address %q: %w", mac, err)
	}

	var secureOn []byte
	if pw != "" {
		secureOn, err = parsePassword(pw)
		if err != nil {
			return fmt.Errorf("invalid SecureOn password: %w", err)
		}
	}

	if bcast == "" {
		bcast = "255.255.255.255"
	}

	fmt.Printf("Sending magic packet to %s via %s:%d", mac, bcast, port)
	if pw != "" {
		fmt.Print(" (with SecureOn password)")
	}
	fmt.Println("...")

	if err := wol.SendMagicPacket(hwAddr, bcast, port, secureOn); err != nil {
		return fmt.Errorf("failed to send magic packet: %w", err)
	}
	fmt.Println("Magic packet sent.")

	if waitPing && ip != "" {
		fmt.Printf("Waiting for %s to come online...\n", ip)
		if wol.WaitForHost(ip, 60) {
			fmt.Printf("%s is online!\n", ip)
		} else {
			fmt.Printf("%s did not respond within 60 seconds.\n", ip)
		}
	}
	return nil
}

func parsePassword(pw string) ([]byte, error) {
	pw = strings.ReplaceAll(pw, "-", ":")
	parts := strings.Split(pw, ":")
	if len(parts) != 6 {
		return nil, fmt.Errorf("password must be 6 hex bytes (e.g. AA:BB:CC:DD:EE:FF)")
	}
	hwAddr, err := net.ParseMAC(pw)
	if err != nil {
		return nil, err
	}
	return []byte(hwAddr), nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&broadcast, "broadcast", "b", "", "broadcast address (default 255.255.255.255)")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "SecureOn password (6 hex bytes, e.g. AA:BB:CC:DD:EE:FF)")
	rootCmd.PersistentFlags().IntVar(&port, "port", 9, "UDP port (default 9)")
	rootCmd.PersistentFlags().BoolVarP(&wait, "wait", "w", false, "wait and check if device comes online after wake")
}
