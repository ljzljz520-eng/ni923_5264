package catalog

import (
	"testing"

	"example.com/toolnav/importer"
	"example.com/toolnav/store"
)

func testCatalog(t *testing.T) *Service {
	t.Helper()
	s, err := store.Open(t.TempDir() + "/catalog.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	c, err := NewService(s)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestImportFilterAndReorder(t *testing.T) {
	catalog := testCatalog(t)
	batch, err := catalog.ImportRows("batch", "fixture", importer.ExampleRows())
	if err != nil || batch.Accepted != 5 {
		t.Fatalf("import failed: %+v %v", batch, err)
	}
	filtered := catalog.Filter(FilterOptions{Tag: "METRICS", MinRank: 0, MaxRank: 4})
	if len(filtered) != 1 || filtered[0].ID != "pulsewatch" {
		t.Fatalf("unexpected filtered tools: %+v", filtered)
	}
	if err := catalog.ReorderByIDs([]string{"policykit", "docforge", "payrail", "pulsewatch", "deployctl"}); err != nil {
		t.Fatal(err)
	}
	if catalog.OrderedTools()[0].ID != "policykit" {
		t.Fatal("reorder did not update ordered view")
	}
	if err := catalog.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(CategoriesPresent(catalog.OrderedTools())) != 5 {
		t.Fatal("category coverage changed")
	}
	if len(ValidateView(catalog.OrderedTools(), catalog.Order)) != 0 {
		t.Fatal("ordered view is invalid")
	}
}
