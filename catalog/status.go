package catalog

import (
	"fmt"
	"sort"

	"example.com/toolnav/governance"
	"example.com/toolnav/model"
)

func (c *Service) UpdateStatus(id string, requested model.Status, policy model.Policy) (governance.Decision, error) {
	tool, ok := c.Tools[id]
	if !ok {
		return governance.Decision{}, fmt.Errorf("tool %q not found", id)
	}
	decision := governance.Evaluate(tool, policy, requested)
	if !decision.Allowed {
		return decision, fmt.Errorf("status change rejected: %s", decision.Reason)
	}
	tool.Status = requested
	c.Tools[id] = tool
	c.Order.Revision++
	if err := c.Store.SaveCatalog(c.Tools, c.Order, nil); err != nil {
		return governance.Decision{}, err
	}
	if c.Audit != nil {
		if _, err := c.Audit.Record("status_changed", []string{id}, c.Order.Revision, decision.Reason); err != nil {
			return governance.Decision{}, err
		}
	}
	return decision, nil
}

func (c *Service) ReplaceTool(tool model.Tool, policy model.Policy) error {
	if err := policy.Validate(tool); err != nil {
		return err
	}
	if _, ok := c.Tools[tool.ID]; !ok {
		return fmt.Errorf("tool %q not found", tool.ID)
	}
	tool.Rank = c.Tools[tool.ID].Rank
	c.Tools[tool.ID] = tool
	c.Order.Revision++
	return c.Store.SaveCatalog(c.Tools, c.Order, nil)
}

func (c *Service) OrderedTools() []model.Tool {
	result := make([]model.Tool, 0, len(c.Order.ToolIDs))
	for rank, id := range c.Order.ToolIDs {
		tool, ok := c.Tools[id]
		if !ok {
			continue
		}
		tool.Rank = rank
		result = append(result, tool)
	}
	return result
}

func (c *Service) SortedTools() []model.Tool {
	result := c.OrderedTools()
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Rank == result[j].Rank {
			return result[i].ID < result[j].ID
		}
		return result[i].Rank < result[j].Rank
	})
	return result
}

func (c *Service) ReorderByIDs(ids []string) error {
	order := model.ToolOrder{ToolIDs: append([]string(nil), ids...), Revision: c.Order.Revision + 1}
	if err := order.ValidateOrder(c.Tools); err != nil {
		return err
	}
	c.Order = order
	for rank, id := range ids {
		tool := c.Tools[id]
		tool.Rank = rank
		c.Tools[id] = tool
	}
	return c.Store.SaveCatalog(c.Tools, c.Order, nil)
}

func (c *Service) TagsFor(id string) []string {
	tool, ok := c.Tools[id]
	if !ok {
		return nil
	}
	return append([]string(nil), tool.Tags...)
}
