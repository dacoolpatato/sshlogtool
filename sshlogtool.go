// sshlogtool.go - Lightweight SSH login monitor & logger for Linux
// Parses /var/log/auth.log for SSH login events,
// tracks username, IP, method (password/publickey), and sudo command actions.

package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "time"
)

const (
    authLogPath   = "/var/log/auth.log"
    timeLayout    = "Jan 2 15:04:05"
)

type SSHLogin struct {
    Timestamp time.Time
    User      string
    IP        string
    Method    string
    Port      string
    Actions   []Action
}

type Action struct {
    Timestamp time.Time
    Type      string // read or changed
    Command   string
}

var (
    historyFlag = flag.Bool("history", false, "Show SSH login history")
    lastFlag    = flag.Bool("last", false, "Show last SSH login")
    watchFlag   = flag.Bool("watch", false, "Watch SSH logins live")
)

func main() {
    flag.Parse()

    if *historyFlag {
        entries, err := parseAuthLog()
        if err != nil {
            fmt.Println("Error parsing log:", err)
            os.Exit(1)
        }
        printHistory(entries)
        return
    }

    if *lastFlag {
        entries, err := parseAuthLog()
        if err != nil {
            fmt.Println("Error parsing log:", err)
            os.Exit(1)
        }
        if len(entries) == 0 {
            fmt.Println("No SSH logins found.")
            return
        }
        printEntry(entries[len(entries)-1])
        return
    }

    if *watchFlag {
        watchLogs()
        return
    }

    fmt.Println("Usage: sshlogtool [-history|-last|-watch]")
    os.Exit(1)
}

// parseAuthLog reads /var/log/auth.log and returns SSHLogin entries
func parseAuthLog() ([]SSHLogin, error) {
    file, err := os.Open(authLogPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var entries []SSHLogin
    var currentLogin *SSHLogin

    sshLoginRegex := regexp.MustCompile(`Accepted (\w+) for (\w+) from ([\d\.]+) port (\d+) ssh2`)
    sudoCmdRegex := regexp.MustCompile(`sudo: +(\w+) : .*COMMAND=(.+)`)

    now := time.Now()

    for scanner.Scan() {
        line := scanner.Text()
        if len(line) < 16 {
            continue
        }
        timestampStr := line[:15]
        ts, err := time.Parse(timeLayout, timestampStr)
        if err != nil {
            continue
        }
        ts = ts.AddDate(now.Year(), 0, 0)

        if sshLoginRegex.MatchString(line) {
            matches := sshLoginRegex.FindStringSubmatch(line)
            method := matches[1]
            user := matches[2]
            ip := matches[3]
            port := matches[4]

            entry := SSHLogin{
                Timestamp: ts,
                User:      user,
                IP:        ip,
                Method:    method,
                Port:      port,
                Actions:   []Action{},
            }
            entries = append(entries, entry)
            currentLogin = &entries[len(entries)-1]
        } else if sudoCmdRegex.MatchString(line) && currentLogin != nil {
            matches := sudoCmdRegex.FindStringSubmatch(line)
            user := matches[1]
            command := matches[2]

            if user == currentLogin.User {
                actionType := "unknown"
                lowerCmd := strings.ToLower(command)
                if strings.Contains(lowerCmd, "cat") || strings.Contains(lowerCmd, "less") || strings.Contains(lowerCmd, "grep") {
                    actionType = "read"
                } else if strings.Contains(lowerCmd, "nano") || strings.Contains(lowerCmd, "vim") || strings.Contains(lowerCmd, "echo") {
                    actionType = "changed"
                }

                act := Action{
                    Timestamp: ts,
                    Type:      actionType,
                    Command:   command,
                }
                currentLogin.Actions = append(currentLogin.Actions, act)
            }
        }
    }

    return entries, nil
}

func printHistory(entries []SSHLogin) {
    if len(entries) == 0 {
        fmt.Println("No SSH login history found.")
        return
    }
    for _, e := range entries {
        printEntry(e)
        fmt.Println()
    }
}

func printEntry(e SSHLogin) {
    fmt.Printf("[%s] SSH login detected:\n", e.Timestamp.Format("2006-01-02 15:04:05"))
    fmt.Printf("  • User: %s\n", e.User)
    fmt.Printf("  • Source IP: %s\n", e.IP)
    fmt.Printf("  • Method: %s\n", e.Method)
    fmt.Printf("  • Port: %s\n", e.Port)
    if len(e.Actions) == 0 {
        fmt.Printf("  • Actions: (none tracked)\n")
    } else {
        fmt.Printf("  • Actions:\n")
        for _, a := range e.Actions {
            icon := "?"
            if a.Type == "read" {
                icon = "📖"
            } else if a.Type == "changed" {
                icon = "✏️"
            }
            fmt.Printf("    → [%s] %s %s\n", a.Timestamp.Format("15:04:05"), icon, a.Command)
        }
    }
}

func watchLogs() {
    cmd := exec.Command("tail", "-F", authLogPath)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    if err := cmd.Start(); err != nil {
        fmt.Println("Error:", err)
        return
    }

    reader := bufio.NewReader(stdout)
    sshLoginRegex := regexp.MustCompile(`Accepted (\w+) for (\w+) from ([\d\.]+) port (\d+) ssh2`)

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                time.Sleep(time.Second)
                continue
            }
            fmt.Println("Error reading log:", err)
            break
        }
        line = strings.TrimSpace(line)
        if sshLoginRegex.MatchString(line) {
            matches := sshLoginRegex.FindStringSubmatch(line)
            method := matches[1]
            user := matches[2]
            ip := matches[3]
            port := matches[4]

            fmt.Printf("\n[NEW LOGIN] %s\n  User: %s\n  IP: %s\n  Method: %s\n  Port: %s\n\n",
                time.Now().Format("2006-01-02 15:04:05"),
                user,
                ip,
                method,
                port)
        }
    }
}
