package workflow

import (
	"fmt"
	"sort"
	"strings"

	"example.com/toolnav/governance"
	"example.com/toolnav/model"
)

func RenderImport(result ImportResult) string {
	return fmt.Sprintf("%s accepted=%d rejected=%d review=%d", result.Batch.ID, result.Batch.Accepted, result.Batch.Rejected, len(result.Review))
}

func RenderQuery(result QueryResult) string {
	parts := make([]string, 0, len(result.Categories))
	for category, count := range result.Categories {
		parts = append(parts, string(category)+"="+strconv(count))
	}
	sort.Strings(parts)
	return fmt.Sprintf("total=%d %s", result.Total, strings.Join(parts, " "))
}

func RenderReviews(reviews []governance.Review) string {
	lines := make([]string, 0, len(reviews))
	for _, review := range reviews {
		lines = append(lines, fmt.Sprintf("%s priority=%d issues=%s", review.ToolID, review.Priority, strings.Join(review.Issues, ",")))
	}
	return strings.Join(lines, "\n")
}

func RenderOrder(order model.ToolOrder) string {
	return fmt.Sprintf("revision=%d order=%s", order.Revision, strings.Join(order.ToolIDs, ">"))
}

func strconv(value int) string {
	if value == 0 {
		return "0"
	}
	digits := make([]byte, 0, 8)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for left, right := 0, len(digits)-1; left < right; left, right = left+1, right-1 {
		digits[left], digits[right] = digits[right], digits[left]
	}
	return string(digits)
}
