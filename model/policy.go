package model

import (
	"fmt"
	"strings"
)

type Policy struct {
	AllowedCategories map[Category]bool
	AllowedStatuses   map[Status]bool
	RequireHTTPS      bool
	MaxTags           int
}

func DefaultPolicy() Policy {
	allowedCategories := make(map[Category]bool)
	for _, category := range Categories() {
		allowedCategories[category] = true
	}
	return Policy{AllowedCategories: allowedCategories, AllowedStatuses: map[Status]bool{StatusActive: true, StatusBeta: true, StatusArchived: true}, RequireHTTPS: true, MaxTags: 8}
}

func (p Policy) Validate(tool Tool) error {
	if strings.TrimSpace(tool.ID) == "" || strings.TrimSpace(tool.Name) == "" {
		return fmt.Errorf("identity fields are required")
	}
	if !p.AllowedCategories[tool.Category] {
		return fmt.Errorf("category is not enabled")
	}
	if !p.AllowedStatuses[tool.Status] {
		return fmt.Errorf("status is not enabled")
	}
	if p.RequireHTTPS && !strings.HasPrefix(tool.URL, "https://") {
		return fmt.Errorf("secure URL is required")
	}
	if p.MaxTags >= 0 && len(tool.Tags) > p.MaxTags {
		return fmt.Errorf("too many tags")
	}
	return nil
}

func NormalizeTags(tags []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(tags))
	for _, raw := range tags {
		tag := strings.ToLower(strings.TrimSpace(raw))
		if tag == "" || seen[tag] {
			continue
		}
		seen[tag] = true
		result = append(result, tag)
	}
	return result
}

func (p Policy) Normalize(tool Tool) Tool {
	tool.ID = strings.TrimSpace(strings.ToLower(tool.ID))
	tool.Name = strings.TrimSpace(tool.Name)
	tool.URL = strings.TrimSpace(tool.URL)
	tool.Tags = NormalizeTags(tool.Tags)
	return tool
}

func (p Policy) Explain(tool Tool) []string {
	issues := make([]string, 0)
	if strings.TrimSpace(tool.ID) == "" {
		issues = append(issues, "missing id")
	}
	if strings.TrimSpace(tool.Name) == "" {
		issues = append(issues, "missing name")
	}
	if !p.AllowedCategories[tool.Category] {
		issues = append(issues, "category disabled")
	}
	if !p.AllowedStatuses[tool.Status] {
		issues = append(issues, "status disabled")
	}
	if p.RequireHTTPS && !strings.HasPrefix(tool.URL, "https://") {
		issues = append(issues, "URL is not secure")
	}
	return issues
}
