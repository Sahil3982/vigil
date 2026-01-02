// internal/format/types.go
package format

type CPUStat struct {
	Percent float64 `json:"cpu_percent"`
	Cores   int     `json:"cores,omitempty"`
}

type MemStat struct {
	TotalBytes     uint64  `json:"total_bytes"`
	UsedBytes      uint64  `json:"used_bytes"`
	FreeBytes      uint64  `json:"free_bytes"`
	UsedPercent    float64 `json:"used_percent"`
	AvailableBytes uint64  `json:"available_bytes,omitempty"`
}

type DiskStat struct {
	Path        string  `json:"mount"`
	TotalBytes  uint64  `json:"total_bytes"`
	UsedBytes   uint64  `json:"used_bytes"`
	FreeBytes   uint64  `json:"free_bytes"`
	UsedPercent float64 `json:"used_percent"`
}

type ExecStat struct {
	Command        string  `json:"command"`
	ExitCode       int     `json:"exit_code"`
	ElapsedSeconds float64 `json:"elapsed_seconds"`
	CPUAvgPercent  float64 `json:"cpu_avg_percent"`
	RAMPeakMB      float64 `json:"ram_peak_mb"`
}