package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var app application

	err := app.setup()
	if err != nil {
		return err
	}

	err = app.startServer()
	if err != nil {
		return err
	}

	return nil
}
