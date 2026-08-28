package model

import "fmt"

type ToolOrder struct {
	ToolIDs  []string `json:"tool_ids"`
	Revision int      `json:"revision"`
}

func (o ToolOrder) ValidateOrder(expected map[string]Tool) error {
	if len(o.ToolIDs) != len(expected) {
		return fmt.Errorf("order has %d ids, expected %d", len(o.ToolIDs), len(expected))
	}
	seen := make(map[string]bool, len(o.ToolIDs))
	for _, id := range o.ToolIDs {
		if seen[id] {
			return fmt.Errorf("duplicate tool id %q", id)
		}
		if _, ok := expected[id]; !ok {
			return fmt.Errorf("unknown tool id %q", id)
		}
		seen[id] = true
	}
	return nil
}

func (o ToolOrder) Clone() ToolOrder {
	ids := append([]string(nil), o.ToolIDs...)
	return ToolOrder{ToolIDs: ids, Revision: o.Revision}
}

func (o *ToolOrder) Position(id string) int {
	for i, current := range o.ToolIDs {
		if current == id {
			return i
		}
	}
	return -1
}
