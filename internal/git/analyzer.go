package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Hash    string
	Author  string
	Email   string
	Date    time.Time
	Message string
	Files  []string
}

type Branch struct {
	Name     string
	IsActive bool
}

type FileChange struct {
	File     string
	Type    string
	Additions int
	Deletions int
}

type RepoInfo struct {
	Root         string
	Branches     []Branch
	ActiveBranch string
	Remotes      []string
}

type GitAnalyzer struct {
	repoPath string
}

func NewGitAnalyzer(repoPath string) (*GitAnalyzer, error) {
	if _, err := runGitCommand(repoPath, "rev-parse", "--git-dir"); err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	return &GitAnalyzer{repoPath: repoPath}, nil
}

func (g *GitAnalyzer) GetRepoInfo() (*RepoInfo, error) {
	info := &RepoInfo{
		Root:    g.repoPath,
		Branches: []Branch{},
	}

	branchOutput, err := runGitCommand(g.repoPath, "branch", "-a")
	if err == nil {
		lines := strings.Split(branchOutput, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			isActive := strings.HasPrefix(line, "*")
			name := strings.TrimPrefix(line, "* ")
			info.Branches = append(info.Branches, Branch{
				Name:     name,
				IsActive: isActive,
			})
			if isActive {
				info.ActiveBranch = name
			}
		}
	}

	remoteOutput, err := runGitCommand(g.repoPath, "remote", "-v")
	if err == nil {
		lines := strings.Split(remoteOutput, "\n")
		seen := make(map[string]bool)
		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if !seen[parts[0]] {
					info.Remotes = append(info.Remotes, parts[0])
					seen[parts[0]] = true
				}
			}
		}
	}

	return info, nil
}

func (g *GitAnalyzer) GetCommits(limit int) ([]Commit, error) {
	output, err := runGitCommand(g.repoPath, "log", 
		"--format=%H|%an|%ae|%at|%s",
		"-n", strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}

	var commits []Commit
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 5 {
			timestamp, _ := strconv.ParseInt(parts[3], 10, 64)
			commits = append(commits, Commit{
				Hash:    parts[0],
				Author:  parts[1],
				Email:   parts[2],
				Date:    time.Unix(timestamp, 0),
				Message: parts[4],
			})
		}
	}

	return commits, nil
}

func (g *GitAnalyzer) GetFileHistory(filePath string, limit int) ([]Commit, error) {
	output, err := runGitCommand(g.repoPath, "log", 
		"--format=%H|%an|%ae|%at|%s",
		"-n", strconv.Itoa(limit),
		"--", filePath)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 5 {
			timestamp, _ := strconv.ParseInt(parts[3], 10, 64)
			commits = append(commits, Commit{
				Hash:    parts[0],
				Author:  parts[1],
				Email:   parts[2],
				Date:    time.Unix(timestamp, 0),
				Message: parts[4],
				Files:   []string{filePath},
			})
		}
	}

	return commits, nil
}

func (g *GitAnalyzer) GetChangedFiles(commitHash string) ([]FileChange, error) {
	output, err := runGitCommand(g.repoPath, "show", 
		"--stat", "--format=", commitHash)
	if err != nil {
		return nil, err
	}

	var changes []FileChange
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "files changed") {
			continue
		}

		if strings.Contains(line, "|") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				fileName := strings.TrimSpace(strings.Split(line, "|")[0])
				changeType := "modified"
				if strings.Contains(line, "insertion") || strings.Contains(line, "+") && !strings.Contains(line, "-") {
					changeType = "added"
				} else if strings.Contains(line, "deletion") || strings.Contains(line, "-") && !strings.Contains(line, "+") {
					changeType = "deleted"
				}

				changes = append(changes, FileChange{
					File: fileName,
					Type: changeType,
				})
			}
		}
	}

	return changes, nil
}

func (g *GitAnalyzer) GetDiff(from, to string) (string, error) {
	return runGitCommand(g.repoPath, "diff", from, to)
}

func (g *GitAnalyzer) GetCurrentChanges() (string, error) {
	return runGitCommand(g.repoPath, "status", "--short")
}

func (g *GitAnalyzer) GetStagedChanges() (string, error) {
	return runGitCommand(g.repoPath, "diff", "--cached", "--stat")
}

func (g *GitAnalyzer) GetBlame(filePath string) (string, error) {
	return runGitCommand(g.repoPath, "blame", "--line-porcelain", filePath)
}

func (g *GitAnalyzer) GetContributors() (map[string]int, error) {
	output, err := runGitCommand(g.repoPath, "shortlog", "-sne", "-n")
	if err != nil {
		return nil, err
	}

	contributors := make(map[string]int)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			count := 1
			for _, p := range parts {
				if n, err := strconv.Atoi(p); err == nil {
					count = n
					break
				}
			}

			nameStart := strings.Index(line, " ")
			for i := nameStart + 1; i < len(line); i++ {
				if line[i] == ' ' {
					nameStart = i
					break
				}
			}
			name := strings.TrimSpace(line[nameStart:])

			if idx := strings.Index(name, "<"); idx > 0 {
				name = strings.TrimSpace(name[:idx])
			}

			contributors[name] = count
		}
	}

	return contributors, nil
}

func (g *GitAnalyzer) SearchCommits(query string, limit int) ([]Commit, error) {
	output, err := runGitCommand(g.repoPath, "log", 
		"--all", "--grep="+query,
		"--format=%H|%an|%ae|%at|%s",
		"-n", strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}

	var commits []Commit
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 5 {
			timestamp, _ := strconv.ParseInt(parts[3], 10, 64)
			commits = append(commits, Commit{
				Hash:    parts[0],
				Author:  parts[1],
				Email:   parts[2],
				Date:    time.Unix(timestamp, 0),
				Message: parts[4],
			})
		}
	}

	return commits, nil
}

func (g *GitAnalyzer) GetRecentChanges(days int) ([]Commit, error) {
	since := fmt.Sprintf("--since=%d.days.ago", days)
	output, err := runGitCommand(g.repoPath, "log", 
		"--format=%H|%an|%ae|%at|%s",
		since)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 5 {
			timestamp, _ := strconv.ParseInt(parts[3], 10, 64)
			commits = append(commits, Commit{
				Hash:    parts[0],
				Author:  parts[1],
				Email:   parts[2],
				Date:    time.Unix(timestamp, 0),
				Message: parts[4],
			})
		}
	}

	return commits, nil
}

func runGitCommand(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("git command failed: %s", string(exitErr.Stderr))
		}
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}