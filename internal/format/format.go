// internal/format/format.go
package format

import "io"

// Formatter defines how data is rendered
type Formatter interface {
	CPU(w io.Writer, stat CPUStat) error
	Mem(w io.Writer, stat MemStat) error
	Disk(w io.Writer, stat DiskStat) error
	Exec(w io.Writer, stat ExecStat) error
}

// New returns a formatter based on flags
func New(jsonFlag, quietFlag bool) Formatter {
	if jsonFlag {
		return &JSONFormatter{}
	}
	return &HumanFormatter{Quiet: quietFlag}
}