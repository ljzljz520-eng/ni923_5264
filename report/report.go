package report

import (
	"fmt"
	"strings"

	"example.com/toolnav/backup"
	"example.com/toolnav/model"
)

func RenderCatalog(tools []model.Tool) string {
	lines := make([]string, 0, len(tools))
	for _, tool := range tools {
		lines = append(lines, fmt.Sprintf("%02d %s [%s] %s %s", tool.Rank, tool.ID, tool.Category, tool.Status, tool.URL))
	}
	return strings.Join(lines, "\n")
}

func RenderBackup(snapshot model.BackupSnapshot) string {
	return fmt.Sprintf("backup=%s revision=%d checksum=%s order=%s", snapshot.ID, snapshot.Revision, snapshot.Checksum, backup.OrderedLabels(snapshot))
}

func RenderImport(batch model.ImportBatch) string {
	return fmt.Sprintf("import=%s source=%s rows=%d accepted=%d rejected=%d", batch.ID, batch.Source, batch.Rows, batch.Accepted, batch.Rejected)
}

func RenderRecovery(tools map[string]model.Tool, order model.ToolOrder, backups int) string {
	return fmt.Sprintf("recovered tools=%d order=%d revision=%d backups=%d", len(tools), len(order.ToolIDs), order.Revision, backups)
}

func RenderCategories() string {
	categories := model.Categories()
	labels := make([]string, len(categories))
	for i, category := range categories {
		labels[i] = string(category)
	}
	return strings.Join(labels, ",")
}
