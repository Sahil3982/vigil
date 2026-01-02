package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

var memCmd = &cobra.Command{
	Use:   "mem",
	Short: "Show memory (RAM) usage",
	Run: func(cmd *cobra.Command, args []string) {
		v, err := mem.VirtualMemory()
		if err != nil {
			panic(err)
		}

		if jsonFlag {
			fmt.Printf(`{
  "total_bytes": %d,
  "used_bytes": %d,
  "free_bytes": %d,
  "used_percent": %.2f
}`, v.Total, v.Used, v.Available, v.UsedPercent)
			return
		}

		totalGB := float64(v.Total) / (1024 * 1024 * 1024)
		usedGB := float64(v.Used) / (1024 * 1024 * 1024)
		bar := barFor(v.UsedPercent, 100)

		status := "✅"
		if v.UsedPercent > 85 {
			status = "⚠️"
		}
		if v.UsedPercent > 95 {
			status = "🔥"
		}

		color.Green("▶ RAM: %s %.1f%% (%.1f/%.1f GB) %s",
			bar, v.UsedPercent, usedGB, totalGB, status)
	},
}

func init() {
	rootCmd.AddCommand(memCmd)
}