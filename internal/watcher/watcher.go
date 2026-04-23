package watcher

import (
	"github.com/codebase-brain/internal/indexer"
	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	indexer *indexer.Indexer
	watcher *fsnotify.Watcher
}

func NewFileWatcher(idx *indexer.Indexer) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FileWatcher{
		indexer: idx,
		watcher: watcher,
	}, nil
}

func (w *FileWatcher) Start(root string) error {
	return w.watcher.Add(root)
}

func (w *FileWatcher) OnChange(event fsnotify.Event) {
	if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
		w.indexer.Index(event.Name)
	}
	if event.Has(fsnotify.Remove) {
	}
}

func (w *FileWatcher) Close() error {
	return w.watcher.Close()
}