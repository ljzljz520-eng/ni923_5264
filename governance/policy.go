package governance

import (
	"fmt"
	"sort"
	"strings"

	"example.com/toolnav/model"
)

type Decision struct {
	Allowed bool
	From    model.Status
	To      model.Status
	Reason  string
}

type Review struct {
	ToolID   string
	Category model.Category
	Status   model.Status
	Issues   []string
	Priority int
}

func Evaluate(tool model.Tool, policy model.Policy, requested model.Status) Decision {
	decision := Decision{From: tool.Status, To: requested}
	if !policy.AllowedStatuses[requested] {
		decision.Reason = "requested status is disabled"
		return decision
	}
	if !model.ValidStatus(requested) {
		decision.Reason = "requested status is unknown"
		return decision
	}
	if !CanTransition(tool.Status, requested) {
		decision.Reason = fmt.Sprintf("transition from %s to %s is not allowed", tool.Status, requested)
		return decision
	}
	if issues := policy.Explain(tool); len(issues) > 0 && requested == model.StatusActive {
		decision.Reason = strings.Join(issues, "; ")
		return decision
	}
	decision.Allowed = true
	decision.Reason = "status transition accepted"
	return decision
}

func CanTransition(from, to model.Status) bool {
	if from == to {
		return true
	}
	switch from {
	case model.StatusBeta:
		return to == model.StatusActive || to == model.StatusArchived
	case model.StatusActive:
		return to == model.StatusArchived || to == model.StatusBeta
	case model.StatusArchived:
		return to == model.StatusBeta
	default:
		return false
	}
}

func BuildReview(tools []model.Tool, policy model.Policy) []Review {
	reviews := make([]Review, 0)
	for _, tool := range tools {
		issues := policy.Explain(tool)
		if len(issues) == 0 {
			continue
		}
		priority := len(issues)
		if tool.Status == model.StatusActive {
			priority++
		}
		reviews = append(reviews, Review{ToolID: tool.ID, Category: tool.Category, Status: tool.Status, Issues: append([]string(nil), issues...), Priority: priority})
	}
	sort.Slice(reviews, func(i, j int) bool {
		if reviews[i].Priority == reviews[j].Priority {
			return reviews[i].ToolID < reviews[j].ToolID
		}
		return reviews[i].Priority > reviews[j].Priority
	})
	return reviews
}

func ReviewSummary(reviews []Review) map[string]int {
	result := map[string]int{"total": len(reviews), "high": 0, "medium": 0, "low": 0}
	for _, review := range reviews {
		switch {
		case review.Priority >= 3:
			result["high"]++
		case review.Priority == 2:
			result["medium"]++
		default:
			result["low"]++
		}
	}
	return result
}

func Coverage(tools []model.Tool) map[model.Category]int {
	coverage := make(map[model.Category]int)
	for _, tool := range tools {
		coverage[tool.Category]++
	}
	return coverage
}

func MissingCategories(tools []model.Tool) []model.Category {
	coverage := Coverage(tools)
	missing := make([]model.Category, 0)
	for _, category := range model.Categories() {
		if coverage[category] == 0 {
			missing = append(missing, category)
		}
	}
	return missing
}

func ComplianceScore(tools []model.Tool, policy model.Policy) int {
	if len(tools) == 0 {
		return 100
	}
	valid := 0
	for _, tool := range tools {
		if len(policy.Explain(tool)) == 0 {
			valid++
		}
	}
	return valid * 100 / len(tools)
}
