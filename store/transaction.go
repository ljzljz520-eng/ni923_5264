package store

import (
	"fmt"

	"example.com/toolnav/model"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) DeleteTool(id string) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(bucketTools).Delete([]byte(id)) })
}

func (s *Store) SaveImport(batch model.ImportBatch) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	data, err := encode(batch)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(bucketImports).Put([]byte(batch.ID), data) })
}

func (s *Store) ClearBackups() error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store is closed")
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketBackups)
		keys := make([][]byte, 0)
		if err := bucket.ForEach(func(key, _ []byte) error {
			keys = append(keys, append([]byte(nil), key...))
			return nil
		}); err != nil {
			return err
		}
		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}
