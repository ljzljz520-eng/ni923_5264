package governance

import (
	"testing"

	"example.com/toolnav/model"
)

func TestStatusPolicyAndReview(t *testing.T) {
	tool := model.Tool{ID: "deployctl", Name: "DeployCtl", Category: model.CategoryDeployment, URL: "https://deploy.example", Status: model.StatusBeta}
	decision := Evaluate(tool, model.DefaultPolicy(), model.StatusActive)
	if !decision.Allowed {
		t.Fatalf("expected transition to be allowed: %+v", decision)
	}
	if CanTransition(model.StatusArchived, model.StatusActive) {
		t.Fatal("archived tools should require an intermediate review")
	}
	reviews := BuildReview([]model.Tool{{ID: "bad", Category: model.CategoryDeployment, Status: model.StatusActive}}, model.DefaultPolicy())
	if len(reviews) != 1 || reviews[0].Priority < 2 {
		t.Fatalf("unexpected reviews: %+v", reviews)
	}
}
