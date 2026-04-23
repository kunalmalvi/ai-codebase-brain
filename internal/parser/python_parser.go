package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PythonParser struct{}

func NewPythonParser() *PythonParser {
	return &PythonParser{}
}

func (p *PythonParser) ParseFile(path string) (*ParseResult, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	result := &ParseResult{
		File:     path,
		Language: "python",
		Imports:  []string{},
		Exports:  []string{},
		Symbols:  []Symbol{},
	}

	lines := strings.Split(string(content), "\n")
	indentStack := []int{0}
	inClass := false
	inFunction := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Detect indentation
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		
		// Reset class/function when dedented
		if len(indentStack) > 0 && indent <= indentStack[0] {
			inClass = false
			inFunction = false
			if indent < indentStack[len(indentStack)-1] {
				indentStack = indentStack[:len(indentStack)-1]
			}
		}

		// Imports: import x, from x import y
		if strings.HasPrefix(trimmed, "import ") {
			module := strings.TrimPrefix(trimmed, "import ")
			module = strings.Split(module, " as ")[0]
			module = strings.Split(module, ",")[0]
			module = strings.TrimSpace(module)
			result.Imports = append(result.Imports, module)
		}
		
		if strings.HasPrefix(trimmed, "from ") {
			parts := strings.SplitN(trimmed, " import ", 2)
			if len(parts) == 2 {
				module := strings.TrimPrefix(parts[0], "from ")
				module = strings.TrimSpace(module)
				if module != "" {
					result.Imports = append(result.Imports, module)
				}
			}
		}

		// Class definitions
		if strings.HasPrefix(trimmed, "class ") {
			className := strings.TrimPrefix(trimmed, "class ")
			if idx := strings.Index(className, "("); idx > 0 {
				className = className[:idx]
			}
			if idx := strings.Index(className, ":"); idx > 0 {
				className = className[:idx]
			}
			className = strings.TrimSpace(className)
			
			result.Exports = append(result.Exports, className)
			result.Symbols = append(result.Symbols, Symbol{
				Name: className,
				Type: "class",
				Location: Location{
					File: path,
					Line: i + 1,
				},
			})
			inClass = true
			indentStack = []int{indent}
		}

		// Function/Method definitions
		if strings.HasPrefix(trimmed, "def ") {
			funcName := strings.TrimPrefix(trimmed, "def ")
			if idx := strings.Index(funcName, "("); idx > 0 {
				funcName = funcName[:idx]
			}
			funcName = strings.TrimSpace(funcName)
			
			// Only export top-level functions (not methods)
			if !inClass {
				result.Exports = append(result.Exports, funcName)
			}
			result.Symbols = append(result.Symbols, Symbol{
				Name: funcName,
				Type: "function",
				Location: Location{
					File: path,
					Line: i + 1,
				},
			})
			inFunction = true
			indentStack = append(indentStack, indent)
		}

		// Async functions
		if strings.HasPrefix(trimmed, "async def ") {
			funcName := strings.TrimPrefix(trimmed, "async def ")
			if idx := strings.Index(funcName, "("); idx > 0 {
				funcName = funcName[:idx]
			}
			funcName = strings.TrimSpace(funcName)
			
			if !inClass {
				result.Exports = append(result.Exports, funcName)
			}
			result.Symbols = append(result.Symbols, Symbol{
				Name: funcName,
				Type: "async_function",
				Location: Location{
					File: path,
					Line: i + 1,
				},
			})
			inFunction = true
			indentStack = append(indentStack, indent)
		}

		// Global variables (const assignments at module level)
		if !inClass && !inFunction {
			if strings.HasPrefix(trimmed, "_") && strings.Contains(trimmed, "=") {
				// Private global
			} else if (strings.HasPrefix(trimmed, "CONST") || 
			           strings.HasPrefix(trimmed, "URL") ||
			           strings.HasPrefix(trimmed, "DEFAULT") ||
			           strings.HasPrefix(trimmed, "MAX") ||
			           strings.HasPrefix(trimmed, "MIN")) && 
			          strings.Contains(trimmed, "=") {
				// Constants
				varName := strings.Split(trimmed, "=")[0]
				varName = strings.TrimSpace(varName)
				result.Exports = append(result.Exports, varName)
			}
		}
	}

	result.Dependencies = result.Imports
	return result, nil
}

func (p *PythonParser) SupportedExtensions() []string {
	return []string{".py"}
}

func FindPythonProjectRoot(path string) (string, error) {
	// Check for pyproject.toml, setup.py, requirements.txt
	for _, name := range []string{"pyproject.toml", "setup.py", "requirements.txt", "Pipfile"} {
		fullPath := filepath.Join(path, name)
		if _, err := os.Stat(fullPath); err == nil {
			return path, nil
		}
	}
	
	// Check parent directories
	parent := filepath.Dir(path)
	if parent != path {
		return FindPythonProjectRoot(parent)
	}
	
	return "", fmt.Errorf("no Python project root found")
}