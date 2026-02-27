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

> **Note:** On WSL, you don't need a display server. Run the tool with `-headless=true` (the default for unattended runs). No GUI is required.

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
go run main.go -format json -out holidays -year 2025
