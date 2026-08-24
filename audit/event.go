package audit

import (
	"fmt"
	"sort"
	"strings"

	"example.com/toolnav/model"
)

type EventType string

const (
	EventImported      EventType = "imported"
	EventQueried       EventType = "queried"
	EventMoved         EventType = "moved"
	EventBackedUp      EventType = "backed_up"
	EventRestored      EventType = "restored"
	EventStatusChanged EventType = "status_changed"
)

type AuditEvent = model.AuditEvent

func NewEvent(id string, eventType EventType, actor string, ids []string, revision int, detail string) AuditEvent {
	return AuditEvent{ID: strings.TrimSpace(id), Type: string(eventType), Actor: strings.TrimSpace(actor), ToolIDs: append([]string(nil), ids...), Revision: revision, Detail: strings.TrimSpace(detail)}
}

func ValidType(eventType EventType) bool {
	switch eventType {
	case EventImported, EventQueried, EventMoved, EventBackedUp, EventRestored, EventStatusChanged:
		return true
	default:
		return false
	}
}

func Filter(events []AuditEvent, eventType EventType, actor string) []AuditEvent {
	filtered := make([]AuditEvent, 0, len(events))
	for _, event := range events {
		if eventType != "" && EventType(event.Type) != eventType {
			continue
		}
		if actor != "" && event.Actor != actor {
			continue
		}
		filtered = append(filtered, event.Clone())
	}
	sort.SliceStable(filtered, func(i, j int) bool { return filtered[i].Revision < filtered[j].Revision })
	return filtered
}

func CountByType(events []AuditEvent) map[EventType]int {
	counts := make(map[EventType]int)
	for _, event := range events {
		counts[EventType(event.Type)]++
	}
	return counts
}

func ToolIDs(events []AuditEvent) []string {
	seen := make(map[string]bool)
	ids := make([]string, 0)
	for _, event := range events {
		for _, id := range event.ToolIDs {
			if !seen[id] {
				seen[id] = true
				ids = append(ids, id)
			}
		}
	}
	sort.Strings(ids)
	return ids
}

func ValidateForTools(events []AuditEvent, tools map[string]model.Tool) []string {
	errors := make([]string, 0)
	for _, event := range events {
		if err := event.Validate(); err != nil {
			errors = append(errors, err.Error())
		}
		if !ValidType(EventType(event.Type)) {
			errors = append(errors, fmt.Sprintf("unsupported audit event type %q", event.Type))
		}
		for _, id := range event.ToolIDs {
			if _, ok := tools[id]; !ok {
				errors = append(errors, fmt.Sprintf("audit event %q references unknown tool %q", event.ID, id))
			}
		}
	}
	return errors
}
