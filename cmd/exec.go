// cmd/exec.go
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/sahil3982/vigil/internal/format"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec -- <command> [args...]",
	Short: "Run and profile a command (CPU, RAM, duration)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "✗ No command provided. Usage: vigil exec -- go build")
			os.Exit(1)
		}

		start := time.Now()

		c := exec.Command(args[0], args[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if !quiet {
			fmt.Fprintf(os.Stderr, "▶ Running: %s\n", c.String())
		}

		err := c.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ Failed to start: %v\n", err)
			os.Exit(1)
		}

		// Monitoring
		var maxRAM float64
		var cpuSum float64
		var cpuSamples int

		ticker := time.NewTicker(200 * time.Millisecond)
		done := make(chan bool, 1)

		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					// Memory
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

					// CPU
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
		close(done)

		elapsed := time.Since(start).Seconds()
		avgCPU := 0.0
		if cpuSamples > 0 {
			avgCPU = cpuSum / float64(cpuSamples)
		}

		stat := format.ExecStat{
			Command:        c.String(),
			ExitCode:       c.ProcessState.ExitCode(),
			ElapsedSeconds: elapsed,
			CPUAvgPercent:  avgCPU,
			RAMPeakMB:      maxRAM,
		}

		f := format.New(jsonFlag, quiet)
		if err := f.Exec(os.Stdout, stat); err != nil {
			os.Exit(1)
		}

		if stat.ExitCode != 0 && !jsonFlag && !quiet {
			os.Exit(stat.ExitCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}