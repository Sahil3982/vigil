// cmd/cpu.go
package cmd

import (
	"os"
	"time"

	"github.com/sahil3982/vigil/internal/format"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/spf13/cobra"
)

var cpuCmd = &cobra.Command{
	Use:   "cpu",
	Short: "Show CPU usage percentage",
	Run: func(cmd *cobra.Command, args []string) {
		// Prime the pump (first call returns 0 on some systems)
		_, _ = cpu.Percent(0, false)
		time.Sleep(200 * time.Millisecond)

		percent, err := cpu.Percent(300*time.Millisecond, false)
		if err != nil || len(percent) == 0 {
			// Fallback: try per-core and average
			percent, err = cpu.Percent(500*time.Millisecond, true)
			if err != nil || len(percent) == 0 {
				os.Exit(1)
			}
			sum := 0.0
			for _, p := range percent {
				sum += p
			}
			percent = []float64{sum / float64(len(percent))}
		}

		stat := format.CPUStat{
			Percent: percent[0],
			Cores:   len(percent),
		}

		f := format.New(jsonFlag, quiet)
		if err := f.CPU(os.Stdout, stat); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(cpuCmd)
}
