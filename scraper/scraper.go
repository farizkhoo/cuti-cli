package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Holiday struct {
	Date   string   `json:"date"`
	Day    string   `json:"day"`
	Name   string   `json:"name"`
	States []string `json:"states"`
}

type Scraper struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewScraper initializes chromedp with sensible defaults
func NewScraper(headless bool) *Scraper {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	// Block heavy resources
	_ = chromedp.Run(ctx,
		network.Enable(),
		network.SetBlockedURLs([]string{
			"*.png", "*.jpg", "*.jpeg", "*.gif",
			"*.woff", "*.ttf", "*.svg", "*.css",
		}),
	)

	return &Scraper{ctx: ctx, cancel: cancel}
}

func (s *Scraper) Close() {
	s.cancel()
}

// FetchState scrapes one state page (national excluded)
func (s *Scraper) FetchState(state string, year int) ([]Holiday, error) {
	url := buildURL(state, year)
	log.Printf("🌐 Fetching %s (%d) — %s", state, year, url)

	// per-page timeout
	ctx, cancel := context.WithTimeout(s.ctx, 20*time.Second)
	defer cancel()

	var rows [][]string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("table.publicholidays", chromedp.ByQuery),
		chromedp.Evaluate(fmt.Sprintf(`
			(() => {
				// Find the header for the requested year
				const yearHeader = Array.from(document.querySelectorAll("h2"))
					.find(h => h.innerText.includes("%d"));
				if (!yearHeader) return [];

				// Table immediately after the h2
				const table = yearHeader.nextElementSibling;
				if (!table || !table.classList.contains("publicholidays")) return [];

				const trs = Array.from(table.querySelectorAll("tbody tr"));
				return trs.map(tr => {
					const tds = Array.from(tr.querySelectorAll("td")).map(td => td.innerText.trim());
					return tds;
				});
			})()
		`, year), &rows),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading %s: %w", state, err)
	}

	if len(rows) == 0 {
		log.Printf("⚠️  No rows found for %s in %d; page may have changed", state, year)
		return nil, nil
	}

	var holidays []Holiday
	for _, r := range rows {
		if len(r) < 3 {
			continue
		}
		dateStr := normalizeDate(r[0], year)
		day := r[1]
		name := r[2]

		holidays = append(holidays, Holiday{
			Date:   dateStr,
			Day:    day,
			Name:   name,
			States: []string{normalizeState(state)},
		})
	}

	log.Printf("✅ Fetched %d rows for %s (%d)", len(holidays), state, year)
	return holidays, nil
}

func buildURL(state string, year int) string {
	// explicitly skip national
	return fmt.Sprintf("https://publicholidays.com.my/%s/%d-dates/", state, year)
}

func normalizeDate(dateStr string, year int) string {
	// Example: "1 Jan", "2 February", etc.
	layouts := []string{"2 Jan", "2 January"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return fmt.Sprintf("%d-%02d-%02d", year, int(t.Month()), t.Day())
		}
	}
	// fallback: keep as-is
	log.Printf("⚠️  Failed to parse date: %s", dateStr)
	return fmt.Sprintf("%d-%s", year, strings.ReplaceAll(dateStr, " ", "-"))
}

func normalizeState(st string) string {
	st = strings.ToLower(st)
	st = strings.ReplaceAll(st, " ", "-")
	st = strings.ReplaceAll(st, "&", "and")

	switch st {
	case "malacca":
		return "melaka"
	case "kualalumpur":
		return "kuala-lumpur"
	case "putrajayaand-selangor", "putrajaya-selangor":
		// handled elsewhere in consolidation
		return "putrajaya"
	}
	return st
}

func Consolidate(holidays []Holiday) []Holiday {
	merged := make(map[string]Holiday)

	for _, h := range holidays {
		// Key by date+name (ignore "day" since states may observe on diff days)
		key := h.Date + "|" + h.Name

		if existing, ok := merged[key]; ok {
			existing.States = append(existing.States, h.States...)
			existing.States = unique(existing.States)
			merged[key] = existing
		} else {
			h.States = unique(h.States)
			merged[key] = h
		}
	}

	// Convert back to slice
	result := make([]Holiday, 0, len(merged))
	for _, h := range merged {
		result = append(result, h)
	}

	// Sort by date for readability
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result
}

func unique(input []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range input {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// Save to JSON
func SaveJSON(path string, holidays []Holiday) error {
	data, err := json.MarshalIndent(holidays, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
