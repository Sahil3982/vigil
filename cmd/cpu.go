package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/spf13/cobra"
)

var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Show CPU usage",
	Run: func(cmd *cobra.Command, args []string) {
		// Seed first reading
		_, _ = cpu.Percent(0, false)
		time.Sleep(500 * time.Millisecond)

		percent, err := cpu.Percent(500*time.Millisecond, false)
		if err != nil {
			panic(err)
		}

		p := percent[0]
		if jsonFlag {
			fmt.Printf(`{"cpu_percent": %.2f}`, p)
			return
		}

		bar := barFor(p, 100)
		status := "✅"
		if p > 80 {
			status = "⚠️"
		}
		if p > 95 {
			status = "🔥"
		}
		color.Cyan("▶ CPU: %s %.1f%% %s", bar, p, status)
	},
}

func init() {
	rootCmd.AddCommand(cpuCmd)
}