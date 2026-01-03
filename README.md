# 🕵️ vigil — Lightweight System Monitor for Terminals

> **Check CPU, Memory, Disk — and *profile any command* — in <10ms startup, <8MB binary.**
>
> Built for developers, CI/CD, Raspberry Pi, and air-gapped environments.

![vigil demo](https://github.com/yourname/vigil/raw/main/demo.gif)

## ✨ Features
- `vigil cpu`, `mem`, `disk` — instant system snapshot
- `vigil exec -- <cmd>` — profile CPU/RAM/time of any process
- `--json` flag for scripting
- Cross-platform (Linux, macOS, Windows, ARM64!)
- Zero dependencies — single static binary

## 🚀 Install

### One-liner (Linux/macOS):
```bash
curl -sfL https://raw.githubusercontent.com/sahil3982/vigil/main/install.sh | sh


## 🚀 Install on Windows

### PowerShell (Recommended)
```powershell
iwr -useb https://raw.githubusercontent.com/sahil3982/vigil/main/install.ps1 | iex