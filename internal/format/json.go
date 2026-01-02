// internal/format/json.go
package format

import (
	"encoding/json"
	"io"
)

type JSONFormatter struct{}

func (j *JSONFormatter) CPU(w io.Writer, stat CPUStat) error {
	return json.NewEncoder(w).Encode(stat)
}

func (j *JSONFormatter) Mem(w io.Writer, stat MemStat) error {
	return json.NewEncoder(w).Encode(stat)
}

func (j *JSONFormatter) Disk(w io.Writer, stat DiskStat) error {
	return json.NewEncoder(w).Encode(stat)
}

func (j *JSONFormatter) Exec(w io.Writer, stat ExecStat) error {
	return json.NewEncoder(w).Encode(stat)
}