package backup

import (
	"fmt"
	"sort"

	"example.com/toolnav/model"
)

func CompareOrder(snapshot model.BackupSnapshot, order model.ToolOrder) error {
	if len(snapshot.ToolIDs) != len(order.ToolIDs) {
		return fmt.Errorf("snapshot and order lengths differ")
	}
	for index := range order.ToolIDs {
		if snapshot.ToolIDs[index] != order.ToolIDs[index] {
			return fmt.Errorf("order differs at position %d", index)
		}
	}
	return nil
}

func MissingIDs(snapshot model.BackupSnapshot, tools map[string]model.Tool) []string {
	seen := make(map[string]bool, len(snapshot.ToolIDs))
	for _, id := range snapshot.ToolIDs {
		seen[id] = true
	}
	missing := make([]string, 0)
	for id := range tools {
		if !seen[id] {
			missing = append(missing, id)
		}
	}
	sort.Strings(missing)
	return missing
}
