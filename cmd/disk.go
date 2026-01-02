package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/spf13/cobra"
)

var diskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Show disk usage (root mount)",
	Run: func(cmd *cobra.Command, args []string) {
		usage, err := disk.Usage("/")
		if err != nil {
			// fallback to first partition
			parts, _ := disk.Partitions(false)
			if len(parts) > 0 {
				usage, _ = disk.Usage(parts[0].Mountpoint)
			}
		}
		if usage == nil {
			color.Red("✗ Could not determine disk usage")
			return
		}

		if jsonFlag {
			fmt.Printf(`{
  "mount": "%s",
  "total_bytes": %d,
  "used_bytes": %d,
  "free_bytes": %d,
  "used_percent": %.2f
}`, usage.Path, usage.Total, usage.Used, usage.Free, usage.UsedPercent)
			return
		}

		totalGB := float64(usage.Total) / (1024 * 1024 * 1024)
		usedGB := float64(usage.Used) / (1024 * 1024 * 1024)
		bar := barFor(usage.UsedPercent, 100)

		status := "✅"
		if usage.UsedPercent > 85 {
			status = "⚠️"
		}
		if usage.UsedPercent > 95 {
			status = "🔥"
		}

		color.Yellow("▶ Disk %s: %s %.1f%% (%.1f/%.1f GB) %s",
			usage.Path, bar, usage.UsedPercent, usedGB, totalGB, status)
	},
}

func init() {
	rootCmd.AddCommand(diskCmd)
}