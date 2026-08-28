package app

import (
	"fmt"
	"io"
	"os"

	toolnav "example.com/toolnav"
)

type Runner struct {
	Navigator *toolnav.Navigator
	Out       io.Writer
}

func New(path string, out io.Writer) (*Runner, error) {
	navigator, err := toolnav.Open(path)
	if err != nil {
		return nil, err
	}
	if out == nil {
		out = os.Stdout
	}
	return &Runner{Navigator: navigator, Out: out}, nil
}

func (r *Runner) Run(args []string) error {
	output, err := r.Navigator.Execute(args)
	if err != nil {
		return err
	}
	if output != "" {
		_, err = fmt.Fprintln(r.Out, output)
	}
	return err
}

func (r *Runner) Close() error { return r.Navigator.Close() }
