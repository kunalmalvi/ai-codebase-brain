package storage

import "github.com/codebase-brain/internal/graph"

type Storage interface {
	SaveGraph(projectID string, g *graph.Graph) error
	LoadGraph(projectID string) (*graph.Graph, error)
	Close() error
}