package indexer

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/codebase-brain/internal/config"
	"github.com/codebase-brain/internal/errors"
	"github.com/codebase-brain/internal/graph"
	"github.com/codebase-brain/internal/logger"
	"github.com/codebase-brain/internal/parser"
	"github.com/codebase-brain/internal/storage"
)

type SearchResult struct {
	FilePath string
	Match    string
	Line     int
}

type RelatedResult struct {
	FilePath string
	Type     string
	Distance int
}

type Indexer struct {
	cfg       *config.Config
	parsers   map[string]parser.Parser
	storage   storage.Storage
	projGraph *graph.Graph
	cache     map[string]*parser.ParseResult
	mu        sync.RWMutex
	workers   int
}

func NewIndexer(cfg *config.Config) (*Indexer, error) {
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}

	parsers := make(map[string]parser.Parser)
	parsers["go"] = parser.NewGoParser()
	parsers["js"] = parser.NewJSParser()
	parsers["ts"] = parser.NewJSParser()
	parsers["jsx"] = parser.NewJSParser()
	parsers["tsx"] = parser.NewJSParser()
	parsers["py"] = parser.NewPythonParser()
	parsers["rs"] = parser.NewRustParser()

	store, err := storage.NewBadgerStorage(cfg.ProjectPath)
	if err != nil {
		logger.Error("failed to create storage: %v", err)
		return nil, errors.Wrap(errors.ErrCodeStorageError, "failed to create storage", err)
	}

	return &Indexer{
		cfg:       cfg,
		parsers:   parsers,
		storage:   store,
		projGraph: graph.NewGraph(),
		cache:     make(map[string]*parser.ParseResult),
		workers:   cfg.Workers,
	}, nil
}

func (i *Indexer) Index(root string) (int, error) {
	logger.Info("starting indexing of %s", root)
	
	files, err := i.discoverFiles(root)
	if err != nil {
		logger.Error("failed to discover files: %v", err)
		return 0, errors.Wrap(errors.ErrCodeIndexError, "failed to discover files", err)
	}

	logger.Info("discovered %d files", len(files))

	count := 0
	errCount := 0
	
	// Use worker pool for concurrent parsing
	jobs := make(chan string, len(files))
	results := make(chan error, len(files))
	
	var wg sync.WaitGroup
	
	// Start workers
	for w := 0; w < i.workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				results <- i.processFile(file)
			}
		}()
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for completion
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		if result == nil {
			count++
		} else {
			errCount++
			logger.Debug("failed to process file: %v", result)
		}
	}

	logger.Info("indexing complete: %d files indexed, %d errors", count, errCount)

	// Save to storage
	if err := i.storage.SaveGraph(root, i.projGraph); err != nil {
		logger.Warn("failed to save graph: %v", err)
	}

	return count, nil
}

func (i *Indexer) processFile(file string) error {
	lang := detectLanguage(file)
	p, ok := i.parsers[lang]
	if !ok {
		return nil // Skip unsupported files
	}

	// Check cache first
	i.mu.RLock()
	if cached, ok := i.cache[file]; ok {
		i.mu.RUnlock()
		i.addToGraph(file, cached)
		return nil
	}
	i.mu.RUnlock()

	result, err := p.ParseFile(file)
	if err != nil {
		logger.Debug("failed to parse %s: %v", file, err)
		return err
	}

	// Cache result
	i.mu.Lock()
	i.cache[file] = result
	i.mu.Unlock()

	i.addToGraph(file, result)
	return nil
}

func (i *Indexer) addToGraph(file string, result *parser.ParseResult) {
	symbols := make([]graph.Symbol, len(result.Symbols))
	for idx, s := range result.Symbols {
		symbols[idx] = graph.Symbol{
			Name: s.Name,
			Type: s.Type,
			Location: graph.Location{
				File: s.Location.File,
				Line: s.Location.Line,
			},
		}
	}
	i.projGraph.AddFileNode(file, result.Language, result.Imports, result.Exports, symbols)
}

func (i *Indexer) Search(query string, limit int) []SearchResult {
	results := []SearchResult{}
	for _, node := range i.projGraph.Search(query) {
		if len(results) >= limit {
			break
		}
		results = append(results, SearchResult{
			FilePath: node.Name,
			Match:    query,
		})
	}
	return results
}

func (i *Indexer) GetRelated(filePath string, depth int) []RelatedResult {
	results := []RelatedResult{}
	visited := make(map[string]bool)
	i.collectRelated(filePath, depth, 0, visited, &results)
	return results
}

func (i *Indexer) collectRelated(filePath string, maxDepth, currentDepth int, visited map[string]bool, results *[]RelatedResult) {
	if currentDepth > maxDepth || visited[filePath] {
		return
	}
	visited[filePath] = true

	related := i.projGraph.GetRelated(filePath)
	for _, r := range related {
		*results = append(*results, RelatedResult{
			FilePath: r,
			Type:     "dependency",
			Distance: currentDepth + 1,
		})
		i.collectRelated(r, maxDepth, currentDepth+1, visited, results)
	}
}

func (i *Indexer) GetGraphOverview(limit int) []*graph.Node {
	return i.projGraph.GetNodes()[:min(limit, len(i.projGraph.GetNodes()))]
}

func (i *Indexer) GetFileContext(filePath, contextType string) string {
	node := i.projGraph.GetNode(filePath)
	if node == nil {
		return "File not found in index"
	}

	output := "File: " + filePath + "\n"
	output += "Type: " + string(node.Type) + "\n"

	if contextType == "all" || contextType == "imports" {
		imports := i.projGraph.GetImports(filePath)
		if len(imports) > 0 {
			output += "\nImports:\n"
			for _, imp := range imports {
				output += "  - " + imp + "\n"
			}
		}
	}

	if contextType == "all" || contextType == "exports" {
		exports := i.projGraph.GetExports(filePath)
		if len(exports) > 0 {
			output += "\nExports:\n"
			for _, exp := range exports {
				output += "  - " + exp + "\n"
			}
		}
	}

	if contextType == "all" || contextType == "dependencies" {
		related := i.projGraph.GetRelated(filePath)
		if len(related) > 0 {
			output += "\nDependencies:\n"
			for _, rel := range related {
				output += "  - " + rel + "\n"
			}
		}
	}

	return output
}

func (i *Indexer) discoverFiles(root string) ([]string, error) {
	var files []string
	extensions := []string{".go", ".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs", ".py", ".rs"}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warn("error walking path %s: %v", path, err)
			return nil // Continue walking despite errors
		}

		if info.IsDir() {
			dirName := info.Name()
			if strings.HasPrefix(dirName, ".") || 
			   dirName == "node_modules" || 
			   dirName == "vendor" ||
			   dirName == "__pycache__" ||
			   dirName == ".venv" ||
			   dirName == "venv" ||
			   dirName == "target" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		for _, e := range extensions {
			if ext == e {
				files = append(files, path)
				break
			}
		}
		return nil
	})

	return files, err
}

func (i *Indexer) ClearCache() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.cache = make(map[string]*parser.ParseResult)
}

func detectLanguage(path string) string {
	ext := filepath.Ext(path)
	langMap := map[string]string{
		".go":  "go",
		".js":  "js",
		".jsx": "js",
		".ts":  "ts",
		".tsx": "ts",
		".mjs": "js",
		".cjs": "js",
		".py":  "py",
		".rs":  "rs",
	}
	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}