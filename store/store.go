package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"example.com/toolnav/model"
	bolt "go.etcd.io/bbolt"
)

var (
	bucketTools   = []byte("tools")
	bucketOrder   = []byte("order")
	bucketImports = []byte("imports")
	bucketBackups = []byte("backups")
	bucketAudits  = []byte("audits")
	keyCurrent    = []byte("current")
)

type Store struct {
	path string
	db   *bolt.DB
}

type CatalogState struct {
	Tools   map[string]model.Tool
	Order   model.ToolOrder
	Imports []model.ImportBatch
	Backups []model.BackupSnapshot
	Audits  []model.AuditEvent
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("store path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := bolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, err
	}
	s := &Store{path: path, db: db}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Path() string { return s.path }

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range [][]byte{bucketTools, bucketOrder, bucketImports, bucketBackups, bucketAudits} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return fmt.Errorf("empty stored value")
	}
	return json.Unmarshal(data, target)
}

func (s *Store) SaveCatalog(tools map[string]model.Tool, order model.ToolOrder, batch *model.ImportBatch) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		toolsBucket := tx.Bucket(bucketTools)
		for id, tool := range tools {
			data, err := encode(tool)
			if err != nil {
				return err
			}
			if err := toolsBucket.Put([]byte(id), data); err != nil {
				return err
			}
		}
		orderData, err := encode(order)
		if err != nil {
			return err
		}
		if err := tx.Bucket(bucketOrder).Put(keyCurrent, orderData); err != nil {
			return err
		}
		if batch != nil {
			batchData, err := encode(*batch)
			if err != nil {
				return err
			}
			return tx.Bucket(bucketImports).Put([]byte(batch.ID), batchData)
		}
		return nil
	})
}

func (s *Store) SaveTool(tool model.Tool) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	data, err := encode(tool)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(bucketTools).Put([]byte(tool.ID), data) })
}

func (s *Store) SaveOrder(order model.ToolOrder) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	data, err := encode(order)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(bucketOrder).Put(keyCurrent, data) })
}

func (s *Store) SaveBackup(snapshot model.BackupSnapshot) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	data, err := encode(snapshot)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(bucketBackups).Put([]byte(snapshot.ID), data) })
}

func (s *Store) SaveAudit(event model.AuditEvent) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	data, err := encode(event)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(bucketAudits).Put([]byte(event.ID), data) })
}

func (s *Store) LoadCatalog() (CatalogState, error) {
	if s == nil || s.db == nil {
		return CatalogState{}, fmt.Errorf("store is closed")
	}
	state := CatalogState{Tools: make(map[string]model.Tool)}
	err := s.db.View(func(tx *bolt.Tx) error {
		return readTools(tx, state.Tools)
	})
	if err != nil {
		return CatalogState{}, err
	}
	if err := s.db.View(func(tx *bolt.Tx) error { return readOrder(tx, &state.Order) }); err != nil {
		return CatalogState{}, err
	}
	if err := s.db.View(func(tx *bolt.Tx) error { return readImports(tx, &state.Imports) }); err != nil {
		return CatalogState{}, err
	}
	if err := s.db.View(func(tx *bolt.Tx) error { return readBackups(tx, &state.Backups) }); err != nil {
		return CatalogState{}, err
	}
	if err := s.db.View(func(tx *bolt.Tx) error { return readAudits(tx, &state.Audits) }); err != nil {
		return CatalogState{}, err
	}
	return state, nil
}

func readTools(tx *bolt.Tx, target map[string]model.Tool) error {
	return tx.Bucket(bucketTools).ForEach(func(k, v []byte) error {
		var tool model.Tool
		if err := decode(v, &tool); err != nil {
			return err
		}
		target[string(k)] = tool
		return nil
	})
}

func readOrder(tx *bolt.Tx, target *model.ToolOrder) error {
	data := tx.Bucket(bucketOrder).Get(keyCurrent)
	if len(data) == 0 {
		return nil
	}
	return decode(data, target)
}

func readImports(tx *bolt.Tx, target *[]model.ImportBatch) error {
	return tx.Bucket(bucketImports).ForEach(func(_, v []byte) error {
		var batch model.ImportBatch
		if err := decode(v, &batch); err != nil {
			return err
		}
		*target = append(*target, batch)
		return nil
	})
}

func readBackups(tx *bolt.Tx, target *[]model.BackupSnapshot) error {
	return tx.Bucket(bucketBackups).ForEach(func(_, v []byte) error {
		var snapshot model.BackupSnapshot
		if err := decode(v, &snapshot); err != nil {
			return err
		}
		*target = append(*target, snapshot)
		return nil
	})
}

func readAudits(tx *bolt.Tx, target *[]model.AuditEvent) error {
	return tx.Bucket(bucketAudits).ForEach(func(_, v []byte) error {
		var event model.AuditEvent
		if err := decode(v, &event); err != nil {
			return err
		}
		*target = append(*target, event)
		return nil
	})
}

func (s *Store) LatestBackup() (model.BackupSnapshot, error) {
	state, err := s.LoadCatalog()
	if err != nil {
		return model.BackupSnapshot{}, err
	}
	if len(state.Backups) == 0 {
		return model.BackupSnapshot{}, fmt.Errorf("no backups available")
	}
	return state.Backups[len(state.Backups)-1], nil
}
