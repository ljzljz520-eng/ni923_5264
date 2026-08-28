package model

import (
	"fmt"
	"strings"
)

type AuditEvent struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Actor    string   `json:"actor"`
	ToolIDs  []string `json:"tool_ids"`
	Revision int      `json:"revision"`
	Detail   string   `json:"detail"`
}

func (e AuditEvent) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("audit event id is required")
	}
	if e.Type == "" {
		return fmt.Errorf("audit event type is required")
	}
	if e.Actor == "" {
		return fmt.Errorf("audit actor is required")
	}
	if e.Revision < 0 {
		return fmt.Errorf("audit revision cannot be negative")
	}
	return nil
}

func (e AuditEvent) Clone() AuditEvent {
	e.ToolIDs = append([]string(nil), e.ToolIDs...)
	return e
}

func (e AuditEvent) Mentions(id string) bool {
	for _, current := range e.ToolIDs {
		if strings.TrimSpace(current) == id {
			return true
		}
	}
	return false
}
