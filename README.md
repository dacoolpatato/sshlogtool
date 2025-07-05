# sshlogtool

**sshlogtool** is a CLI-based Linux tool that monitors and logs all SSH login attempts and sudo activity. It captures the username, IP address, login method (password or key), timestamps, and command actions — labeling them as either `READ` or `CHANGED` based on behavior.

Built in Go and packaged as a `.deb`, it’s lightweight, dependency-free, and easy to run or install as a service. For personal use, sysadmins, or curious Linux users who want full visibility.

---

## Table of Contents

- [Overview](#overview)
- [Use Cases](#use-cases)
- [Features](#features)
- [Output Example](#output-example)
- [Installation (.deb)](#installation-deb)
- [Usage](#usage)
- [Build from Source](#build-from-source)
- [System Requirements](#system-requirements)
- [Configuration Notes](#configuration-notes)
- [License](#license)
- [Author](#author)

---

## Overview

sshlogtool scans your system's authentication log (`/var/log/auth.log`) to track who accessed your Linux machine via SSH, what authentication method they used, and whether they performed any sudo-level commands.

Each command is analyzed and categorized:
- `READ` – Passive (e.g. `cat`, `ls`, `grep`)
- `CHANGED` – Modifying (e.g. `apt`, `nano`, `rm`, `systemctl`)

---

## Use Cases

- Monitor your server for unauthorized SSH access
- Get detailed records of user activity during sessions
- Audit system usage without needing heavy logging setups
- Run it manually, or install it as a background systemd service

---

## Features

- SSH login tracking with:
  - Timestamp
  - Username
  - IP address
  - Auth method (Password or Key)
- Sudo command tracking per session
- Command classification: `READ` or `CHANGED`
- Live monitoring or historical lookup
- CLI-based, no GUI, low resource use
- `.deb` installer for `amd64` and `arm64` platforms
- Optional systemd service for continuous logging
- No dependencies beyond Go standard library

---

## Output Example

[2025-07-04 16:01] LOGIN from 192.168.1.24 via Password - USER: dario
→ SUDO: apt update (CHANGED)
→ SUDO: cat /etc/passwd (READ)

[2025-07-04 19:13] LOGIN from 86.1.X.X via SSH Key - USER: pi
→ No sudo usage

yaml
Copy
Edit

---

## Installation (.deb)

If using the `.deb` installer:

```bash
sudo dpkg -i sshlogtool_1.0_amd64.deb
sudo systemctl enable sshlogtool
sudo systemctl start sshlogtool
For ARM-based systems (e.g. Raspberry Pi):

bash
Copy
Edit
sudo dpkg -i sshlogtool_1.0_arm64.deb
Usage
After installing, run from the terminal:

bash
Copy
Edit
sshlogtool -history   # Display full SSH login and sudo history
sshlogtool -last      # Show the most recent login session
sshlogtool -watch     # Monitor live for new logins and actions
Use sudo if permission is denied:

bash
Copy
Edit
sudo sshlogtool -history
Build from Source
Requires Go (https://go.dev/dl)

Clone the repo and build:
bash
Copy
Edit
git clone https://github.com/YOUR_USERNAME/sshlogtool.git
cd sshlogtool
go build -o sshlogtool sshlogtool.go
Cross-compile for Windows (optional):
bash
Copy
Edit
GOOS=windows GOARCH=amd64 go build -o sshlogtool.exe sshlogtool.go
Build .deb Packages
Run the provided build script:

bash
Copy
Edit
chmod +x build_sshlogtool.sh
./build_sshlogtool.sh
Output .deb files will appear in the debbuild/ folder.

System Requirements
Debian-based Linux (Ubuntu, Kali, Raspbian, etc.)

Access to /var/log/auth.log (default on most systems)

Go (only needed to build from source)

For non-Debian distros, you may need to modify the log path in source code.

Configuration Notes
sshlogtool uses /var/log/auth.log for SSH and sudo data

SSH must be running: check with sudo systemctl status ssh

Your distro must log sudo and authentication properly (most do by default)

If you want to track different logs or customize behavior, you can edit and recompile the Go source.

