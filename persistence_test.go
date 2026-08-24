package toolnav

import (
	"testing"

	"example.com/toolnav/importer"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/persistent.db"
	navigator, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := navigator.Import("batch-persist", "fixture", importer.ExampleRows()); err != nil {
		t.Fatal(err)
	}
	if _, err := navigator.CreateBackup("persisted"); err != nil {
		t.Fatal(err)
	}
	if err := navigator.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	state, err := reopened.Store.LoadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	if len(state.Tools) != 5 || len(state.Order.ToolIDs) != 5 || len(state.Imports) != 1 || len(state.Backups) != 1 || len(state.Audits) < 2 {
		t.Fatalf("persistence state was not restored: %+v", state)
	}
}
