# crontrace

> Wrapper around cron jobs that logs execution history, duration, and exit codes to a local SQLite database for auditing.

---

## Installation

```bash
go install github.com/yourusername/crontrace@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/crontrace.git
cd crontrace
go build -o crontrace .
```

---

## Usage

Wrap any command by prepending `crontrace run --` to your cron entry:

```
# crontab -e
* * * * * crontrace run -- /usr/local/bin/backup.sh --quiet
```

Each execution is recorded in a local SQLite database (`~/.crontrace/history.db` by default).

**View execution history:**

```bash
crontrace list
```

**Inspect a specific job:**

```bash
crontrace show --job backup.sh --limit 20
```

**Example output:**

```
JOB              STARTED              DURATION   EXIT CODE
backup.sh        2024-06-01 03:00:01  1.243s     0
backup.sh        2024-05-31 03:00:01  0.987s     0
backup.sh        2024-05-30 03:00:02  2.105s     1
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--db` | `~/.crontrace/history.db` | Path to the SQLite database |
| `--label` | command name | Human-readable label for the job |
| `--timeout` | none | Kill the job after N seconds |

---

## License

MIT © [yourusername](https://github.com/yourusername)