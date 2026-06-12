package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"
)

func run(opts options) (int, error) {
	req, err := http.NewRequest("GET", opts.url, nil)
	if err != nil {
		return 2, err
	}

	for _, strHeader := range opts.headers {
		splitHeader := strings.SplitN(strHeader, ":", 2)
		key := splitHeader[0]
		value := splitHeader[1]

		if key == "" || value == "" {
			log.Fatalf("Header %v is invalid", strHeader)
		}

		req.Header.Add(key, value)
	}

	if opts.verbose {
		fmt.Fprintf(os.Stderr, "> GET %v HTTP/1.1\n", opts.url)

		for _, header := range opts.headers {
			fmt.Fprintf(os.Stderr, "> %v\n", header)
		}
	}

	client := &http.Client{Timeout: opts.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return 2, err
	}
	defer resp.Body.Close()

	if opts.json {
		contentType := resp.Header.Get("Content-Type")

		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return 2, err
		}

		if resp.StatusCode >= 400 || mediaType != "application/json" {
			return 1, fmt.Errorf("Response body was not JSON")
		}
	}

	io.Copy(os.Stdout, resp.Body)

	if opts.verbose {
		fmt.Fprintf(os.Stderr, "< %v\n", resp.Status)
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Fprintf(os.Stderr, "< %v:%v\n", key, value)
			}
		}
	}

	return 0, nil
}
