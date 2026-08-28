package catalog

import (
	"fmt"

	"example.com/toolnav/model"
)

func (c *Service) ValidateTools() []string {
	errors := make([]string, 0)
	for id, tool := range c.Tools {
		if id != tool.ID {
			errors = append(errors, fmt.Sprintf("map key %q differs from tool id %q", id, tool.ID))
		}
		if err := tool.ValidateTool(); err != nil {
			errors = append(errors, err.Error())
		}
	}
	if err := c.Order.ValidateOrder(c.Tools); err != nil {
		errors = append(errors, err.Error())
	}
	return errors
}

func (c *Service) RebuildRanks() error {
	if err := c.Order.ValidateOrder(c.Tools); err != nil {
		return err
	}
	for rank, id := range c.Order.ToolIDs {
		tool := c.Tools[id]
		tool.Rank = rank
		c.Tools[id] = tool
	}
	return c.Store.SaveCatalog(c.Tools, c.Order, nil)
}

func (c *Service) ToolAt(rank int) (model.Tool, bool) {
	if rank < 0 || rank >= len(c.Order.ToolIDs) {
		return model.Tool{}, false
	}
	tool, ok := c.Tools[c.Order.ToolIDs[rank]]
	return tool, ok
}
