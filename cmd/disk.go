// cmd/disk.go
package cmd

import (
	"os"

	"github.com/sahil3982/vigil/internal/format"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Show disk usage (root or first partition)",
	Run: func(cmd *cobra.Command, args []string) {
		// Try root first
		usage, err := disk.Usage("/")
		if err != nil {
			// Fall back to first partition
			parts, err2 := disk.Partitions(false)
			if err2 != nil || len(parts) == 0 {
				os.Exit(1)
			}
			usage, err = disk.Usage(parts[0].Mountpoint)
			if err != nil {
				os.Exit(1)
			}
		}

		stat := format.DiskStat{
			Path:        usage.Path,
			TotalBytes:  usage.Total,
			UsedBytes:   usage.Used,
			FreeBytes:   usage.Free,
			UsedPercent: usage.UsedPercent,
		}

		f := format.New(jsonFlag, quiet)
		if err := f.Disk(os.Stdout, stat); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(diskCmd)
}