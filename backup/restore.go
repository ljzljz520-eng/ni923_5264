package backup

import (
	"fmt"
	"sort"
	"strings"

	"example.com/toolnav/model"
)

type Manifest struct {
	SnapshotID string
	Revision   int
	ToolCount  int
	Checksum   string
	Categories []model.Category
}

type Difference struct {
	Missing []string
	Added   []string
	Moved   []string
}

func BuildManifest(snapshot model.BackupSnapshot, tools map[string]model.Tool) Manifest {
	categories := make(map[model.Category]bool)
	for _, id := range snapshot.ToolIDs {
		if tool, ok := tools[id]; ok {
			categories[tool.Category] = true
		}
	}
	values := make([]model.Category, 0, len(categories))
	for category := range categories {
		values = append(values, category)
	}
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	return Manifest{SnapshotID: snapshot.ID, Revision: snapshot.Revision, ToolCount: len(snapshot.ToolIDs), Checksum: snapshot.Checksum, Categories: values}
}

func (m Manifest) Valid() bool {
	return m.SnapshotID != "" && m.ToolCount >= 0 && m.Checksum != ""
}

func RestoreOrder(snapshot model.BackupSnapshot, tools map[string]model.Tool) (model.ToolOrder, error) {
	if err := VerifyBackup(snapshot, tools); err != nil {
		return model.ToolOrder{}, err
	}
	return model.ToolOrder{ToolIDs: append([]string(nil), snapshot.ToolIDs...), Revision: snapshot.Revision}, nil
}

func (s *Service) RestoreLatest(tools map[string]model.Tool) (model.ToolOrder, error) {
	snapshot, err := s.Latest()
	if err != nil {
		return model.ToolOrder{}, err
	}
	return RestoreOrder(snapshot, tools)
}

func (s *Service) PersistRestored(order model.ToolOrder) error {
	if s == nil || s.Store == nil {
		return fmt.Errorf("backup service is not configured")
	}
	return s.Store.SaveOrder(order)
}

func CompareSnapshots(previous, current model.BackupSnapshot) Difference {
	previousIndex := indexIDs(previous.ToolIDs)
	currentIndex := indexIDs(current.ToolIDs)
	difference := Difference{Missing: make([]string, 0), Added: make([]string, 0), Moved: make([]string, 0)}
	for id, index := range previousIndex {
		currentPosition, ok := currentIndex[id]
		if !ok {
			difference.Missing = append(difference.Missing, id)
		} else if currentPosition != index {
			difference.Moved = append(difference.Moved, id)
		}
	}
	for id := range currentIndex {
		if _, ok := previousIndex[id]; !ok {
			difference.Added = append(difference.Added, id)
		}
	}
	sort.Strings(difference.Missing)
	sort.Strings(difference.Added)
	sort.Strings(difference.Moved)
	return difference
}

func indexIDs(ids []string) map[string]int {
	result := make(map[string]int, len(ids))
	for index, id := range ids {
		result[id] = index
	}
	return result
}

func DuplicateIDs(ids []string) []string {
	seen := make(map[string]bool)
	duplicates := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			duplicates[id] = true
		}
		seen[id] = true
	}
	result := make([]string, 0, len(duplicates))
	for id := range duplicates {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}

func FormatDifference(difference Difference) string {
	parts := []string{
		"missing=" + strings.Join(difference.Missing, ","),
		"added=" + strings.Join(difference.Added, ","),
		"moved=" + strings.Join(difference.Moved, ","),
	}
	return strings.Join(parts, " ")
}
