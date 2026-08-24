package catalog

import (
	"sort"
	"strings"

	"example.com/toolnav/model"
)

type FilterOptions struct {
	Category model.Category
	Status   model.Status
	Tag      string
	Term     string
	MinRank  int
	MaxRank  int
}

func Match(tool model.Tool, options FilterOptions) bool {
	if options.Category != "" && tool.Category != options.Category {
		return false
	}
	if options.Status != "" && tool.Status != options.Status {
		return false
	}
	if options.Tag != "" && !hasTag(tool.Tags, options.Tag) {
		return false
	}
	if options.Term != "" {
		term := strings.ToLower(strings.TrimSpace(options.Term))
		if !strings.Contains(strings.ToLower(tool.ID), term) && !strings.Contains(strings.ToLower(tool.Name), term) {
			return false
		}
	}
	if options.MinRank >= 0 && tool.Rank < options.MinRank {
		return false
	}
	if options.MaxRank >= 0 && tool.Rank > options.MaxRank {
		return false
	}
	return true
}

func hasTag(tags []string, expected string) bool {
	expected = strings.ToLower(strings.TrimSpace(expected))
	for _, tag := range tags {
		if strings.ToLower(tag) == expected {
			return true
		}
	}
	return false
}

func (c *Service) Filter(options FilterOptions) []model.Tool {
	result := make([]model.Tool, 0)
	for _, id := range c.Order.ToolIDs {
		tool, ok := c.Tools[id]
		if ok && Match(tool, options) {
			result = append(result, tool)
		}
	}
	return result
}

func SortByName(tools []model.Tool) []model.Tool {
	result := append([]model.Tool(nil), tools...)
	sort.SliceStable(result, func(i, j int) bool {
		left := strings.ToLower(result[i].Name)
		right := strings.ToLower(result[j].Name)
		if left == right {
			return result[i].ID < result[j].ID
		}
		return left < right
	})
	return result
}

func SortByCategory(tools []model.Tool) []model.Tool {
	result := append([]model.Tool(nil), tools...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Category == result[j].Category {
			return result[i].Rank < result[j].Rank
		}
		return result[i].Category < result[j].Category
	})
	return result
}

func CategoriesPresent(tools []model.Tool) []model.Category {
	seen := make(map[model.Category]bool)
	for _, tool := range tools {
		seen[tool.Category] = true
	}
	result := make([]model.Category, 0, len(seen))
	for category := range seen {
		result = append(result, category)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func RankGaps(order model.ToolOrder) []int {
	gaps := make([]int, 0)
	for index, id := range order.ToolIDs {
		if id == "" {
			gaps = append(gaps, index)
		}
	}
	return gaps
}

func DuplicateRanks(tools []model.Tool) []int {
	counts := make(map[int]int)
	for _, tool := range tools {
		counts[tool.Rank]++
	}
	duplicates := make([]int, 0)
	for rank, count := range counts {
		if count > 1 {
			duplicates = append(duplicates, rank)
		}
	}
	sort.Ints(duplicates)
	return duplicates
}

func ValidateView(tools []model.Tool, order model.ToolOrder) []string {
	errors := make([]string, 0)
	if len(tools) != len(order.ToolIDs) {
		errors = append(errors, "view count does not match order")
	}
	for _, index := range RankGaps(order) {
		errors = append(errors, "empty order position "+viewItoa(index))
	}
	for _, rank := range DuplicateRanks(tools) {
		errors = append(errors, "duplicate rank "+viewItoa(rank))
	}
	return errors
}

func viewItoa(value int) string {
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
