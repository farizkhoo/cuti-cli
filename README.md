# cuti-cli
CLI to parse Malaysia Holidays

# Prerequisites

This tool uses a headless Chrome browser to scrape holiday data. You need Google Chrome installed.

## Installing Chrome on WSL (Ubuntu)

```sh
# Download and install Google Chrome
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo apt install -y ./google-chrome-stable_current_amd64.deb

# Verify installation
google-chrome --version
```

> **Note:** On WSL, you don't need a display server. Pass `-headless=true` for unattended runs (the flag defaults to `false`, which opens a visible Chrome window).

# Usage

An overview of the available flags:

| Flag      | Description                        | Default    |
|-----------|------------------------------------|------------|
| `-year`     | Year to fetch holidays for         | `2025`     |
| `-format`   | Output format: `json` or `csv`     | `json`     |
| `-out`      | Output file name without extension | `holidays` |
| `-headless` | Run Chrome in headless mode        | `false`    |

## Example

```sh
go run main.go -format json -out holidays -year 2025 -headless=true
```

Output is written to `<out>-<year>.<format>`, e.g. `holidays-2025.json`.
