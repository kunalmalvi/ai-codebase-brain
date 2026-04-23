package monorepo

import (
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Path        string   `json:"path"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Language    string   `json:"language"`
	Modules     []string `json:"modules"`
}

type MonorepoDetector struct{}

func NewMonorepoDetector() *MonorepoDetector {
	return &MonorepoDetector{}
}

func (d *MonorepoDetector) Detect(root string) ([]Project, error) {
	var projects []Project

	// Check for common monorepo indicators
	indicators := []string{
		"package.json",           // Node.js/npm workspaces
		"go.mod",                 // Go modules
		"Cargo.toml",             // Rust
		"pyproject.toml",         // Python
		"Pipfile",                // Python Pipenv
		"workspace.toml",         // Turborepo
		"lerna.json",             // Lerna
		"nx.json",                // Nx
		"rush.json",              // Rush
		"pnpm-workspace.yaml",    // pnpm workspaces
	}

	// First check if root itself is a project
	if rootProject := d.detectProject(root, indicators); rootProject != nil {
		projects = append(projects, *rootProject)
	}

	// Walk directories to find sub-projects
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and common non-project dirs
		if info.IsDir() {
			dirName := info.Name()
			if strings.HasPrefix(dirName, ".") ||
			   dirName == "node_modules" ||
			   dirName == "vendor" ||
			   dirName == "target" ||
			   dirName == "dist" ||
			   dirName == "build" ||
			   dirName == ".git" {
				if dirName != root {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check if any project indicator exists in this directory
		for _, indicator := range indicators {
			if info.Name() == indicator {
				projectDir := filepath.Dir(path)
				
				// Don't add if it's the root (already added)
				if projectDir != root && projectDir != filepath.Dir(root) {
					if proj := d.detectProject(projectDir, indicators); proj != nil {
						// Check if not already added
						exists := false
						for _, p := range projects {
							if p.Path == proj.Path {
								exists = true
								break
							}
						}
						if !exists {
							projects = append(projects, *proj)
						}
					}
				}
			}
		}

		return nil
	})

	return projects, err
}

func (d *MonorepoDetector) detectProject(path string, indicators []string) *Project {
	var projectType, language string
	var modules []string

	for _, indicator := range indicators {
		fullPath := filepath.Join(path, indicator)
		if _, err := os.Stat(fullPath); err == nil {
			switch indicator {
			case "package.json", "pnpm-workspace.yaml", "lerna.json", "nx.json", "rush.json", "workspace.toml":
				projectType = "javascript"
				language = "JavaScript/TypeScript"
				
				// Try to find workspaces
				if indicator == "package.json" {
					if pkgs := d.findWorkspaces(path); len(pkgs) > 0 {
						modules = pkgs
					}
				}
				
			case "go.mod":
				projectType = "go"
				language = "Go"
				
			case "Cargo.toml":
				projectType = "rust"
				language = "Rust"
				
			case "pyproject.toml", "Pipfile":
				projectType = "python"
				language = "Python"
			}
			break
		}
	}

	if projectType == "" {
		return nil
	}

	name := filepath.Base(path)
	
	return &Project{
		Path:     path,
		Name:     name,
		Type:     projectType,
		Language: language,
		Modules:  modules,
	}
}

func (d *MonorepoDetector) findWorkspaces(root string) []string {
	var workspaces []string

	patterns := []string{
		"packages/*",
		"apps/*",
		"src/*",
		"modules/*",
	}

	for _, pattern := range patterns {
		globPath := filepath.Join(root, pattern)
		matches, _ := filepath.Glob(globPath)
		for _, match := range matches {
			if info, err := os.Stat(match); err == nil && info.IsDir() {
				workspaces = append(workspaces, match)
			}
		}
	}

	return workspaces
}

func DetectMonorepoType(root string) string {
	indicators := map[string]string{
		"package.json":         "npm/yarn/pnpm workspaces",
		"pnpm-workspace.yaml":  "pnpm workspaces",
		"lerna.json":          "Lerna",
		"nx.json":             "Nx",
		"rush.json":           "Rush",
		"workspace.toml":      "Turborepo",
		"go.mod":              "Go modules",
		"Cargo.toml":          "Rust workspace",
		"pyproject.toml":      "Python (poetry/pip-tools)",
		"BUILD":               "Bazel",
		"WORKSPACE":           "Bazel",
	}

	for indicator, monorepoType := range indicators {
		path := filepath.Join(root, indicator)
		if _, err := os.Stat(path); err == nil {
			return monorepoType
		}
	}

	return ""
}

func (d *MonorepoDetector) DetectProject(path string, indicators []string) *Project {
	return d.detectProject(path, indicators)
}