package context

import "github.com/codebase-brain/internal/graph"

type Ranker struct {
	graph *graph.Graph
}

func NewRanker(g *graph.Graph) *Ranker {
	return &Ranker{graph: g}
}

func (r *Ranker) RankByRelevance(query string) []string {
	nodes := r.graph.Search(query)
	result := make([]string, len(nodes))
	for i, node := range nodes {
		result[i] = node.Name
	}
	return result
}

func (r *Ranker) RankByGraphDistance(filePath string, depth int) []string {
	return r.graph.GetRelated(filePath)
}

func (r *Ranker) RankByRecency(files []string) []string {
	return files
}