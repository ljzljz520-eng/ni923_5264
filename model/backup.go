package model

type BackupSnapshot struct {
	ID             string   `json:"id"`
	CreatedAtLabel string   `json:"created_at_label"`
	Revision       int      `json:"revision"`
	ToolIDs        []string `json:"tool_ids"`
	Checksum       string   `json:"checksum"`
}

func (s BackupSnapshot) Clone() BackupSnapshot {
	return BackupSnapshot{ID: s.ID, CreatedAtLabel: s.CreatedAtLabel, Revision: s.Revision, ToolIDs: append([]string(nil), s.ToolIDs...), Checksum: s.Checksum}
}
