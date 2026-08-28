package importer

import (
	"fmt"
	"strings"

	"example.com/toolnav/model"
)

type Parser struct {
	Source string
}

func New(source string) Parser { return Parser{Source: source} }

func (p Parser) ParseRows(rows []string) ([]model.Tool, []string) {
	tools := make([]model.Tool, 0, len(rows))
	errors := make([]string, 0)
	for index, row := range rows {
		tool, err := p.ParseRow(row)
		if err != nil {
			errors = append(errors, fmt.Sprintf("row %d: %v", index+1, err))
			continue
		}
		tools = append(tools, tool)
	}
	return tools, errors
}

func (p Parser) ParseRow(row string) (model.Tool, error) {
	parts := strings.Split(row, "|")
	if len(parts) != 7 {
		return model.Tool{}, fmt.Errorf("expected 7 fields, got %d", len(parts))
	}
	tags := splitTags(parts[6])
	tool := model.Tool{ID: strings.TrimSpace(parts[0]), Name: strings.TrimSpace(parts[1]), Category: model.Category(strings.TrimSpace(parts[2])), URL: strings.TrimSpace(parts[3]), Status: model.Status(strings.TrimSpace(parts[4])), Rank: parseRank(parts[5]), Tags: tags}
	if err := tool.ValidateTool(); err != nil {
		return model.Tool{}, err
	}
	return tool, nil
}

func parseRank(value string) int {
	var rank int
	for _, r := range strings.TrimSpace(value) {
		if r < '0' || r > '9' {
			return -1
		}
		rank = rank*10 + int(r-'0')
	}
	return rank
}

func splitTags(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if tag := strings.TrimSpace(part); tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

func (p Parser) BuildBatch(id string, rows []string) (model.ImportBatch, []model.Tool) {
	tools, errors := p.ParseRows(rows)
	return model.ImportBatch{ID: id, Source: p.Source, Rows: len(rows), Accepted: len(tools), Rejected: len(errors), Errors: errors}, tools
}
