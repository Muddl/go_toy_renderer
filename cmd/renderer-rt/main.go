package main

import (
	"fmt"
	"os"
	"runtime"
)

func init() {
	// GLFW requires all calls on the OS thread.
	runtime.LockOSThread()
}

func main() {
	cfg, err := parseConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err = validateBackend(cfg.Backend); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if runErr := run(cfg); runErr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", runErr)
		os.Exit(1)
	}
}
