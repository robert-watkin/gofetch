package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRun(test *testing.T) {
	// setup server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	fmt.Printf("[verbose] Test server started: %v\n", ts.URL)

	var headerF headerFlags
	timeoutDuration := 30 * time.Second // create a time.Duration var of 30 seconds

	// loop through the tests
	var tests = []struct {
		tName       string
		tTimeout    time.Duration
		tJson       bool
		tStatusCode int
	}{
		{"Test 1", 30 * time.Second, false, 0},
		{"Test 2", 0 * time.Second, false, 0},
	}

	for _, t := range tests {
		test.Run(t.tName, func(test *testing.T) {
			// setup options for test
			opts := options{
				url:     ts.URL,
				headers: headerF,
				timeout: timeoutDuration,
				json:    false,
				verbose: false,
			}

			fmt.Printf("[verbose] %v running...\n", t.tName)

			// expect this to pass
			statusCode, err := run(opts)
			fmt.Printf("[verbose] checking err (expecting nil):\n%v\n", err)
			if err != nil {
				test.Errorf("run() err expected nil, got %v", err)
			}

			fmt.Printf("[verbose] checking statusCode (expecting %d):\n%d\n", t.tStatusCode, statusCode)
			if statusCode != t.tStatusCode {
				test.Errorf("run() response code expected 0, got %v", statusCode)
			}
		})
	}
}
