package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type RustParser struct{}

func NewRustParser() *RustParser {
	return &RustParser{}
}

var (
	rustImportRegex  = regexp.MustCompile(`^use\s+([a-zA-Z0-9_:]+)`)
	rustModRegex     = regexp.MustCompile(`^mod\s+([a-zA-Z0-9_]+)`)
	rustPubRegex     = regexp.MustCompile(`pub\s+(?:struct|enum|trait|fn|const|mod)\s+(\w+)`)
	rustStructRegex  = regexp.MustCompile(`^struct\s+(\w+)`)
	rustEnumRegex    = regexp.MustCompile(`^enum\s+(\w+)`)
	rustTraitRegex   = regexp.MustCompile(`^trait\s+(\w+)`)
	rustFnRegex      = regexp.MustCompile(`^fn\s+(\w+)`)
)

func (p *RustParser) ParseFile(path string) (*ParseResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	result := &ParseResult{
		File:     path,
		Language: "rust",
		Imports:  []string{},
		Exports:  []string{},
		Symbols:  []Symbol{},
	}

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip comments and attributes
		if strings.HasPrefix(trimmed, "//") || 
		   strings.HasPrefix(trimmed, "#[") ||
		   trimmed == "" {
			continue
		}

		// Use statements (imports)
		if matches := rustImportRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			imp := matches[1]
			// Remove trailing ::* or ::{...}
			if strings.HasSuffix(imp, "::*") {
				imp = strings.TrimSuffix(imp, "::*")
			}
			result.Imports = append(result.Imports, imp)
		}

		// Module declarations
		if matches := rustModRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			result.Imports = append(result.Imports, "mod:"+matches[1])
		}

		// Structs
		if matches := rustStructRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			result.Exports = append(result.Exports, matches[1])
			result.Symbols = append(result.Symbols, Symbol{
				Name: matches[1],
				Type: "struct",
				Location: Location{
					File: path,
					Line: i + 1,
				},
			})
		}

		// Enums
		if matches := rustEnumRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			result.Exports = append(result.Exports, matches[1])
			result.Symbols = append(result.Symbols, Symbol{
				Name: matches[1],
				Type: "enum",
				Location: Location{
					File: path,
					Line: i + 1,
				},
			})
		}

		// Traits
		if matches := rustTraitRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			result.Exports = append(result.Exports, matches[1])
			result.Symbols = append(result.Symbols, Symbol{
				Name: matches[1],
				Type: "trait",
				Location: Location{
					File: path,
					Line: i + 1,
				},
			})
		}

		// Functions
		if matches := rustFnRegex.FindStringSubmatch(trimmed); len(matches) > 1 {
			result.Exports = append(result.Exports, matches[1])
			result.Symbols = append(result.Symbols, Symbol{
				Name: matches[1],
				Type: "function",
				Location: Location{
					File: path,
					Line: i + 1,
				},
			})
		}

		// Impl blocks - just detect them
		_ = strings.HasPrefix(trimmed, "impl")
	}

	result.Dependencies = result.Imports
	return result, nil
}

func (p *RustParser) SupportedExtensions() []string {
	return []string{".rs"}
}

func FindRustProjectRoot(path string) (string, error) {
	// Check for Cargo.toml
	cargoPath := filepath.Join(path, "Cargo.toml")
	if _, err := os.Stat(cargoPath); err == nil {
		return path, nil
	}
	
	// Check parent directories
	parent := filepath.Dir(path)
	if parent != path {
		return FindRustProjectRoot(parent)
	}
	
	return "", fmt.Errorf("no Rust project root found (no Cargo.toml)")
}

// Alternative: Use Go's AST parser for more accurate Rust parsing
// This requires the rust AST parser library

func ParseRustWithAST(path string) (*ParseResult, error) {
	// This is a placeholder - full AST parsing would require
	// a proper Rust parser library like rust-ast or tree-sitter
	// For now, we use regex-based parsing above
	return nil, fmt.Errorf("AST parsing not implemented, using regex")
}