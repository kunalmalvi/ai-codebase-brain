package smartsearch

import (
	"fmt"
	"regexp"
	"strings"
)

type Refactoring struct {
	Type        string   `json:"type"`
	Severity    string   `json:"severity"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Line        int      `json:"line"`
}

type RefactoringDetector struct{}

func NewRefactoringDetector() *RefactoringDetector {
	return &RefactoringDetector{}
}

var refactorings = []struct {
	Name        string
	Regex       *regexp.Regexp
	Title       string
	Description string
	Type        string
	Severity    string
}{
	{
		Name:        "LongParameterList",
		Regex:       regexp.MustCompile(`func\s+\w+\s*\(([^)]{50,})\)`),
		Title:       "Long Parameter List",
		Description: "Consider using a struct to group parameters",
		Type:        "maintainability",
		Severity:    "suggestion",
	},
	{
		Name:        "DuplicateCode",
		Regex:       nil, // Needs cross-file analysis
		Title:       "Potential Duplicate Code",
		Description: "Similar code blocks detected - consider extracting to function",
		Type:        "duplication",
		Severity:    "warning",
	},
	{
		Name:        "NestedDepth",
		Regex:       regexp.MustCompile(`(if|for|switch|else)\s*\{[^}]{100,}\}`),
		Title:       "Deeply Nested Code",
		Description: "Consider extracting nested logic into separate functions",
		Type:        "complexity",
		Severity:    "suggestion",
	},
	{
		Name:        "GodObject",
		Regex:       regexp.MustCompile(`(class|type|struct)\s+(\w+){30,}`),
		Title:       "Potential God Object",
		Description: "Large class/struct - consider splitting into smaller units",
		Type:        "complexity",
		Severity:    "warning",
	},
	{
		Name:        "FeatureEnvy",
		Regex:       nil, // Needs semantic analysis
		Title:       "Feature Envy",
		Description: "Method seems more interested in another class - consider moving it",
		Type:        "design",
		Severity:    "info",
	},
	{
		Name:        "NullCheck",
		Regex:       regexp.MustCompile(`(if\s+.*==\s*null|if\s+.*!=\s*null|if\s+.*nil)`),
		Title:       "Explicit Null Check",
		Description: "Consider using Option/Result type for safer null handling",
		Type:        "null safety",
		Severity:    "suggestion",
	},
}

func (d *RefactoringDetector) AnalyzeFile(content string) []Refactoring {
	lines := strings.Split(content, "\n")
	var results []Refactoring

	for i, line := range lines {
		lineNum := i + 1

		for _, ref := range refactorings {
			if ref.Regex != nil && ref.Regex.MatchString(line) {
				results = append(results, Refactoring{
					Type:        ref.Type,
					Severity:    ref.Severity,
					Title:       ref.Title,
					Description: ref.Description,
					Line:        lineNum,
				})
			}
		}
	}

	// Additional analysis
	results = append(results, d.checkDeadCode(lines)...)
	results = append(results, d.checkComplexityMetrics(content, lines)...)

	return results
}

func (d *RefactoringDetector) checkDeadCode(lines []string) []Refactoring {
	var results []Refactoring

	// Simple dead code detection - unreachable returns
	for i := 0; i < len(lines)-1; i++ {
		current := strings.TrimSpace(lines[i])
		next := strings.TrimSpace(lines[i+1])

		// If we have a return followed by more code (not closing brace)
		if (current == "return" || current == "return nil" || current == "return false" || current == "return true") &&
			len(next) > 0 && !strings.HasPrefix(next, "//") && !strings.HasPrefix(next, "}") {
			results = append(results, Refactoring{
				Type:        "dead code",
				Severity:    "warning",
				Title:       "Unreachable Code",
				Description: "Code after return statement may be unreachable",
				Line:        i + 2,
			})
		}
	}

	return results
}

func (d *RefactoringDetector) checkComplexityMetrics(content string, lines []string) []Refactoring {
	var results []Refactoring

	// Count branches
	branchCount := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "if ") || strings.HasPrefix(trimmed, "else ") ||
			strings.HasPrefix(trimmed, "for ") || strings.HasPrefix(trimmed, "switch ") ||
			strings.HasPrefix(trimmed, "case ") || strings.HasPrefix(trimmed, "default:") {
			branchCount++
		}
	}

	// Cyclomatic complexity
	complexity := branchCount + 1
	if complexity > 20 {
		results = append(results, Refactoring{
			Type:        "complexity",
			Severity:    "warning",
			Title:       "High Cyclomatic Complexity",
			Description: "Function has high complexity (consider splitting)",
			Line:        1,
		})
	}

	return results
}

func (d *RefactoringDetector) GetSuggestions(refactorings []Refactoring) []string {
	var suggestions []string

	for _, r := range refactorings {
		suggestions = append(suggestions, fmt.Sprintf("[%s] %s: %s (line %d)", 
			r.Severity, r.Title, r.Description, r.Line))
	}

	return suggestions
}

func (d *RefactoringDetector) GenerateReport(refactorings []Refactoring) string {
	if len(refactorings) == 0 {
		return "✅ No refactoring opportunities detected. Code looks clean!"
	}

	output := "🔧 Refactoring Suggestions:\n\n"
	
	byType := make(map[string][]Refactoring)
	for _, r := range refactorings {
		byType[r.Type] = append(byType[r.Type], r)
	}
	
	bySeverity := make(map[string][]Refactoring)
	for _, r := range refactorings {
		bySeverity[r.Severity] = append(bySeverity[r.Severity], r)
	}
	
	// Show by severity first
	for _, severity := range []string{"warning", "suggestion", "info"} {
		items, ok := bySeverity[severity]
		if !ok || len(items) == 0 {
			continue
		}
		
		icon := "💡"
		if severity == "warning" {
			icon = "⚠️"
		}
		
		output += fmt.Sprintf("%s %s (%d)\n", icon, strings.ToUpper(severity), len(items))
		for _, r := range items {
			output += fmt.Sprintf("   • %s (line %d)\n", r.Title, r.Line)
			output += fmt.Sprintf("     %s\n", r.Description)
		}
		output += "\n"
	}

	return output
}