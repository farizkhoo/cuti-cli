package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/farizkhoo/cuti-cli/scraper"
)

func main() {
	year := flag.Int("year", 2025, "Year to fetch holidays for")
	format := flag.String("format", "json", "Output format: json or csv")
	out := flag.String("out", "holidays", "Output file name without extension")
	headless := flag.Bool("headless", false, "Run Chrome in headless mode")
	flag.Parse()

	// States + national
	states := []string{
		"national",
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

	switch strings.ToLower(*format) {
	case "json":
		filename := *out + ".json"
		if err := scraper.SaveJSON(filename, final); err != nil {
			log.Fatal(err)
		}
		log.Printf("✅ Holidays written to %s", filename)

	case "csv":
		filename := *out + ".csv"
		if err := saveCSV(filename, final); err != nil {
			log.Fatal(err)
		}
		log.Printf("✅ Holidays written to %s", filename)

	default:
		log.Fatalf("Unsupported format: %s", *format)
	}
}

func saveCSV(path string, holidays []scraper.Holiday) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Header
	if err := w.Write([]string{"Date", "Day", "Name", "States"}); err != nil {
		return err
	}

	for _, h := range holidays {
		row := []string{h.Date, h.Day, h.Name, strings.Join(h.States, ";")}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}
