package backup

import (
	"testing"

	"example.com/toolnav/model"
)

func backupTools() map[string]model.Tool {
	return map[string]model.Tool{
		"a": {ID: "a", Name: "A", Category: model.CategoryDeployment, URL: "https://a.example", Status: model.StatusActive},
		"b": {ID: "b", Name: "B", Category: model.CategoryMonitoring, URL: "https://b.example", Status: model.StatusActive},
	}
}

func TestVerifyAndCompareSnapshots(t *testing.T) {
	tools := backupTools()
	first := model.BackupSnapshot{ID: "one", Revision: 1, ToolIDs: []string{"a", "b"}, Checksum: Checksum([]string{"a", "b"})}
	if err := VerifyBackup(first, tools); err != nil {
		t.Fatal(err)
	}
	second := model.BackupSnapshot{ID: "two", Revision: 2, ToolIDs: []string{"b", "a"}, Checksum: Checksum([]string{"b", "a"})}
	difference := CompareSnapshots(first, second)
	if len(difference.Moved) != 2 || len(DuplicateIDs(second.ToolIDs)) != 0 {
		t.Fatalf("unexpected difference: %+v", difference)
	}
	if _, err := RestoreOrder(second, tools); err != nil {
		t.Fatal(err)
	}
}
