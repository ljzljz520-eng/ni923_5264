package catalog

import (
	"sort"

	"example.com/toolnav/model"
)

type Summary struct {
	Total      int
	ByCategory map[model.Category]int
	ByStatus   map[model.Status]int
	Tagged     int
}

func (c *Service) Summarize() Summary {
	summary := Summary{ByCategory: make(map[model.Category]int), ByStatus: make(map[model.Status]int)}
	for _, tool := range c.Tools {
		summary.Total++
		summary.ByCategory[tool.Category]++
		summary.ByStatus[tool.Status]++
		if len(tool.Tags) > 0 {
			summary.Tagged++
		}
	}
	return summary
}

func (c *Service) CategoriesWithCounts() []string {
	summary := c.Summarize()
	values := make([]string, 0, len(summary.ByCategory))
	for category, count := range summary.ByCategory {
		values = append(values, string(category)+"="+itoa(count))
	}
	sort.Strings(values)
	return values
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	digits := make([]byte, 0, 12)
	for value > 0 {
		digits = append(digits, byte('0'+value%10))
		value /= 10
	}
	for left, right := 0, len(digits)-1; left < right; left, right = left+1, right-1 {
		digits[left], digits[right] = digits[right], digits[left]
	}
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}

func (c *Service) TopTags(limit int) []string {
	counts := make(map[string]int)
	for _, tool := range c.Tools {
		for _, tag := range tool.Tags {
			counts[tag]++
		}
	}
	type pair struct {
		tag   string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for tag, count := range counts {
		pairs = append(pairs, pair{tag: tag, count: count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].tag < pairs[j].tag
		}
		return pairs[i].count > pairs[j].count
	})
	if limit < 0 {
		limit = 0
	}
	if limit > len(pairs) {
		limit = len(pairs)
	}
	result := make([]string, limit)
	for index := range result {
		result[index] = pairs[index].tag
	}
	return result
}
