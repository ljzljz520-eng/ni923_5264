package audit

import (
	"fmt"
	"strings"

	"example.com/toolnav/model"
	"example.com/toolnav/store"
)

type Recorder struct {
	Store *store.Store
	Actor string
}

func NewRecorder(s *store.Store, actor string) *Recorder {
	if strings.TrimSpace(actor) == "" {
		actor = "admin"
	}
	return &Recorder{Store: s, Actor: strings.TrimSpace(actor)}
}

func (r *Recorder) Record(eventType EventType, ids []string, revision int, detail string) (AuditEvent, error) {
	if r == nil || r.Store == nil {
		return AuditEvent{}, fmt.Errorf("audit recorder is not configured")
	}
	eventID := fmt.Sprintf("event-%d-%s", revision, eventType)
	event := NewEvent(eventID, eventType, r.Actor, ids, revision, detail)
	if err := event.Validate(); err != nil {
		return AuditEvent{}, err
	}
	if err := r.Store.SaveAudit(event); err != nil {
		return AuditEvent{}, err
	}
	return event, nil
}

func (r *Recorder) Events() ([]AuditEvent, error) {
	if r == nil || r.Store == nil {
		return nil, fmt.Errorf("audit recorder is not configured")
	}
	state, err := r.Store.LoadCatalog()
	if err != nil {
		return nil, err
	}
	return Filter(state.Audits, "", r.Actor), nil
}

func (r *Recorder) EventsFor(eventType EventType) ([]AuditEvent, error) {
	events, err := r.Events()
	if err != nil {
		return nil, err
	}
	return Filter(events, eventType, ""), nil
}

func (r *Recorder) Summary() (map[EventType]int, error) {
	events, err := r.Events()
	if err != nil {
		return nil, err
	}
	return CountByType(events), nil
}

func Describe(event AuditEvent) string {
	return fmt.Sprintf("%s by %s revision=%d ids=%d detail=%s", event.Type, event.Actor, event.Revision, len(event.ToolIDs), event.Detail)
}

func HasMutation(events []AuditEvent) bool {
	for _, event := range events {
		if EventType(event.Type) == EventMoved || EventType(event.Type) == EventImported || EventType(event.Type) == EventStatusChanged {
			return true
		}
	}
	return false
}

func EventForTool(events []AuditEvent, id string) []AuditEvent {
	result := make([]AuditEvent, 0)
	for _, event := range events {
		if event.Mentions(id) {
			result = append(result, event.Clone())
		}
	}
	return result
}

func ValidateEventType(eventType model.Status) bool {
	return eventType == model.StatusActive || eventType == model.StatusBeta || eventType == model.StatusArchived
}
