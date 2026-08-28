package toolnav

import (
	"example.com/toolnav/backup"
	"example.com/toolnav/catalog"
	"example.com/toolnav/exporter"
	"example.com/toolnav/governance"
	"example.com/toolnav/metrics"
	"example.com/toolnav/model"
	"example.com/toolnav/report"
	"example.com/toolnav/store"
	"example.com/toolnav/workflow"
)

type Navigator struct {
	Store   *store.Store
	Catalog *catalog.Service
	Backup  *backup.Service
}

func Open(path string) (*Navigator, error) {
	s, err := store.Open(path)
	if err != nil {
		return nil, err
	}
	c, err := catalog.NewService(s)
	if err != nil {
		s.Close()
		return nil, err
	}
	return &Navigator{Store: s, Catalog: c, Backup: backup.NewService(s)}, nil
}

func (n *Navigator) Close() error { return n.Store.Close() }

func (n *Navigator) Add(tool model.Tool) error { return n.Catalog.AddTool(tool) }

func (n *Navigator) Import(batchID, source string, rows []string) (model.ImportBatch, error) {
	return n.Catalog.ImportRows(batchID, source, rows)
}

func (n *Navigator) List(category model.Category, status model.Status) []model.Tool {
	return n.Catalog.Query(category, status)
}

func (n *Navigator) Search(term string) []model.Tool { return n.Catalog.Search(term) }

func (n *Navigator) Move(id string, target int, hook func()) error {
	return n.Catalog.Move(id, target, hook)
}

func (n *Navigator) CreateBackup(label string) (model.BackupSnapshot, error) {
	snapshot, err := n.Backup.Create(n.Catalog.Order, n.Catalog.Tools, label)
	if err != nil {
		return model.BackupSnapshot{}, err
	}
	if n.Catalog.Audit != nil {
		if _, err := n.Catalog.Audit.Record("backed_up", snapshot.ToolIDs, snapshot.Revision, label); err != nil {
			return model.BackupSnapshot{}, err
		}
	}
	return snapshot, nil
}

func (n *Navigator) Validate() error { return n.Catalog.Validate() }

func (n *Navigator) CatalogReport() string { return report.RenderCatalog(n.List("", "")) }

func (n *Navigator) BackupReport() (string, error) {
	snapshot, err := n.Backup.Latest()
	if err != nil {
		return "", err
	}
	return report.RenderBackup(snapshot), nil
}

func (n *Navigator) RecoveryReport() (string, error) {
	state, err := n.Store.LoadCatalog()
	if err != nil {
		return "", err
	}
	return report.RenderRecovery(state.Tools, state.Order, len(state.Backups)), nil
}

func (n *Navigator) ImportAndReview(batchID, source string, rows []string, policy model.Policy) (workflow.ImportResult, error) {
	return workflow.ImportAndReview(n.Catalog, batchID, source, rows, policy)
}

func (n *Navigator) UpdateStatus(id string, status model.Status, policy model.Policy) (governance.Decision, error) {
	return n.Catalog.UpdateStatus(id, status, policy)
}

func (n *Navigator) RestoreLatest() (model.ToolOrder, error) {
	order, err := workflow.Restore(n.Catalog, n.Backup)
	if err != nil {
		return model.ToolOrder{}, err
	}
	if n.Catalog.Audit != nil {
		if _, err := n.Catalog.Audit.Record("restored", order.ToolIDs, order.Revision, "latest backup"); err != nil {
			return model.ToolOrder{}, err
		}
	}
	return order, nil
}

func (n *Navigator) Dashboard(policy model.Policy) (metrics.Dashboard, error) {
	state, err := n.Store.LoadCatalog()
	if err != nil {
		return metrics.Dashboard{}, err
	}
	return metrics.BuildDashboard(n.Catalog.OrderedTools(), state.Backups, state.Audits, policy), nil
}

func (n *Navigator) ExportBundle() (exporter.Bundle, error) {
	state, err := n.Store.LoadCatalog()
	if err != nil {
		return exporter.Bundle{}, err
	}
	return exporter.NewBundle(state.Tools, state.Order, state.Backups, state.Audits), nil
}
