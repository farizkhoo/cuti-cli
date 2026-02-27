# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this project does

`cuti-cli` is a Go CLI tool that scrapes Malaysia public holidays from [publicholidays.com.my](https://publicholidays.com.my) for all 16 states using a headless Chrome browser (chromedp), then outputs the consolidated data as JSON or CSV.

## Commands

```sh
# Run the scraper (default: json, year 2025, headless=false)
go run main.go -format json -out holidays -year 2025

# Run with visible Chrome window (for debugging)
go run main.go -headless=false

# Run in headless mode
go run main.go -headless=true

# Build binary
go build -o cuti-cli .

# Tidy dependencies
go mod tidy
```

There are no tests in this project currently.

## Architecture

- **`main.go`** — Entry point. Parses CLI flags (`-year`, `-format`, `-out`, `-headless`), iterates over all 16 Malaysian states, calls the scraper for each, then consolidates and writes output.
- **`scraper/scraper.go`** — All scraping logic. Key types and functions:
  - `Holiday` — struct with `Date`, `Day`, `Name`, `States []string`
  - `NewScraper(headless bool)` — initializes a single shared chromedp browser context with resource blocking (images, fonts, CSS)
  - `FetchState(state, year)` — navigates to `https://publicholidays.com.my/{state}/{year}-dates/`, waits for `.publicholidays` table, extracts rows via JS evaluation
  - `Consolidate([]Holiday)` — merges holidays with the same date+name across states into a single entry with combined `States` slice, then sorts by date
  - `SaveJSON(path, holidays)` — writes indented JSON output

## Key behaviors

- The scraper uses a **single shared browser context** across all state fetches (one `chromedp.NewContext` call, reused per state).
- Each `FetchState` call has a **20-second per-page timeout**.
- Consolidation key is `date|name` — holidays with the same name on the same date across different states are merged.
- The `-headless` flag defaults to `false` (visible browser window), which is unusual — set to `true` for unattended runs.
- JSON output filename includes the year: `{out}-{year}.json`. CSV output does not: `{out}.csv`.
