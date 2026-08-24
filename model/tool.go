package model

import (
	"fmt"
	"net/url"
	"strings"
)

type Category string

const (
	CategoryDeployment    Category = "deployment"
	CategoryMonitoring    Category = "monitoring"
	CategoryPayments      Category = "payments"
	CategoryDocumentation Category = "documentation"
	CategoryCompliance    Category = "compliance"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusBeta     Status = "beta"
	StatusArchived Status = "archived"
)

type Tool struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Category Category `json:"category"`
	URL      string   `json:"url"`
	Status   Status   `json:"status"`
	Rank     int      `json:"rank"`
	Tags     []string `json:"tags"`
}

func (t Tool) ValidateTool() error {
	if strings.TrimSpace(t.ID) == "" {
		return fmt.Errorf("tool id is required")
	}
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("tool name is required")
	}
	if !ValidCategory(t.Category) {
		return fmt.Errorf("unsupported category %q", t.Category)
	}
	u, err := url.Parse(t.URL)
	if err != nil || u.Scheme != "https" || u.Host == "" {
		return fmt.Errorf("tool url must be an https URL")
	}
	if !ValidStatus(t.Status) {
		return fmt.Errorf("unsupported status %q", t.Status)
	}
	if t.Rank < 0 {
		return fmt.Errorf("rank cannot be negative")
	}
	return nil
}

func ValidCategory(c Category) bool {
	switch c {
	case CategoryDeployment, CategoryMonitoring, CategoryPayments, CategoryDocumentation, CategoryCompliance:
		return true
	default:
		return false
	}
}

func ValidStatus(s Status) bool {
	switch s {
	case StatusActive, StatusBeta, StatusArchived:
		return true
	default:
		return false
	}
}

func Categories() []Category {
	return []Category{CategoryDeployment, CategoryMonitoring, CategoryPayments, CategoryDocumentation, CategoryCompliance}
}
