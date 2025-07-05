# sshlogtool

sshlogtool is a CLI-based Linux tool that monitors and logs all SSH login attempts and sudo activity. It captures the username, IP address, login method (password or key), timestamps, and command actions — labeling them as either READ or CHANGED based on behavior.

Built in Go and packaged as a .deb, it’s lightweight, dependency-free, and easy to run or install as a service. For personal use, sysadmins, or curious Linux users who want full visibility.

---

## Table of Contents

- Overview
- Use Cases
- Features
- Output Example
- Installation (.deb)
- Usage
- Build from Source
- System Requirements
- Configuration Notes
- License
- Author

---

## Overview

sshlogtool scans your system's authentication log (/var/log/auth.log) to track who accessed your Linux machine via SSH, what authentication method they used, and whether they performed any sudo-level commands.

Each command is analyzed and categorized:
- READ – Passive (e.g. cat, ls, grep)
- CHANGED – Modifying (e.g. apt, nano, rm, systemctl)

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
- Command classification: READ or CHANGED
- Live monitoring or historical lookup
- CLI-based, no GUI, low resource use
- .deb installer for amd64 and arm64 platforms
- Optional systemd service for continuous logging
- No dependencies beyond Go standard library

---

## Output Example

[2025-07-04 16:01] LOGIN from 192.168.1.24 via Password - USER: dario
   → SUDO: apt update                 (CHANGED)
   → SUDO: cat /etc/passwd           (READ)

[2025-07-04 19:13] LOGIN from 86.1.X.X via SSH Key - USER: pi
   → No sudo usage

---

## Installation (.deb)

If using the .deb installer:

    sudo dpkg -i sshlogtool_1.0_amd64.deb
    sudo systemctl enable sshlogtool
    sudo systemctl start sshlogtool

For ARM-based systems (e.g. Raspberry Pi):

    sudo dpkg -i sshlogtool_1.0_arm64.deb

---

## Usage

After installing, run from the terminal:

    sshlogtool -history   # Display full SSH login and sudo history
    sshlogtool -last      # Show the most recent login session
    sshlogtool -watch     # Monitor live for new logins and actions

Use sudo if permission is denied:

    sudo sshlogtool -history

---

## Build from Source

Requires Go (https://go.dev/dl)

Clone the repo and build:

    git clone https://github.com/YOUR_USERNAME/sshlogtool.git
    cd sshlogtool
    go build -o sshlogtool sshlogtool.go

Optional: cross-compile for Windows:

    GOOS=windows GOARCH=amd64 go build -o sshlogtool.exe sshlogtool.go

---

## Build .deb Packages

Run the provided build script:

    chmod +x build_sshlogtool.sh
    ./build_sshlogtool.sh

Output .deb files will appear in the debbuild/ folder.

---

## System Requirements

- Debian-based Linux (Ubuntu, Kali, Raspbian, etc.)
- Access to /var/log/auth.log (default on most systems)
- Go (only needed to build from source)

For non-Debian distros, you may need to modify the log path in source code.

---

## Configuration Notes

- sshlogtool uses /var/log/auth.log for SSH and sudo data
- SSH must be running: check with sudo systemctl status ssh
- Your distro must log sudo and authentication properly (most do by default)

If you want to track different logs or customize behavior, you can edit and recompile the Go source.

---

## License

Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International

You may:
- Use this tool for personal or educational purposes
- Share the tool with credit to the original author (Dario)

You may not:
- Modify, fork, or rebrand the tool
- Sell or use it in any commercial product
- Remove copyright or license

Full license text: https://creativecommons.org/licenses/by-nc-nd/4.0/

---

## Author

Created by Dario  
2025 – Personal security tools and system utilities  
If you find this tool useful, star the project or open an issue with feedback.
