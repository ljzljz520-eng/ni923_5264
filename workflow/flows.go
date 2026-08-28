package workflow

import (
	"fmt"
	"strings"

	"example.com/toolnav/audit"
	"example.com/toolnav/backup"
	"example.com/toolnav/catalog"
	"example.com/toolnav/governance"
	"example.com/toolnav/model"
)

type ImportResult struct {
	Batch  model.ImportBatch
	Tools  []model.Tool
	Review []governance.Review
	Audit  audit.AuditEvent
}

type QueryResult struct {
	Tools      []model.Tool
	Total      int
	Categories map[model.Category]int
}

type SortBackupResult struct {
	Order    model.ToolOrder
	Snapshot model.BackupSnapshot
	Manifest backup.Manifest
}

func ImportAndReview(c *catalog.Service, batchID, source string, rows []string, policy model.Policy) (ImportResult, error) {
	if c == nil {
		return ImportResult{}, fmt.Errorf("catalog is required")
	}
	batch, err := c.ImportRows(batchID, source, rows)
	if err != nil {
		return ImportResult{}, err
	}
	tools := c.OrderedTools()
	reviews := governance.BuildReview(tools, policy)
	events, err := c.Audit.EventsFor(audit.EventImported)
	if err != nil {
		return ImportResult{}, err
	}
	var event audit.AuditEvent
	if len(events) > 0 {
		event = events[len(events)-1]
	}
	return ImportResult{Batch: batch, Tools: tools, Review: reviews, Audit: event}, nil
}

func Query(c *catalog.Service, category model.Category, status model.Status, term string) QueryResult {
	var tools []model.Tool
	if strings.TrimSpace(term) != "" {
		tools = c.Search(term)
	} else {
		tools = c.Query(category, status)
	}
	counts := make(map[model.Category]int)
	for _, tool := range tools {
		counts[tool.Category]++
	}
	return QueryResult{Tools: append([]model.Tool(nil), tools...), Total: len(tools), Categories: counts}
}

func SortAndBackup(c *catalog.Service, b *backup.Service, id string, target int, label string) (SortBackupResult, error) {
	if c == nil || b == nil {
		return SortBackupResult{}, fmt.Errorf("catalog and backup services are required")
	}
	if err := c.Move(id, target, nil); err != nil {
		return SortBackupResult{}, err
	}
	snapshot, err := b.Create(c.Order, c.Tools, label)
	if err != nil {
		return SortBackupResult{}, err
	}
	return SortBackupResult{Order: c.Order.Clone(), Snapshot: snapshot, Manifest: backup.BuildManifest(snapshot, c.Tools)}, nil
}

func Restore(c *catalog.Service, b *backup.Service) (model.ToolOrder, error) {
	if c == nil || b == nil {
		return model.ToolOrder{}, fmt.Errorf("catalog and backup services are required")
	}
	order, err := b.RestoreLatest(c.Tools)
	if err != nil {
		return model.ToolOrder{}, err
	}
	if err := c.ReorderByIDs(order.ToolIDs); err != nil {
		return model.ToolOrder{}, err
	}
	return c.Order.Clone(), nil
}

func StatusReview(c *catalog.Service, policy model.Policy) []governance.Review {
	if c == nil {
		return nil
	}
	return governance.BuildReview(c.OrderedTools(), policy)
}

func WorkflowNames() []string {
	return []string{"import-review", "query-report", "sort-backup", "restore-reopen"}
}

func ValidateResult(result ImportResult) error {
	if !result.Batch.Complete() {
		return fmt.Errorf("import batch counts do not balance")
	}
	if len(result.Tools) < result.Batch.Accepted {
		return fmt.Errorf("import result has too few tools")
	}
	return nil
}
