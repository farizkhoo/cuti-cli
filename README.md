# cuti-cli
CLI to parse Malaysia Holidays

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
