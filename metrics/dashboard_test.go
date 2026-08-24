package metrics

import (
	"testing"

	"example.com/toolnav/model"
)

func TestDashboard(t *testing.T) {
	tools := []model.Tool{
		{ID: "a", Name: "A", URL: "https://a.example", Category: model.CategoryDeployment, Status: model.StatusActive, Tags: []string{"ship"}},
		{ID: "b", Name: "B", URL: "https://b.example", Category: model.CategoryMonitoring, Status: model.StatusBeta, Tags: []string{"watch", "ship"}},
	}
	dashboard := BuildDashboard(tools, nil, nil, model.DefaultPolicy())
	if dashboard.ToolCount != 2 || dashboard.ActiveCount != 1 || dashboard.BetaCount != 1 || dashboard.ComplianceScore != 100 {
		t.Fatalf("unexpected dashboard: %+v", dashboard)
	}
	if EventHealth(nil) != "waiting-for-import" || len(TopTags(tools, 1)) != 1 {
		t.Fatal("dashboard signals are incorrect")
	}
}
