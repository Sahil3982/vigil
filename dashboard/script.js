// Configuration
const CONFIG = {
  api: {
    metrics: '/api/v1/metrics',
    history: '/api/v1/metrics/history',
    processes: '/api/v1/processes',
    system: '/api/v1/system/info',
    network: '/api/v1/network'
  },
  colors: {
    critical: '#f85149',
    warning: '#d29922',
    normal: '#238636',
    good: '#3fb950',
    info: '#58a6ff'
  },
  thresholds: {
    cpu: { warning: 75, critical: 90 },
    memory: { warning: 80, critical: 90 },
    disk: { warning: 85, critical: 90 }
  }
};

// State
let state = {
  isPaused: false,
  refreshInterval: 5000,
  lastUpdate: null,
  metricsHistory: [],
  networkStats: { prev: null, curr: null },
  cpuChart: null,
  memoryChart: null,
  gaugeChart: null
};

// DOM Elements
const elements = {
  timestamp: document.getElementById('timestamp'),
  hostname: document.getElementById('hostname'),
  uptime: document.getElementById('uptime'),
  alertPanel: document.getElementById('alertPanel'),
  alertList: document.getElementById('alertList'),
  alertCount: document.getElementById('alertCount'),
  refreshInterval: document.getElementById('refreshInterval'),
  pauseBtn: document.getElementById('pauseBtn'),
  processTableBody: document.getElementById('processTableBody'),
  refreshProcesses: document.getElementById('refreshProcesses')
};

// Initialize
document.addEventListener('DOMContentLoaded', () => {
  initializeCharts();
  loadSystemInfo();
  startMetricsPolling();
  setupEventListeners();
  loadProcesses();
});

// Event Listeners
function setupEventListeners() {
  elements.refreshInterval.addEventListener('change', (e) => {
    state.refreshInterval = parseInt(e.target.value);
  });

  elements.pauseBtn.addEventListener('click', () => {
    state.isPaused = !state.isPaused;
    elements.pauseBtn.innerHTML = state.isPaused ? 
      '<i class="fas fa-play"></i> Resume' : 
      '<i class="fas fa-pause"></i> Pause';
    elements.pauseBtn.classList.toggle('paused', state.isPaused);
  });

  elements.refreshProcesses.addEventListener('click', () => {
    loadProcesses();
  });

  document.getElementById('cpuHistoryRange').addEventListener('change', updateCharts);
  document.getElementById('memoryHistoryRange').addEventListener('change', updateCharts);
}

// Metrics Polling
function startMetricsPolling() {
  updateMetrics();
  setInterval(() => {
    if (!state.isPaused) {
      updateMetrics();
    }
  }, 1000);
}

// Update Metrics
async function updateMetrics() {
  try {
    const response = await fetch(CONFIG.api.metrics);
    const data = await response.json();
    
    state.lastUpdate = new Date();
    updateDashboard(data);
    updateAlerts(data.alerts);
    updateChartsData(data);
    
    // Update network rate calculations
    updateNetworkStats(data.network);
    
  } catch (error) {
    console.error('Failed to fetch metrics:', error);
    showError('Failed to fetch metrics. Check server connection.');
  }
}

// Update Dashboard
function updateDashboard(data) {
  // Timestamp
  elements.timestamp.textContent = `Last Updated: ${new Date(data.timestamp).toLocaleTimeString()}`;
  
  // Host Info
  elements.hostname.textContent = data.host.hostname;
  elements.uptime.textContent = `Uptime: ${formatUptime(data.host.uptime_seconds)}`;
  
  // CPU
  const cpuPct = data.cpu.percent;
  document.getElementById('cpu-value').textContent = `${cpuPct.toFixed(1)}%`;
  document.getElementById('cpu-cores').textContent = `Cores: ${data.cpu.cores_physical}`;
  document.getElementById('cpu-freq').textContent = `Freq: ${data.cpu.frequency}`;
  document.getElementById('load-1').textContent = data.cpu.load_average[0].toFixed(2);
  document.getElementById('load-5').textContent = data.cpu.load_average[1].toFixed(2);
  document.getElementById('load-15').textContent = data.cpu.load_average[2].toFixed(2);
  
  updateGauge('cpuGauge', cpuPct);
  
  // Memory
  const memPct = data.memory.percent;
  const memUsed = data.memory.used / 1e9;
  const memTotal = data.memory.total / 1e9;
  const memAvailable = data.memory.available / 1e9;
  const memCached = data.memory.cached / 1e9;
  
  document.getElementById('mem-value').textContent = `${memPct.toFixed(1)}%`;
  document.getElementById('mem-bar').style.width = `${memPct}%`;
  document.getElementById('mem-bar').style.background = getColorForMetric('memory', memPct);
  document.getElementById('mem-total').textContent = `Total: ${memTotal.toFixed(1)} GB`;
  document.getElementById('mem-used').textContent = `${memUsed.toFixed(1)} GB`;
  document.getElementById('mem-available').textContent = `${memAvailable.toFixed(1)} GB`;
  document.getElementById('mem-cached').textContent = `${memCached.toFixed(1)} GB`;
  document.getElementById('mem-swap').textContent = `${data.memory.swap_percent.toFixed(1)}%`;
  
  // Disk
  const diskPct = data.disk.percent;
  const diskUsed = data.disk.used / 1e9;
  const diskTotal = data.disk.total / 1e9;
  const diskFree = data.disk.free / 1e9;
  
  document.getElementById('disk-value').textContent = `${diskPct.toFixed(1)}%`;
  document.getElementById('disk-bar').style.width = `${diskPct}%`;
  document.getElementById('disk-bar').style.background = getColorForMetric('disk', diskPct);
  document.getElementById('disk-used').textContent = `${diskUsed.toFixed(1)} GB`;
  document.getElementById('disk-free').textContent = `${diskFree.toFixed(1)} GB`;
  document.getElementById('disk-total').textContent = `${diskTotal.toFixed(1)} GB`;
  document.getElementById('disk-inodes').textContent = `${data.disk.inodes_percent.toFixed(1)}%`;
  
  // Network
  if (data.network) {
    const txRate = calculateRate('tx', data.network.bytes_sent);
    const rxRate = calculateRate('rx', data.network.bytes_recv);
    
    document.getElementById('network-tx').textContent = `${txRate.toFixed(2)} MB/s`;
    document.getElementById('network-rx').textContent = `${rxRate.toFixed(2)} MB/s`;
    document.getElementById('network-packets-sent').textContent = data.network.packets_sent.toLocaleString();
    document.getElementById('network-packets-recv').textContent = data.network.packets_recv.toLocaleString();
    document.getElementById('network-errors').textContent = data.network.err_in + data.network.err_out;
    document.getElementById('network-drops').textContent = data.network.drop_in + data.network.drop_out;
  }
  
  // System Info
  document.getElementById('system-os').textContent = `${data.host.os} ${data.host.platform_version}`;
  document.getElementById('system-kernel').textContent = data.host.kernel_version;
  document.getElementById('system-platform').textContent = data.host.platform;
  document.getElementById('system-arch').textContent = data.host.platform_family;
  
  // Go Runtime
  document.getElementById('go-goroutines').textContent = data.system.goroutines.toLocaleString();
  document.getElementById('go-cgo').textContent = data.system.cgo_calls.toLocaleString();
  document.getElementById('go-heap').textContent = (data.system.go_mem_heap / 1e6).toFixed(1);
  document.getElementById('go-stack').textContent = (data.system.go_mem_stack / 1e6).toFixed(1);
  document.getElementById('go-gc-count').textContent = data.system.go_gc_count;
  document.getElementById('go-processes').textContent = data.system.process_count;
}

// Network Rate Calculation
function updateNetworkStats(current) {
  if (!state.networkStats.prev) {
    state.networkStats.prev = { ...current, timestamp: Date.now() };
    state.networkStats.curr = { ...current, timestamp: Date.now() };
    return;
  }
  
  state.networkStats.prev = state.networkStats.curr;
  state.networkStats.curr = { ...current, timestamp: Date.now() };
}

function calculateRate(type, currentBytes) {
  if (!state.networkStats.prev) return 0;
  
  const prev = state.networkStats.prev;
  const curr = state.networkStats.curr;
  const bytesKey = type === 'tx' ? 'bytes_sent' : 'bytes_recv';
  
  const bytesDiff = currentBytes - prev[bytesKey];
  const timeDiff = (curr.timestamp - prev.timestamp) / 1000; // in seconds
  
  return timeDiff > 0 ? (bytesDiff / timeDiff / 1e6) : 0; // MB/s
}

// Alerts
function updateAlerts(alerts) {
  if (!alerts || alerts.length === 0) {
    elements.alertList.innerHTML = '<div class="no-alerts">No active alerts</div>';
    elements.alertCount.textContent = '0';
    return;
  }
  
  elements.alertCount.textContent = alerts.length;
  elements.alertList.innerHTML = alerts.map(alert => `
    <div class="alert-item ${alert.level}">
      <div>
        <strong>${alert.level.toUpperCase()}</strong>: ${alert.message}
        <div class="alert-time">${new Date(alert.time).toLocaleTimeString()}</div>
      </div>
      <div class="alert-value">${alert.value.toFixed(1)}%</div>
    </div>
  `).join('');
}

// Charts
function initializeCharts() {
  // CPU Gauge
  const cpuGaugeCtx = document.getElementById('cpuGauge').getContext('2d');
  state.gaugeChart = new Chart(cpuGaugeCtx, {
    type: 'doughnut',
    data: {
      datasets: [{
        data: [0, 100],
        backgroundColor: [CONFIG.colors.normal, '#30363d'],
        borderWidth: 0,
        circumference: 180,
        rotation: 270
      }]
    },
    options: {
      cutout: '80%',
      responsive: false,
      plugins: {
        legend: { display: false },
        tooltip: { enabled: false }
      }
    }
  });
  
  // CPU History Chart
  const cpuHistoryCtx = document.getElementById('cpuHistoryChart').getContext('2d');
  state.cpuChart = new Chart(cpuHistoryCtx, {
    type: 'line',
    data: {
      labels: [],
      datasets: [{
        label: 'CPU Usage %',
        data: [],
        borderColor: CONFIG.colors.info,
        backgroundColor: `${CONFIG.colors.info}20`,
        fill: true,
        tension: 0.4,
        pointRadius: 0
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        x: {
          grid: { color: '#30363d' },
          ticks: { color: '#8b949e' }
        },
        y: {
          beginAtZero: true,
          max: 100,
          grid: { color: '#30363d' },
          ticks: { color: '#8b949e' }
        }
      },
      plugins: {
        legend: { display: false }
      }
    }
  });
  
  // Memory History Chart
  const memoryHistoryCtx = document.getElementById('memoryHistoryChart').getContext('2d');
  state.memoryChart = new Chart(memoryHistoryCtx, {
    type: 'line',
    data: {
      labels: [],
      datasets: [{
        label: 'Memory Usage %',
        data: [],
        borderColor: CONFIG.colors.good,
        backgroundColor: `${CONFIG.colors.good}20`,
        fill: true,
        tension: 0.4,
        pointRadius: 0
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        x: {
          grid: { color: '#30363d' },
          ticks: { color: '#8b949e' }
        },
        y: {
          beginAtZero: true,
          max: 100,
          grid: { color: '#30363d' },
          ticks: { color: '#8b949e' }
        }
      },
      plugins: {
        legend: { display: false }
      }
    }
  });
}

function updateChartsData(data) {
  const now = new Date(data.timestamp);
  const timeLabel = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  
  // Update CPU gauge
  state.gaugeChart.data.datasets[0].data = [data.cpu.percent, 100 - data.cpu.percent];
  state.gaugeChart.data.datasets[0].backgroundColor = [
    getColorForMetric('cpu', data.cpu.percent),
    '#30363d'
  ];
  state.gaugeChart.update('none');
  
  // Add to history
  state.metricsHistory.push({
    timestamp: data.timestamp,
    cpu: data.cpu.percent,
    memory: data.memory.percent,
    disk: data.disk.percent
  });
  
  // Keep last 100 points
  if (state.metricsHistory.length > 100) {
    state.metricsHistory.shift();
  }
  
  updateCharts();
}

function updateCharts() {
  const cpuRange = parseInt(document.getElementById('cpuHistoryRange').value);
  const memoryRange = parseInt(document.getElementById('memoryHistoryRange').value);
  
  const cpuData = filterHistoryByMinutes(state.metricsHistory, cpuRange);
  const memoryData = filterHistoryByMinutes(state.metricsHistory, memoryRange);
  
  // Update CPU chart
  state.cpuChart.data.labels = cpuData.map(d => 
    new Date(d.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  );
  state.cpuChart.data.datasets[0].data = cpuData.map(d => d.cpu);
  state.cpuChart.update();
  
  // Update Memory chart
  state.memoryChart.data.labels = memoryData.map(d => 
    new Date(d.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  );
  state.memoryChart.data.datasets[0].data = memoryData.map(d => d.memory);
  state.memoryChart.update();
}

function filterHistoryByMinutes(history, minutes) {
  const cutoff = new Date(Date.now() - minutes * 60000);
  return history.filter(h => new Date(h.timestamp) > cutoff);
}

// Processes
async function loadProcesses() {
  try {
    const response = await fetch(CONFIG.api.processes);
    const processes = await response.json();
    
    // Sort by CPU usage and take top 10
    const topProcesses = processes
      .sort((a, b) => b.cpu - a.cpu)
      .slice(0, 10);
    
    elements.processTableBody.innerHTML = topProcesses.map(proc => `
      <tr>
        <td>${proc.pid}</td>
        <td>${proc.name || 'Unknown'}</td>
        <td>
          <div class="progress-bar-small">
            <div class="progress-fill-small" style="width: ${Math.min(proc.cpu, 100)}%"></div>
          </div>
          <span>${proc.cpu.toFixed(1)}%</span>
        </td>
        <td>${(proc.mem_rss / 1e6).toFixed(1)} MB</td>
        <td>${(proc.mem_rss / 1e9).toFixed(2)} GB</td>
        <td>${(proc.mem_vms / 1e9).toFixed(2)} GB</td>
      </tr>
    `).join('');
    
  } catch (error) {
    console.error('Failed to load processes:', error);
    elements.processTableBody.innerHTML = '<tr><td colspan="6">Failed to load processes</td></tr>';
  }
}

// System Info
async function loadSystemInfo() {
  try {
    const response = await fetch(CONFIG.api.system);
    const info = await response.json();
    
    document.getElementById('hostname').textContent = info.hostname;
    document.getElementById('system-platform').textContent = info.platform;
    document.getElementById('system-arch').textContent = info.kernelArch;
    
  } catch (error) {
    console.error('Failed to load system info:', error);
  }
}

// Utility Functions
function getColorForMetric(metric, value) {
  const threshold = CONFIG.thresholds[metric];
  if (!threshold) return CONFIG.colors.normal;
  
  if (value >= threshold.critical) return CONFIG.colors.critical;
  if (value >= threshold.warning) return CONFIG.colors.warning;
  return CONFIG.colors.good;
}

function formatUptime(seconds) {
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  
  if (days > 0) return `${days}d ${hours}h`;
  if (hours > 0) return `${hours}h ${mins}m`;
  return `${mins}m`;
}

function formatBytes(bytes) {
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let value = bytes;
  let unitIndex = 0;
  
  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024;
    unitIndex++;
  }
  
  return `${value.toFixed(1)} ${units[unitIndex]}`;
}

function showError(message) {
  const errorDiv = document.createElement('div');
  errorDiv.className = 'alert-item critical';
  errorDiv.innerHTML = `
    <div>
      <strong>ERROR</strong>: ${message}
    </div>
  `;
  
  elements.alertList.prepend(errorDiv);
  
  // Auto remove after 5 seconds
  setTimeout(() => {
    errorDiv.remove();
  }, 5000);
}

// Add this CSS for small progress bars
const style = document.createElement('style');
style.textContent = `
  .progress-bar-small {
    height: 6px;
    background: #30363d;
    border-radius: 3px;
    margin-bottom: 3px;
    overflow: hidden;
  }
  
  .progress-fill-small {
    height: 100%;
    border-radius: 3px;
    background: linear-gradient(90deg, #238636, #58a6ff);
  }
`;
document.head.appendChild(style);