package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec -- <command> [args...]",
	Short: "Run and profile a command (CPU, RAM, time)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			color.Red("✗ No command provided. Usage: vigil exec -- go build")
			os.Exit(1)
		}

		start := time.Now()

		// Start command
		c := exec.Command(args[0], args[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if !quiet {
			color.Cyan("▶ Running: %s", c.String())
		}

		err := c.Start()
		if err != nil {
			color.Red("✗ Failed to start: %v", err)
			os.Exit(1)
		}

		// Monitor in background
		var maxRAM float64
		var cpuSum float64
		var cpuSamples int

		ticker := time.NewTicker(200 * time.Millisecond)
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					// Get process memory
					p, err := process.NewProcess(int32(c.Process.Pid))
					if err == nil {
						memInfo, _ := p.MemoryInfo()
						if memInfo != nil {
							mb := float64(memInfo.RSS) / (1024 * 1024)
							if mb > maxRAM {
								maxRAM = mb
							}
						}
					}

					// Get CPU
					perc, _ := cpu.Percent(0, false)
					if len(perc) > 0 {
						cpuSum += perc[0]
						cpuSamples++
					}
				}
			}
		}()

		err = c.Wait()
		ticker.Stop()
		done <- true

		elapsed := time.Since(start).Seconds()
		avgCPU := 0.0
		if cpuSamples > 0 {
			avgCPU = cpuSum / float64(cpuSamples)
		}

		if jsonFlag {
			fmt.Printf(`{
  "command": %q,
  "exit_code": %d,
  "elapsed_seconds": %.3f,
  "cpu_avg_percent": %.1f,
  "ram_peak_mb": %.1f
}`, c.String(), c.ProcessState.ExitCode(), elapsed, avgCPU, maxRAM)
			return
		}

		status := "✅"
		if err != nil {
			status = "❌"
		}

		color.White("──────────────────────────────────────")
		color.Cyan("▶ Finished in %.2fs %s", elapsed, status)
		color.Green("   CPU: avg %.0f%%", avgCPU)
		color.Green("   RAM: peak %.1f MB", maxRAM)
		color.White("   Exit code: %d", c.ProcessState.ExitCode())
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}