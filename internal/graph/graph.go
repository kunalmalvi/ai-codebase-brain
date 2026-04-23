package graph

import "fmt"

type NodeID string
type NodeType string

const (
	NodeTypeFile     NodeType = "file"
	NodeTypeFunction NodeType = "function"
	NodeTypeClass    NodeType = "class"
	NodeTypeVariable NodeType = "variable"
	NodeTypeImport   NodeType = "import"
)

type Location struct {
	File     string
	Line     int
	Column   int
	EndLine  int
	EndColumn int
}

type Node struct {
	ID       NodeID
	Type     NodeType
	Name     string
	Location Location
	Metadata map[string]interface{}
}

type Edge struct {
	From   NodeID
	To     NodeID
	Type   string
	Weight float64
}

type Symbol struct {
	Name     string
	Type     string
	Location Location
}

type Graph struct {
	nodes map[NodeID]*Node
	edges map[NodeID][]NodeID
	adjacency map[string][]string
	fileNodes map[string]*Node
}

func NewGraph() *Graph {
	return &Graph{
		nodes:     make(map[NodeID]*Node),
		edges:     make(map[NodeID][]NodeID),
		adjacency: make(map[string][]string),
		fileNodes: make(map[string]*Node),
	}
}

func (g *Graph) AddFileNode(path, language string, imports, exports []string, symbols []Symbol) {
	nodeID := NodeID(path)
	node := &Node{
		ID:   nodeID,
		Type: NodeTypeFile,
		Name: path,
		Location: Location{
			File: path,
		},
		Metadata: map[string]interface{}{
			"language": language,
			"imports":  imports,
			"exports":  exports,
			"symbols":  symbols,
		},
	}
	g.nodes[nodeID] = node
	g.fileNodes[path] = node

	for _, imp := range imports {
		g.adjacency[path] = append(g.adjacency[path], imp)
	}
}

func (g *Graph) AddNode(node *Node) {
	g.nodes[node.ID] = node
}

func (g *Graph) AddEdge(from, to NodeID, edgeType string, weight float64) {
	g.edges[from] = append(g.edges[from], to)
	_ = edgeType
	_ = weight
}

func (g *Graph) GetNode(id string) *Node {
	return g.fileNodes[id]
}

func (g *Graph) GetNodes() []*Node {
	result := make([]*Node, 0, len(g.fileNodes))
	for _, node := range g.fileNodes {
		result = append(result, node)
	}
	return result
}

func (g *Graph) GetRelated(filePath string) []string {
	return g.adjacency[filePath]
}

func (g *Graph) GetImports(filePath string) []string {
	if node, ok := g.fileNodes[filePath]; ok {
		if imports, ok := node.Metadata["imports"].([]string); ok {
			return imports
		}
	}
	return []string{}
}

func (g *Graph) GetExports(filePath string) []string {
	if node, ok := g.fileNodes[filePath]; ok {
		if exports, ok := node.Metadata["exports"].([]string); ok {
			return exports
		}
	}
	return []string{}
}

func (g *Graph) Search(query string) []*Node {
	results := []*Node{}
	for _, node := range g.fileNodes {
		if contains(node.Name, query) {
			results = append(results, node)
		}
	}
	return results
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (g *Graph) String() string {
	return fmt.Sprintf("Graph{nodes: %d}", len(g.nodes))
}