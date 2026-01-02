// internal/format/human.go
package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

type HumanFormatter struct {
	Quiet bool
}

func (h *HumanFormatter) bar(value, max float64) string {
	perc := value / max
	width := 10
	filled := int(perc * float64(width))
	empty := width - filled
	return "[" + strings.Repeat("■", filled) + strings.Repeat("□", empty) + "]"
}

func (h *HumanFormatter) statusIcon(percent float64) string {
	if percent > 95 {
		return color.RedString("🔥")
	}
	if percent > 80 {
		return color.YellowString("⚠️")
	}
	return color.GreenString("✅")
}

func (h *HumanFormatter) CPU(w io.Writer, stat CPUStat) error {
	if h.Quiet {
		_, err := fmt.Fprintf(w, "%.1f", stat.Percent)
		return err
	}
	bar := h.bar(stat.Percent, 100)
	status := h.statusIcon(stat.Percent)
	_, err := color.New(color.FgCyan).Fprintf(w, "▶ CPU: %s %.1f%% %s\n", bar, stat.Percent, status)
	return err
}

func (h *HumanFormatter) Mem(w io.Writer, stat MemStat) error {
	if h.Quiet {
		_, err := fmt.Fprintf(w, "%.1f", stat.UsedPercent)
		return err
	}
	bar := h.bar(stat.UsedPercent, 100)
	status := h.statusIcon(stat.UsedPercent)
	totalGB := float64(stat.TotalBytes) / (1024 * 1024 * 1024)
	usedGB := float64(stat.UsedBytes) / (1024 * 1024 * 1024)
	_, err := color.New(color.FgGreen).Fprintf(w, "▶ RAM: %s %.1f%% (%.1f/%.1f GB) %s\n",
		bar, stat.UsedPercent, usedGB, totalGB, status)
	return err
}

func (h *HumanFormatter) Disk(w io.Writer, stat DiskStat) error {
	if h.Quiet {
		_, err := fmt.Fprintf(w, "%.1f", stat.UsedPercent)
		return err
	}
	bar := h.bar(stat.UsedPercent, 100)
	status := h.statusIcon(stat.UsedPercent)
	totalGB := float64(stat.TotalBytes) / (1024 * 1024 * 1024)
	usedGB := float64(stat.UsedBytes) / (1024 * 1024 * 1024)
	_, err := color.New(color.FgYellow).Fprintf(w, "▶ Disk %s: %s %.1f%% (%.1f/%.1f GB) %s\n",
		stat.Path, bar, stat.UsedPercent, usedGB, totalGB, status)
	return err
}

func (h *HumanFormatter) Exec(w io.Writer, stat ExecStat) error {
	status := "✅"
	if stat.ExitCode != 0 {
		status = "❌"
	}
	if h.Quiet {
		_, err := fmt.Fprintf(w, "%d", stat.ExitCode)
		return err
	}
	color.New(color.FgWhite).Fprintf(w, "──────────────────────────────────────\n")
	color.New(color.FgCyan).Fprintf(w, "▶ Finished in %.2fs %s\n", stat.ElapsedSeconds, status)
	color.New(color.FgGreen).Fprintf(w, "   CPU: avg %.0f%%\n", stat.CPUAvgPercent)
	color.New(color.FgGreen).Fprintf(w, "   RAM: peak %.1f MB\n", stat.RAMPeakMB)
	color.New(color.FgWhite).Fprintf(w, "   Exit code: %d\n", stat.ExitCode)
	return nil
}
