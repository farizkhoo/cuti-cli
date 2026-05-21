---
name: test
description: Verify cuti-cli compiles, vets clean, and (optionally) runs end-to-end against publicholidays.com.my. Use when the user says "test this", "run the tests", or wants confidence a change didn't break the scraper.
---

# test

This project has no `go test` suite yet, so "testing" means a layered check from cheapest to most thorough. Run each layer in order and stop reporting at the first failure — the user wants to know *what broke*, not a wall of green ticks.

## Layer 1 — static checks (always run)

```sh
go vet ./...
go build -o /tmp/cuti-cli-build-check . && rm /tmp/cuti-cli-build-check
```

Both must succeed. `go build` to a throwaway path avoids polluting the working tree.

## Layer 2 — unit tests (run if any `*_test.go` exists)

```sh
go test ./...
```

If there are no test files, say so explicitly — don't claim "tests passed" when nothing ran.

## Layer 3 — smoke run against the live site (only if the user asks for end-to-end / "real" testing, or if scraper logic changed)

Scraping all 16 states takes a few minutes. For a fast smoke, edit nothing — just check that one state fetches cleanly by running the full binary with headless on:

```sh
go run main.go -headless=true -year=2025 -out=/tmp/cuti-smoke
```

Watch the logs:
- `✅ Fetched N rows for <state> (2025)` for every state → scraper is healthy
- `⚠️  No rows found for <state>` on most states → page structure likely changed (selectors in `scraper/scraper.go` need updating)
- `⛔ Failed to fetch …` on most states → network, Chrome, or Cloudflare issue

Then sanity-check the output:

```sh
jq 'length, .[0]' /tmp/cuti-smoke-2025.json
```

Expect ~15–25 unique holidays and entries shaped like `{date, day, name, states}`. Clean up `/tmp/cuti-smoke-2025.json` afterwards.

## Reporting

Tell the user which layers ran and which were skipped, and why. Don't run layer 3 unsolicited — it hits a third-party site and takes minutes.
