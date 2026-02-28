//go:build headless

package main

import "errors"

// run is a no-op stub used in headless (CI) environments.
func run(_ Config) error {
	return errors.New("windowed rendering not available in headless build")
}
