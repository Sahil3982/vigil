// internal/format/csv.go
package format

import (
	"fmt"
	"io"
	"time"
)

type CSVFormatter struct{}

func (c *CSVFormatter) CPU(w io.Writer, stat CPUStat) error {
	fmt.Fprintf(w, "%s,cpu,%.2f,%d\n", time.Now().Format("2006-01-02 15:04:05"), stat.Percent, stat.Cores)
	return nil
}

func (c *CSVFormatter) Mem(w io.Writer, stat MemStat) error {
	fmt.Fprintf(w, "%s,mem,%.2f,%d,%d\n", time.Now().Format("2006-01-02 15:04:05"), stat.UsedPercent, stat.UsedBytes, stat.TotalBytes)
	return nil
}