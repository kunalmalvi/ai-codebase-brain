package smartsearch

import (
	"regexp"
	"strings"
)

type Pattern struct {
	Name        string
	Description string
	Regex       *regexp.Regexp
	Severity    string // "info", "warning", "suggestion"
	Category    string
}

var Patterns = []*Pattern{
	{
		Name:        "HardcodedSecret",
		Description: "Potential hardcoded secret or API key detected",
		Regex:       regexp.MustCompile(`(?i)(api[_-]?key|secret|password|token|auth)[^\n]{0,30}['""]?\s*[=:]\s*['""][a-zA-Z0-9_\-]{20,}`),
		Severity:    "warning",
		Category:    "security",
	},
	{
		Name:        "TODO",
		Description: "TODO comment found - work in progress",
		Regex:       regexp.MustCompile(`(?i)\b(todo|fixme|hack|xxx)\b`),
		Severity:    "info",
		Category:    "quality",
	},
	{
		Name:        "ConsoleLog",
		Description: "Console log statement found - remove in production",
		Regex:       regexp.MustCompile(`(?i)\b(console\.(log|debug|info|warn|error)|print|println|fmt\.Print)\b`),
		Severity:    "suggestion",
		Category:    "quality",
	},
	{
		Name:        "EmptyCatch",
		Description: "Empty catch block - errors are being swallowed",
		Regex:       regexp.MustCompile(`(?i)(catch\s*\([^)]*\)\s*\{\s*\}|except\s*:\s*pass)`),
		Severity:    "warning",
		Category:    "quality",
	},
	{
		Name:        "LongFunction",
		Description: "Function is very long (>100 lines) - consider splitting",
		Regex:       nil, // Handled specially
		Severity:    "suggestion",
		Category:    "maintainability",
	},
	{
		Name:        "NestedCallbacks",
		Description: "Deeply nested callbacks - consider refactoring",
		Regex:       nil, // Handled specially
		Severity:    "suggestion",
		Category:    "maintainability",
	},
	{
		Name:        "SQLInjection",
		Description: "Potential SQL injection vulnerability",
		Regex:       regexp.MustCompile(`(?i)(execute|query|exec)\s*\(\s*['"]\s*\+.*(request|params|post|get)`),
		Severity:    "warning",
		Category:    "security",
	},
	{
		Name:        "ErrorNotReturned",
		Description: "Error is not being returned or handled",
		Regex:       regexp.MustCompile(`(?i)(err\s*:=|err\s*=|error\s*:=|Exception\s+)\s*(?!nil|error|nil)`),
		Severity:    "warning",
		Category:    "quality",
	},
	{
		Name:        "UnclosedResource",
		Description: "Resource may not be properly closed",
		Regex:       regexp.MustCompile(`(?i)(Open\(|NewReader\(|Create\(|Get\()[^)]*(?!defer|defer\s+.*Close)`),
		Severity:    "suggestion",
		Category:    "quality",
	},
	{
		Name:        "MagicNumber",
		Description: "Magic number found - use a named constant",
		Regex:       regexp.MustCompile(`\b([0-9]{3,}|0x[0-9A-Fa-f]{3,})\b`),
		Severity:    "suggestion",
		Category:    "style",
	},
}

type PatternDetector struct {
	patterns []*Pattern
}

func NewPatternDetector() *PatternDetector {
	return &PatternDetector{
		patterns: Patterns,
	}
}

type PatternMatch struct {
	Pattern   string
	Line      int
	Content   string
	Severity  string
	Category  string
}

func (d *PatternDetector) AnalyzeFile(content string) []PatternMatch {
	lines := strings.Split(content, "\n")
	var matches []PatternMatch
	
	for i, line := range lines {
		lineNum := i + 1
		
		// Check for line-based patterns
		for _, pattern := range d.patterns {
			if pattern.Regex == nil {
				continue // Skip patterns that need special handling
			}
			
			if pattern.Regex.MatchString(line) {
				matches = append(matches, PatternMatch{
					Pattern:  pattern.Name,
					Line:     lineNum,
					Content:  strings.TrimSpace(line),
					Severity: pattern.Severity,
					Category: pattern.Category,
				})
			}
		}
		
		// Special checks
		if len(line) > 200 && !strings.HasPrefix(strings.TrimSpace(line), "//") && 
		   !strings.HasPrefix(strings.TrimSpace(line), "#") {
			matches = append(matches, PatternMatch{
				Pattern:  "LongLine",
				Line:     lineNum,
				Content:  "Line exceeds 200 characters",
				Severity: "suggestion",
				Category: "style",
			})
		}
	}
	
	// Check for long functions (special)
	functionLengths := d.countFunctionLengths(content)
	for line, length := range functionLengths {
		if length > 100 {
			matches = append(matches, PatternMatch{
				Pattern:  "LongFunction",
				Line:     line,
				Content:  "Function has ~100+ lines, consider splitting",
				Severity: "suggestion",
				Category: "maintainability",
			})
		}
	}
	
	return matches
}

func (d *PatternDetector) countFunctionLengths(content string) map[int]int {
	lengths := make(map[int]int)
	lines := strings.Split(content, "\n")
	
	inFunc := false
	funcStart := 0
	
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Detect function start
		isFuncLine := false
		funcIndicators := []string{"func ", "def ", "function ", "fn ", "public ", "private ", "async "}
		for _, ind := range funcIndicators {
			if strings.HasPrefix(trimmed, ind) || strings.HasPrefix(trimmed, "func (") {
				isFuncLine = true
				break
			}
		}
		
		if isFuncLine && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "#") {
			if inFunc && funcStart > 0 {
				lengths[funcStart] = i - funcStart
			}
			inFunc = true
			funcStart = i + 1
		}
		
		// Track closing braces/blocks
		if inFunc && (trimmed == "}" || trimmed == "}" || strings.HasPrefix(trimmed, "return") || trimmed == "end") {
			if funcStart > 0 {
				lengths[funcStart] = i - funcStart + 1
			}
			inFunc = false
		}
	}
	
	return lengths
}

func (d *PatternDetector) GetSecurityIssues(matches []PatternMatch) []PatternMatch {
	var security []PatternMatch
	for _, m := range matches {
		if m.Category == "security" {
			security = append(security, m)
		}
	}
	return security
}

func (d *PatternDetector) GetQualityIssues(matches []PatternMatch) []PatternMatch {
	var quality []PatternMatch
	for _, m := range matches {
		if m.Category == "quality" {
			quality = append(quality, m)
		}
	}
	return quality
}