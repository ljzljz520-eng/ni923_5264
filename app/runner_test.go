package app

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunnerCommands(t *testing.T) {
	var output bytes.Buffer
	runner, err := New(t.TempDir()+"/app.db", &output)
	if err != nil {
		t.Fatal(err)
	}
	defer runner.Close()
	if err := runner.Run([]string{"sample"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "all rows accepted") {
		t.Fatal("sample command did not render import result")
	}
	output.Reset()
	if err := runner.Run([]string{"list", "category=payments"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "payrail") {
		t.Fatal("list command did not render category")
	}
}
