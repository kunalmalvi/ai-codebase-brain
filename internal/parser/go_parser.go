package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type GoParser struct{}

func NewGoParser() *GoParser {
	return &GoParser{}
}

func (p *GoParser) ParseFile(path string) (*ParseResult, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	result := &ParseResult{
		File:     path,
		Language: "go",
		Imports:  []string{},
		Exports:  []string{},
		Symbols:  []Symbol{},
	}

	for _, imp := range node.Imports {
		if imp.Path != nil {
			importPath := strings.Trim(imp.Path.Value, `"`)
			result.Imports = append(result.Imports, importPath)
		}
	}

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name != nil {
				sym := Symbol{
					Name: d.Name.Name,
					Type: "function",
					Location: Location{
						File: path,
						Line: fset.Position(d.Pos()).Line,
					},
				}
				result.Symbols = append(result.Symbols, sym)

				if d.Recv == nil && d.Name.IsExported() {
					result.Exports = append(result.Exports, d.Name.Name)
				}
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if s.Name != nil {
						result.Symbols = append(result.Symbols, Symbol{
							Name: s.Name.Name,
							Type: "type",
							Location: Location{
								File: path,
								Line: fset.Position(s.Pos()).Line,
							},
						})
						if s.Name.IsExported() {
							result.Exports = append(result.Exports, s.Name.Name)
						}
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						if name.IsExported() {
							result.Exports = append(result.Exports, name.Name)
							result.Symbols = append(result.Symbols, Symbol{
								Name: name.Name,
								Type: "variable",
								Location: Location{
									File: path,
									Line: fset.Position(name.Pos()).Line,
								},
							})
						}
					}
				}
			}
		}
	}

	result.Dependencies = result.Imports

	return result, nil
}

func (p *GoParser) SupportedExtensions() []string {
	return []string{".go"}
}

func FindGoModules(root string) ([]string, error) {
	var modules []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}
		if info.Name() == "go.mod" {
			modules = append(modules, filepath.Dir(path))
		}
		return nil
	})

	return modules, err
}