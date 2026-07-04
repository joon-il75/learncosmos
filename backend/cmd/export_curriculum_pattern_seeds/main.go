package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/learnweaver/backend/internal/domain/curriculum"
)

func main() {
	outPath := flag.String("out", "", "output JSON file path")
	version := flag.String("version", "v1", "pattern version")
	language := flag.String("language", "all", "seed language: ko, en, or all")
	flag.Parse()

	if *outPath == "" {
		fmt.Fprintln(os.Stderr, "missing required -out")
		os.Exit(2)
	}

	seeds, err := curriculum.BuildCurriculumPatternSeedsForLanguage(*version, *language)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build seeds: %v\n", err)
		os.Exit(1)
	}
	if len(seeds) == 0 {
		fmt.Fprintln(os.Stderr, "no seeds generated")
		os.Exit(1)
	}

	payload, err := json.MarshalIndent(seeds, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal seeds: %v\n", err)
		os.Exit(1)
	}
	payload = append(payload, '\n')

	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create output dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outPath, payload, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("wrote %d curriculum pattern seeds to %s\n", len(seeds), *outPath)
}
