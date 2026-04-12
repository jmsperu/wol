package cmd

import (
	"fmt"
	"os"
	"sort"
	"sync"
	"text/tabwriter"

	"github.com/jmsperu/wol/internal/config"
	"github.com/jmsperu/wol/internal/wol"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check which saved devices are online",
	Long:  "Checks reachability of all saved devices that have an IP address configured.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if len(cfg.Devices) == 0 {
			fmt.Println("No saved devices. Use 'wol add' to save one.")
			return nil
		}

		type result struct {
			name   string
			mac    string
			ip     string
			online bool
		}

		names := make([]string, 0, len(cfg.Devices))
		for name := range cfg.Devices {
			names = append(names, name)
		}
		sort.Strings(names)

		results := make([]result, len(names))
		var wg sync.WaitGroup

		for i, name := range names {
			d := cfg.Devices[name]
			results[i] = result{name: name, mac: d.MAC, ip: d.IP}
			if d.IP == "" {
				continue
			}
			wg.Add(1)
			go func(idx int, ip string) {
				defer wg.Done()
				results[idx].online = wol.IsHostUp(ip)
			}(i, d.IP)
		}
		wg.Wait()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tMAC\tIP\tSTATUS")
		for _, r := range results {
			ip := r.ip
			if ip == "" {
				ip = "-"
			}
			status := "unknown"
			if r.ip != "" {
				if r.online {
					status = "online"
				} else {
					status = "offline"
				}
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.name, r.mac, ip, status)
		}
		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
