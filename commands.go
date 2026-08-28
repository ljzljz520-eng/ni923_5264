package toolnav

import (
	"fmt"
	"strconv"
	"strings"

	"example.com/toolnav/audit"
	"example.com/toolnav/exporter"
	"example.com/toolnav/importer"
	"example.com/toolnav/model"
	"example.com/toolnav/report"
)

func (n *Navigator) Execute(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("command is required")
	}
	switch args[0] {
	case "sample":
		batch, err := n.Import("batch-sample", "sample", importer.ExampleRows())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s\n%s", batch.Summary(), n.CatalogReport()), nil
	case "list":
		category, status := parseFilters(args[1:])
		return n.ReportFiltered(category, status), nil
	case "search":
		if len(args) < 2 {
			return "", fmt.Errorf("search term is required")
		}
		return n.ReportTools(n.Search(strings.Join(args[1:], " "))), nil
	case "move":
		if len(args) != 3 {
			return "", fmt.Errorf("move requires id and position")
		}
		position, err := strconv.Atoi(args[2])
		if err != nil {
			return "", fmt.Errorf("invalid position: %w", err)
		}
		if err := n.Move(args[1], position, nil); err != nil {
			return "", err
		}
		return n.CatalogReport(), nil
	case "backup":
		label := "manual"
		if len(args) > 1 {
			label = strings.Join(args[1:], "-")
		}
		snapshot, err := n.CreateBackup(label)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("backup %s", snapshot.ID), nil
	case "restore":
		order, err := n.RestoreLatest()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("restored revision=%d order=%s", order.Revision, strings.Join(order.ToolIDs, ">")), nil
	case "status":
		if len(args) != 3 {
			return "", fmt.Errorf("status requires id and status")
		}
		decision, err := n.UpdateStatus(args[1], model.Status(args[2]), model.DefaultPolicy())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("status %s -> %s: %s", decision.From, decision.To, decision.Reason), nil
	case "dashboard":
		dashboard, err := n.Dashboard(model.DefaultPolicy())
		if err != nil {
			return "", err
		}
		return dashboard.Format(), nil
	case "audit":
		events, err := n.Catalog.Audit.Events()
		if err != nil {
			return "", err
		}
		lines := make([]string, 0, len(events))
		for _, event := range events {
			lines = append(lines, audit.Describe(event))
		}
		return strings.Join(lines, "\n"), nil
	case "export":
		bundle, err := n.ExportBundle()
		if err != nil {
			return "", err
		}
		data, err := exporter.JSON(bundle)
		if err != nil {
			return "", err
		}
		return string(data), nil
	default:
		return "", fmt.Errorf("unknown command %q", args[0])
	}
}

func parseFilters(args []string) (model.Category, model.Status) {
	var category model.Category
	var status model.Status
	for _, arg := range args {
		if strings.HasPrefix(arg, "category=") {
			category = model.Category(strings.TrimPrefix(arg, "category="))
		}
		if strings.HasPrefix(arg, "status=") {
			status = model.Status(strings.TrimPrefix(arg, "status="))
		}
	}
	return category, status
}

func (n *Navigator) ReportFiltered(category model.Category, status model.Status) string {
	return n.ReportTools(n.List(category, status))
}

func (n *Navigator) ReportTools(tools []model.Tool) string { return report.RenderCatalog(tools) }
