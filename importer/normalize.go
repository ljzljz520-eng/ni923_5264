package importer

import (
	"strings"
	"unicode"

	"example.com/toolnav/model"
)

func NormalizeRow(row string) string {
	fields := strings.Split(row, "|")
	for index, field := range fields {
		fields[index] = strings.TrimSpace(field)
	}
	return strings.Join(fields, "|")
}

func Slug(value string) string {
	var builder strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && builder.Len() > 0 {
			builder.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(builder.String(), "-")
}

func BuildTool(id, name string, category model.Category, rawURL string, status model.Status, tags []string) (model.Tool, error) {
	tool := model.Tool{ID: Slug(id), Name: strings.TrimSpace(name), Category: category, URL: strings.TrimSpace(rawURL), Status: status, Tags: model.NormalizeTags(tags)}
	return tool, tool.ValidateTool()
}

func ValidateAll(tools []model.Tool) []string {
	errors := make([]string, 0)
	seen := make(map[string]bool)
	for index, tool := range tools {
		if err := tool.ValidateTool(); err != nil {
			errors = append(errors, err.Error())
		}
		if seen[tool.ID] {
			errors = append(errors, "duplicate id at row "+string(rune(index+'0')))
		}
		seen[tool.ID] = true
	}
	return errors
}
