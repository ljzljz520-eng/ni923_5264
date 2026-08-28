package audit

import (
	"testing"

	"example.com/toolnav/model"
)

func TestEventFilteringAndValidation(t *testing.T) {
	events := []AuditEvent{
		NewEvent("e2", EventMoved, "admin", []string{"b"}, 2, "move"),
		NewEvent("e1", EventImported, "admin", []string{"a", "b"}, 1, "import"),
	}
	if errors := ValidateForTools(events, map[string]model.Tool{"a": {}, "b": {}}); len(errors) != 0 {
		t.Fatal(errors)
	}
	filtered := Filter(events, EventMoved, "admin")
	if len(filtered) != 1 || filtered[0].ID != "e2" {
		t.Fatalf("unexpected filtered events: %+v", filtered)
	}
	if CountByType(events)[EventImported] != 1 || len(ToolIDs(events)) != 2 {
		t.Fatal("event summary is incomplete")
	}
}
