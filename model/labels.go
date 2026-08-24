package model

import "strings"

func CategoryLabel(category Category) string {
	parts := strings.Split(string(category), "-")
	for index, part := range parts {
		if len(part) == 0 {
			continue
		}
		parts[index] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}

func StatusLabel(status Status) string {
	if status == StatusArchived {
		return "Archived"
	}
	if status == StatusBeta {
		return "Beta"
	}
	return "Active"
}

func ToolKey(tool Tool) string { return strings.ToLower(strings.TrimSpace(tool.ID)) }

func SameTool(left, right Tool) bool {
	return left.ID == right.ID && left.Name == right.Name && left.Category == right.Category && left.URL == right.URL && left.Status == right.Status && left.Rank == right.Rank && strings.Join(left.Tags, ",") == strings.Join(right.Tags, ",")
}
