package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codebase-brain/internal/graph"
	"github.com/dgraph-io/badger/v4"
)

type BadgerStorage struct {
	db *badger.DB
}

func NewBadgerStorage(projectPath string) (*BadgerStorage, error) {
	dir := filepath.Join(projectPath, ".ai-coder-context")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage dir: %w", err)
	}

	opts := badger.DefaultOptions(dir).
		WithLogger(nil)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger: %w", err)
	}

	return &BadgerStorage{db: db}, nil
}

func (s *BadgerStorage) SaveGraph(projectID string, g *graph.Graph) error {
	return s.db.Update(func(txn *badger.Txn) error {
		nodes := g.GetNodes()
		for _, node := range nodes {
			key := []byte(fmt.Sprintf("node:%s:%s", projectID, node.Name))
			val := []byte(fmt.Sprintf("%v", node.Metadata))
			if err := txn.Set(key, val); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BadgerStorage) LoadGraph(projectID string) (*graph.Graph, error) {
	g := graph.NewGraph()

	err := s.db.View(func(txn *badger.Txn) error {
		prefix := []byte(fmt.Sprintf("node:%s:", projectID))
		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		for iter.Seek(prefix); iter.ValidForPrefix(prefix); iter.Next() {
			item := iter.Item()
			key := string(item.Key())
			val, err := item.ValueCopy(nil)
			if err != nil {
				continue
			}

			nodeName := key[len(prefix):]
			g.AddFileNode(nodeName, "unknown", []string{}, []string{}, nil)
			_ = val
		}
		return nil
	})

	return g, err
}

func (s *BadgerStorage) Close() error {
	return s.db.Close()
}