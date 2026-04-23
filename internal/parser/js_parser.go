package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type JSParser struct{}

func NewJSParser() *JSParser {
	return &JSParser{}
}

var (
	importRegex    = regexp.MustCompile(`^import\s+(?:(?:\{[^}]*\}|\*|\w+)\s+from\s+)?['"]([^'"]+)['"]`)
	requireRegex   = regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	exportRegex    = regexp.MustCompile(`export\s+(?:default\s+)?(?:const|let|var|function|class)\s+(\w+)`)
	exportAllRegex = regexp.MustCompile(`export\s+\*\s+from\s+['"]([^'"]+)['"]`)
)

func (p *JSParser) ParseFile(path string) (*ParseResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	result := &ParseResult{
		File:     path,
		Language: detectJSType(path),
		Imports:  []string{},
		Exports:  []string{},
		Symbols:  []Symbol{},
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)

		if matches := importRegex.FindStringSubmatch(line); len(matches) > 1 {
			result.Imports = append(result.Imports, matches[1])
		}

		if matches := requireRegex.FindStringSubmatch(line); len(matches) > 1 {
			result.Imports = append(result.Imports, matches[1])
		}

		if matches := exportRegex.FindStringSubmatch(line); len(matches) > 1 {
			result.Exports = append(result.Exports, matches[1])
			result.Symbols = append(result.Symbols, Symbol{
				Name: matches[1],
				Type: "export",
				Location: Location{
					File: path,
					Line: i + 1,
				},
			})
		}

		if matches := exportAllRegex.FindStringSubmatch(line); len(matches) > 1 {
			result.Imports = append(result.Imports, matches[1])
		}
	}

	result.Dependencies = result.Imports
	return result, nil
}

func (p *JSParser) SupportedExtensions() []string {
	return []string{".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs"}
}

func detectJSType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".ts", ".tsx":
		return "ts"
	default:
		return "js"
	}
}

func FindPackageJSON(root string) (string, error) {
	path := filepath.Join(root, "package.json")
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return "", fmt.Errorf("no package.json found in %s", root)
}