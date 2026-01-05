#  vigil — Lightweight System Monitor for Terminals

> **Check CPU, Memory, Disk — and *profile any command* — in <10ms startup, <8MB binary.**
>
> Built for developers, CI/CD, Raspberry Pi, and air-gapped environments.

![vigil demo](https://github.com/yourname/vigil/raw/main/demo.gif  )

## ✨ Features
- `vigil cpu`, `mem`, `disk` — instant system snapshot
- `vigil exec -- <cmd>` — profile CPU/RAM/time of any process
- `--json` flag for scripting
- Cross-platform (Linux, macOS, Windows, ARM64!)
- Zero dependencies — single static binary

## 🚀 Install

### One-liner (Linux/macOS):
```bash
curl -sfL https://raw.githubusercontent.com/sahil3982/vigil/main/install.sh   | sh
```
### Windows

#### PowerShell (Recommended):
```powershell
iwr -useb https://raw.githubusercontent.com/sahil3982/vigil/main/install.ps1 | iex
```

### Or Git Bash
```bash
curl -sfL https://raw.githubusercontent.com/sahil3982/vigil/main/install.sh | sh
```

## 🖥️ Usage Examples

### Quick System Check
```bash
# Check CPU usage
$ vigil cpu
▶ CPU: [■■■■■■□□□□] 62.3% 

# Check memory usage
$ vigil mem
▶ RAM: [■■■■■■■□□□] 72.1% (11.5/16.0 GB) 

# Check disk usage
$ vigil disk
▶ Disk /: [■■■■■■■□□□] 72.1% (215.4/300.0 GB) 
```

### Profile Any Command
```bash
# Profile a build process
$ vigil exec -- go build main.go
▶ Running: go build main.go
──────────────────────────────────────
▶ Finished in 2.41s 
   CPU: avg 88%
   RAM: peak 1240.5 MB
   Exit code: 0

# Profile tests
$ vigil exec -- go test ./...
# See how much RAM your tests consume!

# Profile any process
$ vigil exec -- npm install
$ vigil exec -- docker build -t myapp .
```

### JSON Output for Automation
```bash
# Get structured data for scripts
$ vigil cpu --json
{"cpu_percent": 42.3}

$ vigil exec --json -- go test
{
  "command": "go test",
  "elapsed_seconds": 2.41,
  "cpu_avg_percent": 88.2,
  "ram_peak_mb": 1240.5,
  "exit_code": 0
}
```

## 🌐 Live Dashboard

Run a web-based dashboard to monitor your system in real-time:

```bash
# Start the dashboard server
vigil serve --port=3000
```

Then open: [http://localhost:3000](http://localhost:3000) in your browser.

### Remote Monitoring (Secure)
To monitor a remote server (e.g., AWS EC2, Raspberry Pi):
```bash
# On the remote server
vigil serve --port=3000

# On your local machine (SSH tunnel)
ssh -L 3000:localhost:3000 user@remote-server
```

Then open [http://localhost:3000](http://localhost:3000) to see the remote server's metrics securely.

## 🛠️ Why `vigil`?

### For Developers
- **Instant debugging**: "Why is my build slow?" → `vigil exec -- go build`
- **CI performance**: Catch memory leaks before they hit production
- **Cross-platform**: Works on your laptop, CI runners, and cloud servers

### For DevOps & System Admins
- **Zero setup**: No agents, no databases, no configuration
- **Lightweight**: <8MB binary, no dependencies
- **Scriptable**: JSON output for automation and alerting

### For Raspberry Pi & Edge Computing
- **Low overhead**: Minimal CPU/RAM usage
- **Self-contained**: Single binary, no network needed
- **Dashboard**: Visual monitoring without heavy tools

## 📜 License
MIT © [Sahil](https://github.com/sahil3982)

## 🤝 Contributing
PRs welcome! Check out the issues or suggest new features.

