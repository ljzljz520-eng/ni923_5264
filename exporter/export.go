package exporter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"example.com/toolnav/audit"
	"example.com/toolnav/backup"
	"example.com/toolnav/model"
)

type Bundle struct {
	Tools   []model.Tool           `json:"tools"`
	Order   model.ToolOrder        `json:"order"`
	Backups []model.BackupSnapshot `json:"backups"`
	Audits  []model.AuditEvent     `json:"audits"`
}

func NewBundle(tools map[string]model.Tool, order model.ToolOrder, backups []model.BackupSnapshot, audits []model.AuditEvent) Bundle {
	ordered := make([]model.Tool, 0, len(order.ToolIDs))
	for _, id := range order.ToolIDs {
		if tool, ok := tools[id]; ok {
			ordered = append(ordered, tool)
		}
	}
	return Bundle{Tools: ordered, Order: order.Clone(), Backups: cloneBackups(backups), Audits: cloneAudits(audits)}
}

func cloneBackups(values []model.BackupSnapshot) []model.BackupSnapshot {
	result := make([]model.BackupSnapshot, 0, len(values))
	for _, value := range values {
		result = append(result, value.Clone())
	}
	return result
}

func cloneAudits(values []model.AuditEvent) []model.AuditEvent {
	result := make([]model.AuditEvent, 0, len(values))
	for _, value := range values {
		value.ToolIDs = append([]string(nil), value.ToolIDs...)
		result = append(result, value)
	}
	return result
}

func JSON(bundle Bundle) ([]byte, error) {
	return json.MarshalIndent(bundle, "", "  ")
}

func ParseJSON(data []byte) (Bundle, error) {
	var bundle Bundle
	if len(data) == 0 {
		return Bundle{}, fmt.Errorf("export data is empty")
	}
	if err := json.Unmarshal(data, &bundle); err != nil {
		return Bundle{}, err
	}
	return bundle, Validate(bundle)
}

func CSV(tools []model.Tool, order model.ToolOrder) string {
	lines := []string{"id,name,category,status,rank,tags"}
	for rank, id := range order.ToolIDs {
		tool, ok := findTool(tools, id)
		if !ok {
			continue
		}
		lines = append(lines, strings.Join([]string{tool.ID, quote(tool.Name), string(tool.Category), string(tool.Status), strconv(rank), quote(strings.Join(tool.Tags, ","))}, ","))
	}
	return strings.Join(lines, "\n")
}

func findTool(tools []model.Tool, id string) (model.Tool, bool) {
	for _, tool := range tools {
		if tool.ID == id {
			return tool, true
		}
	}
	return model.Tool{}, false
}

func quote(value string) string {
	if strings.ContainsAny(value, ",\"\n") {
		return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
	}
	return value
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

func Validate(bundle Bundle) error {
	tools := make(map[string]model.Tool, len(bundle.Tools))
	for _, tool := range bundle.Tools {
		if _, exists := tools[tool.ID]; exists {
			return fmt.Errorf("duplicate exported tool %q", tool.ID)
		}
		if err := tool.ValidateTool(); err != nil {
			return err
		}
		tools[tool.ID] = tool
	}
	if err := bundle.Order.ValidateOrder(tools); err != nil {
		return err
	}
	for _, snapshot := range bundle.Backups {
		if err := backup.VerifyBackup(snapshot, tools); err != nil {
			return err
		}
	}
	if errors := audit.ValidateForTools(bundle.Audits, tools); len(errors) > 0 {
		return fmt.Errorf("invalid audit bundle: %s", strings.Join(errors, "; "))
	}
	return nil
}

func SortBundle(bundle *Bundle) {
	if bundle == nil {
		return
	}
	sort.SliceStable(bundle.Backups, func(i, j int) bool { return bundle.Backups[i].Revision < bundle.Backups[j].Revision })
	sort.SliceStable(bundle.Audits, func(i, j int) bool { return bundle.Audits[i].Revision < bundle.Audits[j].Revision })
}

func Summary(bundle Bundle) string {
	return fmt.Sprintf("tools=%d backups=%d audits=%d revision=%d", len(bundle.Tools), len(bundle.Backups), len(bundle.Audits), bundle.Order.Revision)
}
