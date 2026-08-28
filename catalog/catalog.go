package catalog

import (
	"fmt"
	"strings"

	"example.com/toolnav/audit"
	"example.com/toolnav/importer"
	"example.com/toolnav/model"
	"example.com/toolnav/store"
)

type Service struct {
	Store *store.Store
	Audit *audit.Recorder
	Tools map[string]model.Tool
	Order model.ToolOrder
}

func NewService(s *store.Store) (*Service, error) {
	state, err := s.LoadCatalog()
	if err != nil {
		return nil, err
	}
	if state.Tools == nil {
		state.Tools = make(map[string]model.Tool)
	}
	if len(state.Order.ToolIDs) == 0 && len(state.Tools) > 0 {
		state.Order = model.ToolOrder{ToolIDs: sortedIDs(state.Tools), Revision: 1}
	}
	return &Service{Store: s, Audit: audit.NewRecorder(s, "admin"), Tools: state.Tools, Order: state.Order}, nil
}

func sortedIDs(tools map[string]model.Tool) []string {
	ids := make([]string, 0, len(tools))
	for id := range tools {
		ids = append(ids, id)
	}
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[j] < ids[i] {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
	}
	return ids
}

func (c *Service) AddTool(tool model.Tool) error {
	if err := tool.ValidateTool(); err != nil {
		return err
	}
	if _, exists := c.Tools[tool.ID]; exists {
		return fmt.Errorf("tool %q already exists", tool.ID)
	}
	tool.Rank = len(c.Order.ToolIDs)
	c.Tools[tool.ID] = tool
	c.Order.ToolIDs = append(c.Order.ToolIDs, tool.ID)
	c.Order.Revision++
	if err := c.Store.SaveCatalog(c.Tools, c.Order, nil); err != nil {
		return err
	}
	_, err := c.Audit.Record(audit.EventStatusChanged, []string{tool.ID}, c.Order.Revision, "tool added")
	return err
}

func (c *Service) ImportRows(batchID string, source string, rows []string) (model.ImportBatch, error) {
	parser := importer.New(source)
	batch, tools := parser.BuildBatch(batchID, rows)
	for _, tool := range tools {
		if _, exists := c.Tools[tool.ID]; exists {
			batch.Rejected++
			batch.Accepted--
			batch.Errors = append(batch.Errors, fmt.Sprintf("duplicate tool id %q", tool.ID))
			continue
		}
		tool.Rank = len(c.Order.ToolIDs)
		c.Tools[tool.ID] = tool
		c.Order.ToolIDs = append(c.Order.ToolIDs, tool.ID)
	}
	c.Order.Revision++
	if err := c.Store.SaveCatalog(c.Tools, c.Order, &batch); err != nil {
		return batch, err
	}
	_, err := c.Audit.Record(audit.EventImported, c.Order.ToolIDs, c.Order.Revision, batch.Source)
	return batch, err
}

func (c *Service) Query(category model.Category, status model.Status) []model.Tool {
	result := make([]model.Tool, 0, len(c.Order.ToolIDs))
	for _, id := range c.Order.ToolIDs {
		tool, ok := c.Tools[id]
		if !ok {
			continue
		}
		if category != "" && tool.Category != category {
			continue
		}
		if status != "" && tool.Status != status {
			continue
		}
		result = append(result, tool)
	}
	return result
}

func (c *Service) Move(id string, target int, beforeMutation func()) error {
	index := c.Order.Position(id)
	if index < 0 {
		return fmt.Errorf("tool %q is not in order", id)
	}
	if target < 0 {
		target = 0
	}
	if target >= len(c.Order.ToolIDs) {
		target = len(c.Order.ToolIDs) - 1
	}
	if index == target {
		return nil
	}
	working := c.Order.ToolIDs
	removed := working[index]
	copy(working[index:], working[index+1:])
	working = working[:len(working)-1]
	if beforeMutation != nil {
		beforeMutation()
	}
	working = append(working, "")
	copy(working[target+1:], working[target:])
	working[target] = removed
	c.Order.ToolIDs = working
	c.Order.Revision++
	for rank, toolID := range working {
		tool := c.Tools[toolID]
		tool.Rank = rank
		c.Tools[toolID] = tool
	}
	if err := c.Store.SaveCatalog(c.Tools, c.Order, nil); err != nil {
		return err
	}
	_, err := c.Audit.Record(audit.EventMoved, c.Order.ToolIDs, c.Order.Revision, fmt.Sprintf("moved %s to %d", id, target))
	return err
}

func (c *Service) Validate() error { return c.Order.ValidateOrder(c.Tools) }

func (c *Service) Count() int { return len(c.Tools) }

func (c *Service) Search(term string) []model.Tool {
	term = strings.ToLower(strings.TrimSpace(term))
	if term == "" {
		return c.Query("", "")
	}
	result := make([]model.Tool, 0)
	for _, tool := range c.Query("", "") {
		if strings.Contains(strings.ToLower(tool.Name), term) || strings.Contains(strings.ToLower(tool.ID), term) {
			result = append(result, tool)
		}
	}
	return result
}
