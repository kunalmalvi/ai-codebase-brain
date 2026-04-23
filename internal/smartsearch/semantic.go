package smartsearch

import (
	"regexp"
	"strings"
)

type CodeConcept struct {
	Name        string
	Keywords    []string
	Patterns    []string
	Category    string
}

var Concepts = []CodeConcept{
	{
		Name:     "database",
		Keywords: []string{"db", "database", "sql", "query", "repository", "model", "schema", "migration", "table"},
		Patterns: []string{".*Query$", ".*Repository$", ".*Model$", ".*Service$"},
		Category: "data",
	},
	{
		Name:     "api",
		Keywords: []string{"api", "rest", "http", "endpoint", "route", "handler", "controller", "request", "response"},
		Patterns: []string{".*Handler$", ".*Controller$", ".*Route$", ".*Endpoint$", "/api/"},
		Category: "web",
	},
	{
		Name:     "authentication",
		Keywords: []string{"auth", "login", "jwt", "oauth", "session", "token", "password", "credential", "permission"},
		Patterns: []string{".*Auth.*", ".*Middleware$", ".*Guard$"},
		Category: "security",
	},
	{
		Name:     "testing",
		Keywords: []string{"test", "mock", "spec", "expect", "describe", "it", "should"},
		Patterns: []string{".*_test\\..*", ".*\\.test\\..*", ".*\\.spec\\..*", "test_.*"},
		Category: "quality",
	},
	{
		Name:     "config",
		Keywords: []string{"config", "setting", "env", "environment", "option", "flag"},
		Patterns: []string{".*config\\..*", ".*\\.env.*", "settings.*"},
		Category: "infrastructure",
	},
	{
		Name:     "error_handling",
		Keywords: []string{"error", "exception", "try", "catch", "throw", "raise", "fail"},
		Patterns: []string{".*Error$", ".*Exception$"},
		Category: "quality",
	},
	{
		Name:     "async",
		Keywords: []string{"async", "await", "promise", "future", "goroutine", "channel", "concurrent"},
		Patterns: []string{},
		Category: "performance",
	},
	{
		Name:     "crud",
		Keywords: []string{"create", "read", "update", "delete", "save", "fetch", "list", "get", "put", "post", "patch"},
		Patterns: []string{".*CRUD$", ".*Resource$"},
		Category: "data",
	},
	{
		Name:     "validation",
		Keywords: []string{"validate", "validate", "schema", "type", "check", "verify", "assert"},
		Patterns: []string{".*Validator$", ".*Validation$"},
		Category: "quality",
	},
	{
		Name:     "logging",
		Keywords: []string{"log", "debug", "info", "warn", "error", "trace"},
		Patterns: []string{},
		Category: "infrastructure",
	},
	{
		Name:     "cache",
		Keywords: []string{"cache", "redis", "memcached", "store", "session"},
		Patterns: []string{".*Cache$", ".*Store$"},
		Category: "performance",
	},
	{
		Name:     "message_queue",
		Keywords: []string{"queue", "message", "kafka", "rabbitmq", "pubsub", "event", "listener"},
		Patterns: []string{".*Queue$", ".*Consumer$", ".*Producer$"},
		Category: "integration",
	},
}

type SemanticSearch struct {
	concepts []CodeConcept
}

func NewSemanticSearch() *SemanticSearch {
	return &SemanticSearch{
		concepts: Concepts,
	}
}

func (s *SemanticSearch) DetectConcepts(filePath, content string) []string {
	lowerContent := strings.ToLower(content)
	lowerPath := strings.ToLower(filePath)
	
	var detected []string
	seen := make(map[string]bool)
	
	for _, concept := range s.concepts {
		// Check keywords
		for _, kw := range concept.Keywords {
			if strings.Contains(lowerContent, kw) || strings.Contains(lowerPath, kw) {
				if !seen[concept.Name] {
					detected = append(detected, concept.Name)
					seen[concept.Name] = true
				}
			}
		}
		
		// Check patterns
		for _, pattern := range concept.Patterns {
			matched, _ := regexp.MatchString(pattern, filePath)
			if matched && !seen[concept.Name] {
				detected = append(detected, concept.Name)
				seen[concept.Name] = true
			}
		}
	}
	
	return detected
}

func (s *SemanticSearch) SearchByConcept(query string) []string {
	lowerQuery := strings.ToLower(query)
	
	var results []string
	for _, concept := range s.concepts {
		// Check if query matches concept name
		if strings.Contains(concept.Name, lowerQuery) {
			results = append(results, concept.Name)
			continue
		}
		
		// Check if query matches any keyword
		for _, kw := range concept.Keywords {
			if strings.Contains(kw, lowerQuery) || strings.Contains(lowerQuery, kw) {
				results = append(results, concept.Name)
				break
			}
		}
	}
	
	return results
}

func (s *SemanticSearch) GetRelatedConcepts(conceptName string) []string {
	for _, concept := range s.concepts {
		if concept.Name == conceptName {
			// Return concepts from same category
			var related []string
			for _, c := range s.concepts {
				if c.Category == concept.Category && c.Name != concept.Name {
					related = append(related, c.Name)
				}
			}
			return related
		}
	}
	return nil
}

func (s *SemanticSearch) GetCategory(conceptName string) string {
	for _, concept := range s.concepts {
		if concept.Name == conceptName {
			return concept.Category
		}
	}
	return ""
}