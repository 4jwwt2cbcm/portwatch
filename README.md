# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and log file:

```bash
portwatch start --interval 60 --log /var/log/portwatch.log
```

Take a one-time snapshot of currently open ports:

```bash
portwatch scan
```

**Example output:**

```
[2024-01-15 08:32:11] INFO  Baseline snapshot taken: 12 open ports
[2024-01-15 08:33:11] INFO  Scan complete: no changes detected
[2024-01-15 08:34:11] ALERT New port opened: 0.0.0.0:8888 (PID 4821)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `30` | Scan interval in seconds |
| `--log` | stdout | Path to log file |
| `--alert-cmd` | — | Shell command to run on change |

---

## License

MIT © 2024 yourusername