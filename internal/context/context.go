package context

import "github.com/codebase-brain/internal/graph"

type ContextRequest struct {
	Query     string
	Hint      string
	FilePath  string
	Limit     int
}

type ContextItem struct {
	FilePath  string
	Relevance float64
	Type      string
}

type ContextService struct {
	graph   *graph.Graph
	ranker  *Ranker
}

func NewContextService(g *graph.Graph) *ContextService {
	return &ContextService{
		graph:  g,
		ranker: NewRanker(g),
	}
}

func (c *ContextService) GetContext(req ContextRequest) []ContextItem {
	results := []ContextItem{}

	if req.FilePath != "" {
		related := c.graph.GetRelated(req.FilePath)
		for _, rel := range related {
			results = append(results, ContextItem{
				FilePath:  rel,
				Relevance: 1.0,
				Type:      "dependency",
			})
		}
	}

	if req.Query != "" {
		matched := c.graph.Search(req.Query)
		for _, node := range matched {
			results = append(results, ContextItem{
				FilePath:  node.Name,
				Relevance: 0.8,
				Type:      string(node.Type),
			})
		}
	}

	if req.Limit > 0 && len(results) > req.Limit {
		results = results[:req.Limit]
	}

	return results
}