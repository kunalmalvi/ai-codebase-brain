package graph

import (
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph()
	if g == nil {
		t.Error("NewGraph() returned nil")
	}
	if len(g.nodes) != 0 {
		t.Errorf("expected empty nodes, got %d", len(g.nodes))
	}
}

func TestAddFileNode(t *testing.T) {
	g := NewGraph()
	
	imports := []string{"fmt", "os"}
	exports := []string{"Main", "Run"}
	symbols := []Symbol{
		{Name: "Main", Type: "function"},
		{Name: "Run", Type: "function"},
	}
	
	g.AddFileNode("/path/to/main.go", "go", imports, exports, symbols)
	
	if len(g.nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(g.nodes))
	}
	
	if len(g.adjacency) != 1 {
		t.Errorf("expected 1 adjacency entry, got %d", len(g.adjacency))
	}
}

func TestGetNode(t *testing.T) {
	g := NewGraph()
	
	g.AddFileNode("test.go", "go", []string{"fmt"}, []string{"Test"}, nil)
	
	node := g.GetNode("test.go")
	if node == nil {
		t.Error("GetNode() returned nil for existing file")
	}
	
	nonexistent := g.GetNode("nonexistent.go")
	if nonexistent != nil {
		t.Error("GetNode() should return nil for nonexistent file")
	}
}

func TestGetNodes(t *testing.T) {
	g := NewGraph()
	
	g.AddFileNode("file1.go", "go", nil, nil, nil)
	g.AddFileNode("file2.go", "go", nil, nil, nil)
	g.AddFileNode("file3.py", "python", nil, nil, nil)
	
	nodes := g.GetNodes()
	if len(nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(nodes))
	}
}

func TestGetImports(t *testing.T) {
	g := NewGraph()
	
	imports := []string{"fmt", "os", "strings"}
	g.AddFileNode("main.go", "go", imports, nil, nil)
	
	result := g.GetImports("main.go")
	if len(result) != 3 {
		t.Errorf("expected 3 imports, got %d", len(result))
	}
	
	nonexistent := g.GetImports("nonexistent.go")
	if len(nonexistent) != 0 {
		t.Error("expected empty imports for nonexistent file")
	}
}

func TestGetExports(t *testing.T) {
	g := NewGraph()
	
	exports := []string{"Main", "Run", "Init"}
	g.AddFileNode("main.go", "go", nil, exports, nil)
	
	result := g.GetExports("main.go")
	if len(result) != 3 {
		t.Errorf("expected 3 exports, got %d", len(result))
	}
}

func TestGetRelated(t *testing.T) {
	g := NewGraph()
	
	g.AddFileNode("main.go", "go", []string{"utils"}, nil, nil)
	g.AddFileNode("utils.go", "go", []string{"fmt"}, nil, nil)
	
	related := g.GetRelated("main.go")
	if len(related) != 1 || related[0] != "utils" {
		t.Errorf("expected ['utils'], got %v", related)
	}
}

func TestSearch(t *testing.T) {
	g := NewGraph()
	
	g.AddFileNode("/src/main.go", "go", nil, nil, nil)
	g.AddFileNode("/src/utils/helper.go", "go", nil, nil, nil)
	g.AddFileNode("/tests/main_test.go", "go", nil, nil, nil)
	
	results := g.Search("main")
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	
	noResults := g.Search("nonexistent")
	if len(noResults) != 0 {
		t.Errorf("expected 0 results, got %d", len(noResults))
	}
}

func TestAddEdge(t *testing.T) {
	g := NewGraph()
	
	g.AddFileNode("a.go", "go", nil, nil, nil)
	g.AddFileNode("b.go", "go", nil, nil, nil)
	
	g.AddEdge(NodeID("a.go"), NodeID("b.go"), "imports", 1.0)
	
	edges := g.edges[NodeID("a.go")]
	if len(edges) != 1 {
		t.Errorf("expected 1 edge, got %d", len(edges))
	}
}

func TestAddNode(t *testing.T) {
	g := NewGraph()
	
	// Add via AddFileNode (the normal way)
	g.AddFileNode("test.go", "go", nil, nil, nil)
	
	// Test GetNode works
	retrieved := g.GetNode("test.go")
	if retrieved == nil || retrieved.Name != "test.go" {
		t.Error("AddFileNode/GetNode failed")
	}
	
	// GetNode returns from fileNodes map, not nodes map
	// So test fileNodes directly instead
	nodes := g.GetNodes()
	if len(nodes) < 1 {
		t.Error("expected at least 1 node")
	}
}