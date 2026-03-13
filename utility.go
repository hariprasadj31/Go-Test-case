package main


import (
"strings"
)
func IsNullOrWhiteSpace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}