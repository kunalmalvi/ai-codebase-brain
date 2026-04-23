package parser

import "fmt"

type Parser interface {
	ParseFile(path string) (*ParseResult, error)
	SupportedExtensions() []string
}

type ParseResult struct {
	File         string
	Language     string
	Imports      []string
	Exports      []string
	Symbols      []Symbol
	Dependencies []string
}

type Symbol struct {
	Name     string
	Type     string
	Location Location
}

type Location struct {
	File     string
	Line     int
	Column   int
	EndLine  int
	EndColumn int
}

func DetectLanguage(filePath string) string {
	ext := ""
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '.' {
			ext = filePath[i:]
			break
		}
	}

	languageMap := map[string]string{
		".go":     "go",
		".js":     "js",
		".jsx":    "js",
		".ts":     "ts",
		".tsx":    "ts",
		".mjs":    "js",
		".cjs":    "js",
		".py":     "py",
		".rs":     "rs",
	}

	if lang, ok := languageMap[ext]; ok {
		return lang
	}
	return ""
}

var _ = fmt.Sprintf