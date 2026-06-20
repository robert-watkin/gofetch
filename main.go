package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type options struct {
	url     string
	headers headerFlags
	timeout time.Duration
	json    bool
	verbose bool
}

func parseArgs() (options, error) {
	var opts options

	fs := flag.NewFlagSet("gofetch", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Var(&opts.headers, "header", "Multiple optional headers")
	fs.BoolVar(&opts.json, "json", false, "Process response body as JSON")
	fs.BoolVar(&opts.verbose, "verbose", false, "Output request and response headers")
	fs.DurationVar(&opts.timeout, "timeout", 30*time.Second, "request timeout")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return opts, fmt.Errorf("Failed to parse options: %v", err)
	}

	if fs.NArg() != 1 {
		return opts, errors.New("usage: gofetch [flags] URL")
	}
	opts.url = fs.Arg(0)

	return opts, nil
}

func main() {
	// struct to hold the options returned
	opts, err := parseArgs()
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}

	code, err := run(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gofetch: %v\n", err)
	}
	os.Exit(code)
}
