package metrics

import (
	"fmt"
	"sort"
	"strings"

	"example.com/toolnav/audit"
	"example.com/toolnav/governance"
	"example.com/toolnav/model"
)

type Dashboard struct {
	ToolCount        int
	ActiveCount      int
	BetaCount        int
	ArchivedCount    int
	CategoryCoverage int
	ComplianceScore  int
	BackupCount      int
	AuditCount       int
	TopTags          []string
}

func BuildDashboard(tools []model.Tool, backups []model.BackupSnapshot, events []model.AuditEvent, policy model.Policy) Dashboard {
	dashboard := Dashboard{ToolCount: len(tools), BackupCount: len(backups), AuditCount: len(events), TopTags: TopTags(tools, 5)}
	for _, tool := range tools {
		switch tool.Status {
		case model.StatusActive:
			dashboard.ActiveCount++
		case model.StatusBeta:
			dashboard.BetaCount++
		case model.StatusArchived:
			dashboard.ArchivedCount++
		}
	}
	dashboard.CategoryCoverage = len(model.Categories()) - len(governance.MissingCategories(tools))
	dashboard.ComplianceScore = governance.ComplianceScore(tools, policy)
	return dashboard
}

func (d Dashboard) Healthy() bool {
	return d.ToolCount > 0 && d.CategoryCoverage >= 3 && d.ComplianceScore >= 80
}

func (d Dashboard) StatusCounts() map[model.Status]int {
	return map[model.Status]int{model.StatusActive: d.ActiveCount, model.StatusBeta: d.BetaCount, model.StatusArchived: d.ArchivedCount}
}

func (d Dashboard) Format() string {
	return fmt.Sprintf("tools=%d active=%d beta=%d archived=%d categories=%d compliance=%d backups=%d audits=%d tags=%s", d.ToolCount, d.ActiveCount, d.BetaCount, d.ArchivedCount, d.CategoryCoverage, d.ComplianceScore, d.BackupCount, d.AuditCount, strings.Join(d.TopTags, ","))
}

func TopTags(tools []model.Tool, limit int) []string {
	counts := make(map[string]int)
	for _, tool := range tools {
		for _, tag := range tool.Tags {
			counts[tag]++
		}
	}
	type pair struct {
		value string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for value, count := range counts {
		pairs = append(pairs, pair{value: value, count: count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].value < pairs[j].value
		}
		return pairs[i].count > pairs[j].count
	})
	if limit < 0 {
		limit = 0
	}
	if limit > len(pairs) {
		limit = len(pairs)
	}
	result := make([]string, limit)
	for index := range result {
		result[index] = pairs[index].value
	}
	return result
}

func EventHealth(events []model.AuditEvent) string {
	counts := audit.CountByType(events)
	if counts[audit.EventImported] == 0 {
		return "waiting-for-import"
	}
	if counts[audit.EventBackedUp] == 0 {
		return "backup-required"
	}
	if counts[audit.EventMoved] > counts[audit.EventBackedUp] {
		return "backup-stale"
	}
	return "operational"
}

func CoverageLabel(d Dashboard) string {
	if d.CategoryCoverage == len(model.Categories()) {
		return "complete"
	}
	return fmt.Sprintf("%d/%d categories", d.CategoryCoverage, len(model.Categories()))
}
