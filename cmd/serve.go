// cmd/serve.go
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
)

var (
	port           int
	historyLimit   int
	metricsHistory []map[string]interface{}
	historyMutex   sync.RWMutex
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server for comprehensive metrics dashboard",
	Long: `Start a comprehensive HTTP dashboard with:
  • Real-time system metrics
  • Historical data tracking
  • Process monitoring
  • Network statistics
  • Alerting capabilities`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize history collector
		go collectHistoryWorker()

		// API Endpoints
		http.HandleFunc("/api/v1/metrics", handleMetrics)
		http.HandleFunc("/api/v1/metrics/history", handleMetricsHistory)
		http.HandleFunc("/api/v1/system/info", handleSystemInfo)
		http.HandleFunc("/api/v1/processes", handleProcesses)
		http.HandleFunc("/api/v1/health", handleHealthCheck)
		http.HandleFunc("/api/v1/network", handleNetworkStats)

		// Serve static dashboard
		fs := http.FileServer(http.Dir("./dashboard"))
		http.Handle("/", fs)

		// Start server
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		green := color.New(color.FgGreen).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()

		fmt.Printf("\n🚀 %s Vigil Metrics Dashboard\n", green("Starting"))
		fmt.Printf("📡 %s http://localhost:%d\n", cyan("Dashboard URL:"), port)
		fmt.Printf("📊 %s http://localhost:%d/api/v1/metrics\n", cyan("Live Metrics:"), port)
		fmt.Printf("📈 %s http://localhost:%d/api/v1/metrics/history\n", cyan("History API:"), port)
		fmt.Printf("💻 %s http://localhost:%d/api/v1/system/info\n", cyan("System Info:"), port)
		fmt.Printf("⚙️  %s http://localhost:%d/api/v1/processes\n", cyan("Process List:"), port)
		fmt.Printf(" %s\n\n", cyan("Use Ctrl+C to stop"))

		if err := http.ListenAndServe(addr, nil); err != nil {
			color.Red(" Failed to start server: %v", err)
			os.Exit(1)
		}
	},
}

func collectMetrics() map[string]interface{} {
	// CPU Metrics
	cpuPercent, _ := cpu.Percent(0, false)
	cpuCount, _ := cpu.Counts(true)
	cpuFreq, _ := cpu.Info()
	var cpuFreqStr string
	if len(cpuFreq) > 0 {
		cpuFreqStr = fmt.Sprintf("%.2f GHz", cpuFreq[0].Mhz/1000)
	}

	// Memory Metrics
	memInfo, _ := mem.VirtualMemory()
	swapInfo, _ := mem.SwapMemory()

	// Disk Metrics
	diskInfo, _ := disk.Usage("/")
	diskIO, _ := disk.IOCounters()

	// Network Metrics
	netIO, _ := net.IOCounters(false)
	var netStats map[string]interface{}
	if len(netIO) > 0 {
		netStats = map[string]interface{}{
			"bytes_sent":   netIO[0].BytesSent,
			"bytes_recv":   netIO[0].BytesRecv,
			"packets_sent": netIO[0].PacketsSent,
			"packets_recv": netIO[0].PacketsRecv,
			"err_in":       netIO[0].Errin,
			"err_out":      netIO[0].Errout,
			"drop_in":      netIO[0].Dropin,
			"drop_out":     netIO[0].Dropout,
		}
	}

	// Host Info
	hostInfo, _ := host.Info()
	uptime, _ := host.Uptime()

	// Process Count
	processes, _ := process.Processes()
	processCount := len(processes)

	// Go Runtime Metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"timestamp": time.Now().UTC(),
		"host": map[string]interface{}{
			"hostname":         hostInfo.Hostname,
			"os":               hostInfo.OS,
			"platform":         hostInfo.Platform,
			"platform_family":  hostInfo.PlatformFamily,
			"platform_version": hostInfo.PlatformVersion,
			"kernel_version":   hostInfo.KernelVersion,
			"uptime_seconds":   uptime,
		},
		"cpu": map[string]interface{}{
			"percent":        cpuPercent[0],
			"cores_physical": cpuCount,
			"cores_logical":  runtime.NumCPU(),
			"frequency":      cpuFreqStr,
			"load_average":   getLoadAverage(),
		},
		"memory": map[string]interface{}{
			"total":        memInfo.Total,
			"available":    memInfo.Available,
			"used":         memInfo.Used,
			"free":         memInfo.Free,
			"percent":      memInfo.UsedPercent,
			"swap_total":   swapInfo.Total,
			"swap_used":    swapInfo.Used,
			"swap_percent": swapInfo.UsedPercent,
			"cached":       memInfo.Cached,
			"buffers":      memInfo.Buffers,
		},
		"disk": map[string]interface{}{
			"total":          diskInfo.Total,
			"free":           diskInfo.Free,
			"used":           diskInfo.Used,
			"percent":        diskInfo.UsedPercent,
			"inodes_total":   diskInfo.InodesTotal,
			"inodes_used":    diskInfo.InodesUsed,
			"inodes_free":    diskInfo.InodesFree,
			"inodes_percent": diskInfo.InodesUsedPercent,
			"io_read_bytes":  getDiskIOBytes(diskIO, "read"),
			"io_write_bytes": getDiskIOBytes(diskIO, "write"),
		},
		"network": netStats,
		"system": map[string]interface{}{
			"goroutines":    runtime.NumGoroutine(),
			"cgo_calls":     runtime.NumCgoCall(),
			"process_count": processCount,
			"go_mem_alloc":  m.Alloc,
			"go_mem_sys":    m.Sys,
			"go_mem_heap":   m.HeapAlloc,
			"go_mem_stack":  m.StackInuse,
			"go_gc_count":   m.NumGC,
			"go_gc_pause":   m.PauseTotalNs,
		},
		"alerts": checkAlerts(cpuPercent[0], memInfo.UsedPercent, diskInfo.UsedPercent),
	}
}

func getLoadAverage() []float64 {
	avg, err := load.Avg()
	if err != nil {
		return []float64{0, 0, 0}
	}
	return []float64{avg.Load1, avg.Load5, avg.Load15}
}

func getDiskIOBytes(io map[string]disk.IOCountersStat, op string) uint64 {
	var total uint64
	for _, counter := range io {
		if op == "read" {
			total += counter.ReadBytes
		} else {
			total += counter.WriteBytes
		}
	}
	return total
}

func checkAlerts(cpu, mem, disk float64) []map[string]interface{} {
	var alerts []map[string]interface{}
	now := time.Now()

	if cpu > 90 {
		alerts = append(alerts, map[string]interface{}{
			"level":   "critical",
			"metric":  "cpu",
			"message": fmt.Sprintf("CPU usage critical: %.1f%%", cpu),
			"value":   cpu,
			"time":    now,
		})
	} else if cpu > 75 {
		alerts = append(alerts, map[string]interface{}{
			"level":   "warning",
			"metric":  "cpu",
			"message": fmt.Sprintf("CPU usage high: %.1f%%", cpu),
			"value":   cpu,
			"time":    now,
		})
	}

	if mem > 90 {
		alerts = append(alerts, map[string]interface{}{
			"level":   "critical",
			"metric":  "memory",
			"message": fmt.Sprintf("Memory usage critical: %.1f%%", mem),
			"value":   mem,
			"time":    now,
		})
	} else if mem > 80 {
		alerts = append(alerts, map[string]interface{}{
			"level":   "warning",
			"metric":  "memory",
			"message": fmt.Sprintf("Memory usage high: %.1f%%", mem),
			"value":   mem,
			"time":    now,
		})
	}

	if disk > 90 {
		alerts = append(alerts, map[string]interface{}{
			"level":   "critical",
			"metric":  "disk",
			"message": fmt.Sprintf("Disk usage critical: %.1f%%", disk),
			"value":   disk,
			"time":    now,
		})
	}

	return alerts
}

func collectHistoryWorker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics := collectMetrics()
		historyMutex.Lock()
		metricsHistory = append(metricsHistory, metrics)
		if len(metricsHistory) > historyLimit {
			metricsHistory = metricsHistory[len(metricsHistory)-historyLimit:]
		}
		historyMutex.Unlock()
	}
}

// API Handlers
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := collectMetrics()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	json.NewEncoder(w).Encode(metrics)
}

func handleMetricsHistory(w http.ResponseWriter, r *http.Request) {
	historyMutex.RLock()
	defer historyMutex.RUnlock()

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	if limit > len(metricsHistory) {
		limit = len(metricsHistory)
	}

	response := map[string]interface{}{
		"count":   limit,
		"history": metricsHistory[len(metricsHistory)-limit:],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	info, _ := host.Info()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func handleProcesses(w http.ResponseWriter, r *http.Request) {
	processes, err := process.Processes()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("Failed to get processes: %v", err),
		})
		return
	}

	var processList []map[string]interface{}

	for _, p := range processes {
		name, _ := p.Name()
		mem, _ := p.MemoryInfo()
		cpu, _ := p.CPUPercent()

		// Handle nil memory info
		var rss, vms uint64
		if mem != nil {
			rss = mem.RSS
			vms = mem.VMS
		}

		processList = append(processList, map[string]interface{}{
			"pid":     p.Pid,
			"name":    name,
			"cpu":     cpu,
			"mem_rss": rss,
			"mem_vms": vms,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processList)
}

func handleNetworkStats(w http.ResponseWriter, r *http.Request) {
	stats, _ := net.IOCounters(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"time":    time.Now().UTC(),
		"uptime":  time.Since(startTime).String(),
		"version": "1.0.0",
	})
}

var startTime = time.Now()

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port for dashboard")
	serveCmd.Flags().IntVar(&historyLimit, "history-limit", 1000, "Maximum number of historical data points")
	rootCmd.AddCommand(serveCmd)
}
