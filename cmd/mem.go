// cmd/mem.go
package cmd

import (
	"os"

	"github.com/sahil3982/vigil/internal/format"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

var memCmd = &cobra.Command{
	Use:   "mem",
	Short: "Show memory (RAM) usage",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := mem.VirtualMemory()
		if err != nil {
			os.Exit(1)
		}

		stat := format.MemStat{
			TotalBytes:     v.Total,
			UsedBytes:      v.Used,
			FreeBytes:      v.Free,
			AvailableBytes: v.Available,
			UsedPercent:    v.UsedPercent,
		}

		f := format.New(jsonFlag, quiet)
		if err := f.Mem(os.Stdout, stat); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(memCmd)
}