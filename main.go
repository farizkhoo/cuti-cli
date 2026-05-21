package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/farizkhoo/cuti-cli/scraper"
)

func main() {
	year := flag.Int("year", 2025, "Year to fetch holidays for")
	format := flag.String("format", "json", "Output format: json or csv")
	out := flag.String("out", "holidays", "Output file name without extension")
	headless := flag.Bool("headless", false, "Run Chrome in headless mode")
	flag.Parse()

	normalizedFormat := strings.ToLower(*format)
	if normalizedFormat != "json" && normalizedFormat != "csv" {
		log.Fatalf("Unsupported format: %s (expected json or csv)", *format)
	}

	// States only (national excluded)
	states := []string{
		"johor", "kedah", "kelantan", "kuala-lumpur",
		"labuan", "melaka", "negeri-sembilan", "pahang",
		"penang", "perak", "perlis", "putrajaya",
		"sabah", "sarawak", "selangor", "terengganu",
	}

	s := scraper.NewScraper(*headless)
	defer s.Close()

	var all []scraper.Holiday
	for i, st := range states {
		log.Printf("🌐 [%d/%d] Fetching %s (%d)…", i+1, len(states), st, *year)

		holidays, err := s.FetchState(st, *year)
		if err != nil {
			log.Printf("⛔ Failed to fetch %s (%d): %v", st, *year, err)
			continue
		}
		all = append(all, holidays...)
	}

	final := scraper.Consolidate(all)

	filename := fmt.Sprintf("%s-%d.%s", *out, *year, normalizedFormat)
	var saveErr error
	switch normalizedFormat {
	case "json":
		saveErr = scraper.SaveJSON(filename, final)
	case "csv":
		saveErr = scraper.SaveCSV(filename, final)
	}
	if saveErr != nil {
		log.Fatal(saveErr)
	}
	log.Printf("✅ Holidays written to %s", filename)
}

