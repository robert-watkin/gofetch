package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
)

// run executes the fetch logic.
// It returns the exit code and any error.
// We use os.Stderr directly here for CLI output (diagnostics always go to stderr).
func run(opts options) (int, error) {
	resp, err := fetch(context.Background(), opts)
	if err != nil {
		// Check for HTTP error (4xx/5xx) first — this should exit with code 1
		var httpErr *HttpError
		if errors.As(err, &httpErr) {
			return 1, err
		}
		// Everything else (network, timeout, DNS, etc.) is code 2
		return 2, err
	}
	defer resp.Body.Close()

	// --json validation (do this before printing body)
	if opts.json {
		contentType := resp.Header.Get("Content-Type")
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return 1, fmt.Errorf("failed to parse Content-Type: %w", err)
		}
		if mediaType != "application/json" {
			return 1, fmt.Errorf("response body was not JSON (got %s)", mediaType)
		}
	}

	// Print body to stdout
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return 2, err
	}

	return 0, nil
}

// fetch performs the HTTP GET.
// It returns a response only for 2xx statuses.
// Non-2xx and transport errors are returned as typed errors.
func fetch(ctx context.Context, opts options) (*http.Response, error) {
	// Proper context + timeout setup.
	// Always capture the cancel func and defer it.
	ctx, cancel := context.WithTimeout(ctx, opts.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, opts.url, nil)
	if err != nil {
		return nil, &NetworkError{URL: opts.url, Err: err}
	}

	// Add headers (this is done after creating req, before Do)
	for _, strHeader := range opts.headers {
		parts := strings.SplitN(strHeader, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format (must be Name: value): %s", strHeader)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			return nil, fmt.Errorf("invalid header (empty key or value): %s", strHeader)
		}
		req.Header.Add(key, value)
	}

	// Print outgoing request info for --verbose (before actually sending)
	if opts.verbose {
		fmt.Fprintf(os.Stderr, "> %s %s HTTP/1.1\n", req.Method, req.URL.RequestURI())
		for k, vs := range req.Header {
			for _, v := range vs {
				fmt.Fprintf(os.Stderr, "> %s: %s\n", k, v)
			}
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &NetworkError{URL: opts.url, Err: err}
	}

	// --verbose: response info goes to stderr
	if opts.verbose {
		fmt.Fprintf(os.Stderr, "< %s\n", resp.Status)
		for key, values := range resp.Header {
			for _, value := range values {
				fmt.Fprintf(os.Stderr, "< %s: %s\n", key, value)
			}
		}
	}

	// Treat non-success as an error (using our custom type)
	if resp.StatusCode >= 400 {
		// We close the body here because the caller won't get the resp.
		resp.Body.Close()
		return nil, &HttpError{
			URL:        opts.url,
			StatusCode: resp.StatusCode,
			Err:        fmt.Errorf("HTTP %d", resp.StatusCode),
		}
	}

	return resp, nil
}
