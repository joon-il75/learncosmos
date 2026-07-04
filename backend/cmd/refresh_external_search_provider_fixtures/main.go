package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/learnweaver/backend/internal/evaluation/lessonsearchprerun"
)

func main() {
	opts := parseFlags()
	lessonsearchprerun.LoadEnv()
	if err := lessonsearchprerun.Run(context.Background(), opts, lessonsearchprerun.BuildDefaultProviders()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func parseFlags() lessonsearchprerun.Options {
	var opts lessonsearchprerun.Options
	flag.StringVar(&opts.FixturesPath, "fixtures", lessonsearchprerun.DefaultFixturePath, "input provider fixture JSON path")
	flag.StringVar(&opts.OutputPath, "out", lessonsearchprerun.DefaultOutputPath, "output review report JSON path")
	flag.StringVar(&opts.Provider, "provider", "all", "provider filter: youtube, naver_blog, all")
	flag.StringVar(&opts.FixtureID, "fixture-id", "", "limit refresh to a single fixture id")
	flag.IntVar(&opts.Limit, "limit", 0, "maximum fixtures to process; 0 means all")
	flag.IntVar(&opts.TopN, "top-n", 5, "maximum provider candidates per fixture")
	flag.StringVar(&opts.Mode, "mode", lessonsearchprerun.DefaultMode, "mode: refresh-existing, discover-gap")
	flag.DurationVar(&opts.Timeout, "timeout", 10*time.Second, "provider request timeout")
	flag.Parse()
	return opts
}
