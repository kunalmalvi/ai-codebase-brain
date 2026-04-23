package smartsearch

import (
	"fmt"
	"strings"
)

type FileSummary struct {
	FilePath     string   `json:"file_path"`
	Language     string   `json:"language"`
	Purpose      string   `json:"purpose"`
	MainFeatures []string `json:"main_features"`
	Concepts     []string `json:"concepts"`
	Dependencies []string `json:"dependencies"`
	Exports      []string `json:"exports"`
	Complexity   string   `json:"complexity"` // "simple", "moderate", "complex"
	LineCount    int      `json:"line_count"`
}

type Summarizer struct {
	semantic *SemanticSearch
}

func NewSummarizer() *Summarizer {
	return &Summarizer{
		semantic: NewSemanticSearch(),
	}
}

func (s *Summarizer) SummarizeFile(filePath, content string, imports, exports []string) *FileSummary {
	lines := strings.Split(content, "\n")
	lineCount := len(lines)
	
	summary := &FileSummary{
		FilePath:     filePath,
		LineCount:    lineCount,
		Dependencies: imports,
		Exports:      exports,
	}
	
	// Detect language from file path
	summary.Language = detectLanguageFromPath(filePath)
	
	// Detect concepts
	summary.Concepts = s.semantic.DetectConcepts(filePath, content)
	
	// Determine purpose based on concepts and content
	summary.Purpose = s.determinePurpose(summary.Concepts, content)
	
	// Extract main features
	summary.MainFeatures = s.extractFeatures(content, summary.Concepts)
	
	// Determine complexity
	summary.Complexity = s.assessComplexity(lineCount, len(imports), len(exports), len(summary.Concepts))
	
	return summary
}

func (s *Summarizer) determinePurpose(concepts []string, content string) string {
	if len(concepts) == 0 {
		// Try to guess from content
		lower := strings.ToLower(content)
		if strings.Contains(lower, "package main") {
			return "Main application entry point"
		}
		if strings.Contains(lower, "function") || strings.Contains(lower, "def ") {
			return "Utility/helper module"
		}
		return "General purpose file"
	}
	
	// Primary concept determines purpose
	primary := concepts[0]
	switch primary {
	case "api":
		return "API endpoint/handler module"
	case "database":
		return "Data access/repository layer"
	case "authentication":
		return "Authentication/authorization module"
	case "config":
		return "Configuration module"
	case "testing":
		return "Test suite"
	case "crud":
		return "CRUD operations module"
	case "cache":
		return "Caching layer"
	case "message_queue":
		return "Message queue/event handling"
	default:
		return fmt.Sprintf("%s module", primary)
	}
}

func (s *Summarizer) extractFeatures(content string, concepts []string) []string {
	var features []string
	lines := strings.Split(content, "\n")
	
	// Look for function/method definitions
	funcCount := 0
	classCount := 0
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "func ") || strings.HasPrefix(trimmed, "def ") || 
		   strings.HasPrefix(trimmed, "function ") || strings.HasPrefix(trimmed, "fn ") {
			funcCount++
		}
		if strings.HasPrefix(trimmed, "class ") || strings.HasPrefix(trimmed, "struct ") ||
		   strings.HasPrefix(trimmed, "type ") {
			classCount++
		}
	}
	
	if funcCount > 0 {
		features = append(features, fmt.Sprintf("%d function(s)", funcCount))
	}
	if classCount > 0 {
		features = append(features, fmt.Sprintf("%d type(s) defined", classCount))
	}
	
	// Add concept-based features
	for _, concept := range concepts {
		switch concept {
		case "api":
			features = append(features, "Handles HTTP requests")
		case "database":
			features = append(features, "Data persistence operations")
		case "authentication":
			features = append(features, "Auth/nauthorization logic")
		case "validation":
			features = append(features, "Input validation")
		case "async":
			features = append(features, "Async/await operations")
		case "logging":
			features = append(features, "Logging capabilities")
		}
	}
	
	if len(features) == 0 {
		features = append(features, "General utility")
	}
	
	return features
}

func (s *Summarizer) assessComplexity(lineCount, importCount, exportCount, conceptCount int) string {
	score := 0
	
	// Lines of code
	if lineCount > 500 {
		score += 3
	} else if lineCount > 200 {
		score += 2
	} else if lineCount > 100 {
		score += 1
	}
	
	// Dependencies
	if importCount > 20 {
		score += 3
	} else if importCount > 10 {
		score += 2
	} else if importCount > 5 {
		score += 1
	}
	
	// Exports/Public API
	if exportCount > 20 {
		score += 3
	} else if exportCount > 10 {
		score += 2
	} else if exportCount > 5 {
		score += 1
	}
	
	// Concepts (indicates multiple responsibilities)
	if conceptCount > 4 {
		score += 3
	} else if conceptCount > 2 {
		score += 2
	} else if conceptCount > 0 {
		score += 1
	}
	
	// Determine complexity
	if score >= 8 {
		return "complex"
	} else if score >= 4 {
		return "moderate"
	}
	return "simple"
}

func detectLanguageFromPath(path string) string {
	ext := ""
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			ext = path[i:]
			break
		}
	}
	
	langMap := map[string]string{
		".go":     "Go",
		".js":     "JavaScript",
		".ts":     "TypeScript",
		".jsx":    "React (JSX)",
		".tsx":    "React (TSX)",
		".py":     "Python",
		".rs":     "Rust",
		".java":   "Java",
		".rb":     "Ruby",
		".php":    "PHP",
		".cs":     "C#",
		".cpp":    "C++",
		".c":      "C",
		".swift":  "Swift",
		".kt":     "Kotlin",
	}
	
	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return "Unknown"
}

func (s *Summarizer) GenerateQuickSummary(summary *FileSummary) string {
	output := fmt.Sprintf("📄 %s\n", summary.FilePath)
	output += fmt.Sprintf("   Language: %s | Lines: %d | Complexity: %s\n\n", 
		summary.Language, summary.LineCount, summary.Complexity)
	
	output += fmt.Sprintf("🎯 Purpose: %s\n\n", summary.Purpose)
	
	if len(summary.MainFeatures) > 0 {
		output += "✨ Features:\n"
		for _, f := range summary.MainFeatures {
			output += fmt.Sprintf("   • %s\n", f)
		}
		output += "\n"
	}
	
	if len(summary.Concepts) > 0 {
		output += "🔍 Detected Concepts: "
		output += strings.Join(summary.Concepts, ", ")
		output += "\n"
	}
	
	return output
}