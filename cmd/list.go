package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/jmsperu/wol/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List saved devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if len(cfg.Devices) == 0 {
			fmt.Println("No saved devices. Use 'wol add' to save one.")
			return nil
		}

		names := make([]string, 0, len(cfg.Devices))
		for name := range cfg.Devices {
			names = append(names, name)
		}
		sort.Strings(names)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tMAC\tIP\tBROADCAST")
		for _, name := range names {
			d := cfg.Devices[name]
			ip := d.IP
			if ip == "" {
				ip = "-"
			}
			bcast := d.Broadcast
			if bcast == "" {
				bcast = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, d.MAC, ip, bcast)
		}
		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
