package main

import (
	"flag"
	"fmt"
	"os"

	"example.com/toolnav/app"
)

func main() {
	path := flag.String("db", "toolnav.db", "path to the bbolt database")
	flag.Parse()
	runner, err := app.New(*path, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer runner.Close()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stdout, app.Help())
		return
	}
	if err := runner.Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
