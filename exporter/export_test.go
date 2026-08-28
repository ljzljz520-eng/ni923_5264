package exporter

import (
	"strings"
	"testing"

	"example.com/toolnav/model"
)

func TestExporterRoundTrip(t *testing.T) {
	tools := []model.Tool{{ID: "a", Name: "A", Category: model.CategoryDeployment, URL: "https://a.example", Status: model.StatusActive, Rank: 0}}
	bundle := NewBundle(map[string]model.Tool{"a": tools[0]}, model.ToolOrder{ToolIDs: []string{"a"}, Revision: 1}, nil, nil)
	data, err := JSON(bundle)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseJSON(data)
	if err != nil || len(parsed.Tools) != 1 {
		t.Fatalf("round trip failed: %+v %v", parsed, err)
	}
	if !strings.Contains(CSV(tools, bundle.Order), "id,name") {
		t.Fatal("CSV header missing")
	}
}
