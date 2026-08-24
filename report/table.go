package report

import (
	"fmt"
	"strings"

	"example.com/toolnav/model"
)

func RenderCompact(tools []model.Tool) string {
	rows := make([]string, 0, len(tools)+1)
	rows = append(rows, "id|name|category|status|rank")
	for _, tool := range tools {
		rows = append(rows, fmt.Sprintf("%s|%s|%s|%s|%d", tool.ID, tool.Name, tool.Category, tool.Status, tool.Rank))
	}
	return strings.Join(rows, "\n")
}

func RenderTool(tool model.Tool) string {
	return fmt.Sprintf("%s (%s) %s", tool.Name, tool.Category, tool.Status)
}
