package smartsearch

import (
	"fmt"
	"os"
	"strings"
)

type ContextSuggestion struct {
	FilePath     string   `json:"file_path"`
	Reason       string   `json:"reason"`
	Priority     int      `json:"priority"`
	SuggestionType string  `json:"type"` // "related", "dependency", "test", "config"
}

type ContextAdvisor struct {
	semantic *SemanticSearch
	summarizer *Summarizer
}

func NewContextAdvisor() *ContextAdvisor {
	return &ContextAdvisor{
		semantic:   NewSemanticSearch(),
		summarizer: NewSummarizer(),
	}
}

func (a *ContextAdvisor) GetSuggestionsForTask(filePath, content, taskDescription string) []ContextSuggestion {
	var suggestions []ContextSuggestion
	
	// Detect concepts in current file
	concepts := a.semantic.DetectConcepts(filePath, content)
	
	// Get related concepts
	for _, concept := range concepts {
		related := a.semantic.GetRelatedConcepts(concept)
		for _, r := range related {
			suggestions = append(suggestions, ContextSuggestion{
				FilePath:     "related:" + r,
				Reason:       "Related to detected concept: " + concept,
				Priority:     3,
				SuggestionType: "related",
			})
		}
	}
	
	// Detect what the user is trying to do
	taskLower := strings.ToLower(taskDescription)
	
	// Auth-related tasks
	if strings.Contains(taskLower, "auth") || strings.Contains(taskLower, "login") || strings.Contains(taskLower, "password") {
		suggestions = append(suggestions, ContextSuggestion{
			FilePath:     "concept:authentication",
			Reason:       "Task involves authentication",
			Priority:     10,
			SuggestionType: "context",
		})
		suggestions = append(suggestions, ContextSuggestion{
			FilePath:     "concept:security",
			Reason:       "May need security-related context",
			Priority:     9,
			SuggestionType: "context",
		})
	}
	
	// Database tasks
	if strings.Contains(taskLower, "database") || strings.Contains(taskLower, "db") || 
	   strings.Contains(taskLower, "save") || strings.Contains(taskLower, "query") {
		suggestions = append(suggestions, ContextSuggestion{
			FilePath:     "concept:database",
			Reason:       "Task involves database operations",
			Priority:     10,
			SuggestionType: "context",
		})
	}
	
	// API tasks
	if strings.Contains(taskLower, "api") || strings.Contains(taskLower, "endpoint") || 
	   strings.Contains(taskLower, "request") || strings.Contains(taskLower, "http") {
		suggestions = append(suggestions, ContextSuggestion{
			FilePath:     "concept:api",
			Reason:       "Task involves API endpoints",
			Priority:     10,
			SuggestionType: "context",
		})
	}
	
	// Testing tasks
	if strings.Contains(taskLower, "test") || strings.Contains(taskLower, "mock") {
		// Suggest related test files
		suggestions = append(suggestions, ContextSuggestion{
			FilePath:     "type:test",
			Reason:       "Task involves testing",
			Priority:     10,
			SuggestionType: "test",
		})
	}
	
	// Refactoring tasks
	if strings.Contains(taskLower, "refactor") || strings.Contains(taskLower, "improve") || 
	   strings.Contains(taskLower, "clean") {
		suggestions = append(suggestions, ContextSuggestion{
			FilePath:     "concept:error_handling",
			Reason:       "May need to review error handling",
			Priority:     8,
			SuggestionType: "context",
		})
		suggestions = append(suggestions, ContextSuggestion{
			FilePath:     "type:refactoring",
			Reason:       "Running refactoring analysis",
			Priority:     10,
			SuggestionType: "refactoring",
		})
	}
	
	// Performance tasks
	if strings.Contains(taskLower, "performance") || strings.Contains(taskLower, "optimize") || 
	   strings.Contains(taskLower, "speed") {
		suggestions = append(suggestions, ContextSuggestion{
			FilePath:     "concept:cache",
			Reason:       "May benefit from caching context",
			Priority:     9,
			SuggestionType: "context",
		})
		suggestions = append(suggestions, ContextSuggestion{
			FilePath:     "concept:async",
			Reason:       "May need async/concurrent patterns",
			Priority:     8,
			SuggestionType: "context",
		})
	}
	
	return suggestions
}

func (a *ContextAdvisor) RecommendContextFiles(filePath, task string, projectFiles []string) []ContextSuggestion {
	content := ""
	if data, err := os.ReadFile(filePath); err == nil {
		content = string(data)
	}
	
	// Get suggestions
	suggestions := a.GetSuggestionsForTask(filePath, content, task)
	
	// Filter project files based on suggestions
	var recommendations []ContextSuggestion
	
	for _, s := range suggestions {
		if strings.HasPrefix(s.FilePath, "concept:") {
			concept := strings.TrimPrefix(s.FilePath, "concept:")
			// Find files with this concept
			for _, pf := range projectFiles {
				pfContent, _ := os.ReadFile(pf)
				pfConcepts := a.semantic.DetectConcepts(pf, string(pfContent))
				for _, pc := range pfConcepts {
					if pc == concept {
						recommendations = append(recommendations, ContextSuggestion{
							FilePath:     pf,
							Reason:       s.Reason,
							Priority:     s.Priority,
							SuggestionType: "related",
						})
					}
				}
			}
		}
	}
	
	return recommendations
}

func (a *ContextAdvisor) GenerateTaskSummary(suggestions []ContextSuggestion) string {
	if len(suggestions) == 0 {
		return "No specific context suggestions for this task."
	}
	
	output := "💡 Context Suggestions for this task:\n\n"
	
	// Group by type
	byType := make(map[string][]ContextSuggestion)
	for _, s := range suggestions {
		byType[s.SuggestionType] = append(byType[s.SuggestionType], s)
	}
	
	for _, suggestionType := range []string{"context", "related", "test", "refactoring"} {
		items, ok := byType[suggestionType]
		if !ok {
			continue
		}
		
		output += "📌 " + suggestionType + ":\n"
		seen := make(map[string]bool)
		for _, s := range items {
			if !seen[s.FilePath] {
				output += fmt.Sprintf("   • %s\n      Reason: %s\n", s.FilePath, s.Reason)
				seen[s.FilePath] = true
			}
		}
		output += "\n"
	}
	
	return output
}

func (a *ContextAdvisor) GetOptimalContextWindow(taskDescription string, availableFiles []string) []string {
	suggestions := a.GetSuggestionsForTask("", "", taskDescription)
	
	// Score each file
	type scoredFile struct {
		path   string
		score  int
	}
	
	var scored []scoredFile
	
	for _, f := range availableFiles {
		score := 0
		
		// Check if file matches any suggested concepts
		for _, s := range suggestions {
			if strings.HasPrefix(s.FilePath, "concept:") {
				concept := strings.TrimPrefix(s.FilePath, "concept:")
				content, _ := os.ReadFile(f)
				concepts := a.semantic.DetectConcepts(f, string(content))
				for _, c := range concepts {
					if c == concept {
						score += s.Priority
					}
					// Also check related
					related := a.semantic.GetRelatedConcepts(c)
					for _, r := range related {
						if r == concept {
							score += s.Priority / 2
						}
					}
				}
			}
		}
		
		// Boost recently modified files
		
		if score > 0 {
			scored = append(scored, scoredFile{path: f, score: score})
		}
	}
	
	// Sort by score
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}
	
	// Return top results
	result := make([]string, min(10, len(scored)))
	for i, s := range scored {
		if i >= 10 {
			break
		}
		result[i] = s.path
	}
	
	return result
}