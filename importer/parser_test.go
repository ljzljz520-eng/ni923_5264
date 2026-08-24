package importer

import (
	"testing"

	"example.com/toolnav/model"
)

func TestParserBuildBatch(t *testing.T) {
	parser := New("fixture")
	batch, tools := parser.BuildBatch("batch", []string{ExampleRows()[0], "bad|row"})
	if batch.Rows != 2 || batch.Accepted != 1 || batch.Rejected != 1 || len(tools) != 1 {
		t.Fatalf("unexpected batch: %+v", batch)
	}
	if NormalizeRow(" a | b ") != "a|b" || Slug("Pay Rail") != "pay-rail" {
		t.Fatal("normalization failed")
	}
	tool, err := BuildTool("New Tool", "New Tool", model.CategoryPayments, "https://new.example", model.StatusActive, []string{"Bills", "bills"})
	if err != nil || tool.ID != "new-tool" || len(tool.Tags) != 1 {
		t.Fatalf("unexpected built tool: %+v %v", tool, err)
	}
}
