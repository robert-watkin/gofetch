package main

import (
	"errors"
	"fmt"
)

var (
	ErrNetwork = errors.New("network error")
	ErrHttp    = errors.New("http error")
)

type NetworkError struct {
	URL string
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("network error fetching %s: %v", e.URL, e.Err)
}
func (e *NetworkError) Is(target error) bool { return target == ErrNetwork }
func (e *NetworkError) Unwrap() error        { return e.Err }

type HttpError struct {
	URL        string
	StatusCode int
	Err        error
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("HTTP error fetching %s: %s (status %d)", e.URL, e.Err, e.StatusCode)
}
func (e *HttpError) Is(target error) bool { return target == ErrHttp }
func (e *HttpError) Unwrap() error        { return e.Err }
