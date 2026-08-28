package backup

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"example.com/toolnav/model"
	"example.com/toolnav/store"
)

type Service struct {
	Store *store.Store
}

func NewService(s *store.Store) *Service { return &Service{Store: s} }

func CreateBackup(s *store.Store, order model.ToolOrder, tools map[string]model.Tool, label string) (model.BackupSnapshot, error) {
	if err := order.ValidateOrder(tools); err != nil {
		return model.BackupSnapshot{}, err
	}
	ids := append([]string(nil), order.ToolIDs...)
	snapshot := model.BackupSnapshot{ID: fmt.Sprintf("backup-%d", order.Revision), CreatedAtLabel: label, Revision: order.Revision, ToolIDs: ids}
	snapshot.Checksum = Checksum(snapshot.ToolIDs)
	if err := s.SaveBackup(snapshot); err != nil {
		return model.BackupSnapshot{}, err
	}
	return snapshot, nil
}

func (s *Service) Create(order model.ToolOrder, tools map[string]model.Tool, label string) (model.BackupSnapshot, error) {
	return CreateBackup(s.Store, order, tools, label)
}

func VerifyBackup(snapshot model.BackupSnapshot, tools map[string]model.Tool) error {
	if len(snapshot.ToolIDs) != len(tools) {
		return fmt.Errorf("backup has %d ids, expected %d", len(snapshot.ToolIDs), len(tools))
	}
	seen := make(map[string]bool, len(snapshot.ToolIDs))
	for _, id := range snapshot.ToolIDs {
		if seen[id] {
			return fmt.Errorf("backup contains duplicate %q", id)
		}
		if _, ok := tools[id]; !ok {
			return fmt.Errorf("backup contains unknown %q", id)
		}
		seen[id] = true
	}
	if snapshot.Checksum != Checksum(snapshot.ToolIDs) {
		return fmt.Errorf("backup checksum mismatch")
	}
	return nil
}

func Checksum(ids []string) string {
	hash := sha256.Sum256([]byte(strings.Join(ids, "|")))
	return hex.EncodeToString(hash[:])
}

func OrderedLabels(snapshot model.BackupSnapshot) string {
	return strings.Join(snapshot.ToolIDs, " > ")
}

func (s *Service) Latest() (model.BackupSnapshot, error) { return s.Store.LatestBackup() }
