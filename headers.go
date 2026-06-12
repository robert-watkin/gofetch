package main

import (
	"fmt"
	"strings"
)

type headerFlags []string

func (h headerFlags) String() string {
	return strings.Join(h, ",")
}

func (h *headerFlags) Set(value string) error {
	if !strings.Contains(value, ":") {
		return fmt.Errorf("header %q must be in \"Name: value\" form", value)
	}
	*h = append(*h, value)
	return nil
}
