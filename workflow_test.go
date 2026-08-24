package toolnav

import (
	"testing"

	"example.com/toolnav/backup"
	"example.com/toolnav/governance"
	"example.com/toolnav/importer"
	"example.com/toolnav/model"
	"example.com/toolnav/workflow"
)

func openTestNavigator(t *testing.T) *Navigator {
	t.Helper()
	navigator, err := Open(t.TempDir() + "/tools.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = navigator.Close() })
	return navigator
}

func importSample(t *testing.T, navigator *Navigator) {
	t.Helper()
	batch, err := navigator.Import("batch-1", "fixture", importer.ExampleRows())
	if err != nil {
		t.Fatal(err)
	}
	if !batch.Complete() || batch.Accepted != 5 {
		t.Fatalf("unexpected import batch: %+v", batch)
	}
}

func TestWorkflowImportQueryStatus(t *testing.T) {
	navigator := openTestNavigator(t)
	result, err := navigator.ImportAndReview("batch-flow", "fixture", importer.ExampleRows(), model.DefaultPolicy())
	if err != nil {
		t.Fatal(err)
	}
	if err := workflow.ValidateResult(result); err != nil {
		t.Fatal(err)
	}
	if result.Batch.Accepted != 5 || len(result.Tools) != 5 {
		t.Fatalf("unexpected result: %+v", result.Batch)
	}
	query := workflow.Query(navigator.Catalog, model.CategoryMonitoring, "", "")
	if query.Total != 1 || query.Tools[0].ID != "pulsewatch" {
		t.Fatalf("unexpected query: %+v", query)
	}
	decision, err := navigator.UpdateStatus("pulsewatch", model.StatusActive, model.DefaultPolicy())
	if err != nil || !decision.Allowed {
		t.Fatalf("status update failed: %+v %v", decision, err)
	}
	dashboard, err := navigator.Dashboard(model.DefaultPolicy())
	if err != nil || dashboard.ToolCount != 5 || dashboard.CategoryCoverage != 5 {
		t.Fatalf("unexpected dashboard: %+v %v", dashboard, err)
	}
}

func TestWorkflowSortAndBackup(t *testing.T) {
	navigator := openTestNavigator(t)
	importSample(t, navigator)
	result, err := workflow.SortAndBackup(navigator.Catalog, navigator.Backup, "policykit", 0, "after-sort")
	if err != nil {
		t.Fatal(err)
	}
	if result.Order.ToolIDs[0] != "policykit" {
		t.Fatalf("unexpected order: %+v", result.Order)
	}
	if err := backup.VerifyBackup(result.Snapshot, navigator.Catalog.Tools); err != nil {
		t.Fatal(err)
	}
	if !result.Manifest.Valid() || result.Manifest.ToolCount != 5 {
		t.Fatalf("unexpected manifest: %+v", result.Manifest)
	}
}

func TestWorkflowRestoreAndReopen(t *testing.T) {
	navigator := openTestNavigator(t)
	importSample(t, navigator)
	initial := append([]string(nil), navigator.Catalog.Order.ToolIDs...)
	if _, err := navigator.CreateBackup("before-change"); err != nil {
		t.Fatal(err)
	}
	if err := navigator.Move("deployctl", 4, nil); err != nil {
		t.Fatal(err)
	}
	restored, err := navigator.RestoreLatest()
	if err != nil {
		t.Fatal(err)
	}
	if len(restored.ToolIDs) != len(initial) {
		t.Fatalf("unexpected restored order: %+v", restored)
	}
	for index := range initial {
		if initial[index] != restored.ToolIDs[index] {
			t.Fatalf("restored order differs: initial=%v restored=%v", initial, restored.ToolIDs)
		}
	}
	if len(governance.Coverage(navigator.Catalog.OrderedTools())) != 5 {
		t.Fatal("restored catalog lost category coverage")
	}
}

func TestToolBackupSeesConsistentOrder(t *testing.T) {
	navigator := openTestNavigator(t)
	importSample(t, navigator)
	var snapshot model.BackupSnapshot
	var backupErr error
	err := navigator.Move("deployctl", 4, func() {
		snapshot, backupErr = navigator.CreateBackup("during-drag")
	})
	if err != nil {
		t.Fatal(err)
	}
	if backupErr != nil {
		t.Fatalf("backup during drag failed: %v", backupErr)
	}
	if err := backup.VerifyBackup(snapshot, navigator.Catalog.Tools); err != nil {
		t.Fatalf("backup has inconsistent order: %v", err)
	}
}
