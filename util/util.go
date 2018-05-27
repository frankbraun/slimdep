// Package util contains utility functions.
package util

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
)

// GetGOPATH returns the Go path (environment variable or default).
func GetGOPATH() (string, error) {
	path := os.Getenv("GOPATH")
	if path == "" {
		path = os.Getenv("HOME")
		if path != "" {
			path = filepath.Join(path, "go")
		}
	}
	if path == "" {
		return "", errors.New("cannot determine GOPATH")
	}
	return path, nil
}

// UniqueStrings sorts the given string array a and makes the content unique.
func UniqueStrings(a []string) []string {
	var r []string
	sort.Strings(a)
	if len(a) > 0 {
		r = append(r, a[0])
		for i := 1; i < len(a); i++ {
			if a[i] != r[len(r)-1] {
				r = append(r, a[i])
			}
		}
	}
	return r
}
