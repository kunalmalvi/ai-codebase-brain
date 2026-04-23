package mcp

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/codebase-brain/internal/config"
	"github.com/codebase-brain/internal/errors"
	"github.com/codebase-brain/internal/git"
	"github.com/codebase-brain/internal/indexer"
	"github.com/codebase-brain/internal/logger"
	"github.com/codebase-brain/internal/monorepo"
	"github.com/codebase-brain/internal/session"
	"github.com/codebase-brain/internal/smartsearch"
	"github.com/codebase-brain/internal/token"
	"github.com/mark3labs/mcp-go/mcp"
)

var idx *indexer.Indexer
var sessMgr *session.SessionManager
var tokenOptimizer *token.TokenOptimizer
var codeCompressor *token.CodeCompressor
var monorepoDetector *monorepo.MonorepoDetector
var gitAnalyzer *git.GitAnalyzer
var cfg *config.Config

func init() {
	cfg = config.Default()
	var err error
	idx, err = indexer.NewIndexer(cfg)
	if err != nil {
		logger.Warn("failed to create indexer: %v", err)
	}
	
	sessMgr, err = session.NewSessionManager(cfg.ProjectPath)
	if err != nil {
		logger.Warn("failed to create session manager: %v", err)
	}
	
	tokenOptimizer = token.NewTokenOptimizer(true)
	codeCompressor = token.NewCodeCompressor()
	monorepoDetector = monorepo.NewMonorepoDetector()
	
	gitAnalyzer, err = git.NewGitAnalyzer(cfg.ProjectPath)
	if err != nil {
		logger.Warn("git analyzer not initialized: %v", err)
	}

	semanticSearch = smartsearch.NewSemanticSearch()
	codeSummarizer = smartsearch.NewSummarizer()
	patternDetector = smartsearch.NewPatternDetector()
	
	logger.Info("MCP handlers initialized")
}

func validateRequired(args map[string]interface{}, fields ...string) error {
	for _, field := range fields {
		val, exists := args[field]
		if !exists {
			return errors.InvalidInput(field, "is required")
		}
		str, isString := val.(string)
		if !isString || strings.TrimSpace(str) == "" {
			return errors.InvalidInput(field, "is required")
		}
	}
	return nil
}

func getStringArg(args map[string]interface{}, key string, defaultVal string) string {
	if val, ok := args[key].(string); ok && val != "" {
		return val
	}
	return defaultVal
}

func getIntArg(args map[string]interface{}, key string, defaultVal int) int {
	if val, ok := args[key].(float64); ok {
		return int(val)
	}
	return defaultVal
}

func handleSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "query"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	query := getStringArg(args, "query", "")
	scope := getStringArg(args, "scope", ".")
	limit := getIntArg(args, "limit", 10)

	logger.Debug("search request: query=%s, scope=%s, limit=%d", query, scope, limit)

	if idx == nil {
		return mcp.NewToolResultText(fmt.Sprintf("Search results for '%s' (limited to %d):\n\nNo indexed project. Run index_project first.", query, limit)), nil
	}

	results := idx.Search(query, limit)
	if len(results) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No results found for: %s", query)), nil
	}

	output := fmt.Sprintf("Found %d results for '%s':\n\n", len(results), query)
	for i, r := range results {
		output += fmt.Sprintf("%d. %s\n   %s\n", i+1, r.FilePath, r.Match)
	}
	return mcp.NewToolResultText(output), nil
}

func handleGetRelated(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "file_path"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	filePath := getStringArg(args, "file_path", "")
	depth := getIntArg(args, "depth", 2)

	logger.Debug("get_related: file=%s, depth=%d", filePath, depth)

	if idx == nil {
		return mcp.NewToolResultText("No indexed project. Run index_project first."), nil
	}

	related := idx.GetRelated(filePath, depth)
	if len(related) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No related files found for: %s", filePath)), nil
	}

	output := fmt.Sprintf("Related files to %s (depth %d):\n\n", filePath, depth)
	for i, r := range related {
		output += fmt.Sprintf("%d. %s (type: %s)\n", i+1, r.FilePath, r.Type)
	}
	return mcp.NewToolResultText(output), nil
}

func handleGetProjectGraph(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	limit := getIntArg(args, "limit", 100)

	if idx == nil {
		return mcp.NewToolResultText("No indexed project. Run index_project first."), nil
	}

	nodes := idx.GetGraphOverview(limit)
	if len(nodes) == 0 {
		return mcp.NewToolResultText("Project graph is empty. Run index_project first."), nil
	}

	output := fmt.Sprintf("Project graph (%d nodes):\n\n", len(nodes))
	for _, n := range nodes {
		output += fmt.Sprintf("- %s (%s)\n", n.Name, n.Type)
	}
	return mcp.NewToolResultText(output), nil
}

func handleIndexProject(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	path := getStringArg(args, "path", ".")

	logger.Info("indexing project at: %s", path)

	if idx == nil {
		return mcp.NewToolResultError("Indexer not initialized"), nil
	}

	count, err := idx.Index(path)
	if err != nil {
		logger.Error("indexing failed: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Indexing failed: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully indexed %d files", count)), nil
}

func handleGetFileContext(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "file_path"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	filePath := getStringArg(args, "file_path", "")
	contextType := getStringArg(args, "context_type", "all")

	if idx == nil {
		return mcp.NewToolResultText("No indexed project. Run index_project first."), nil
	}

	result := idx.GetFileContext(filePath, contextType)
	return mcp.NewToolResultText(result), nil
}

func handleCreateSession(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	projectID := getStringArg(args, "project_id", "default")

	if sessMgr == nil {
		return mcp.NewToolResultError("Session manager not initialized"), nil
	}

	s := sessMgr.CreateSession(projectID)
	return mcp.NewToolResultText(fmt.Sprintf("Created session: %s\nProject: %s\nCreated at: %s", s.ID, s.ProjectID, s.CreatedAt)), nil
}

func handleGetSession(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "session_id"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	sessionID := getStringArg(args, "session_id", "")

	if sessMgr == nil {
		return mcp.NewToolResultError("Session manager not initialized"), nil
	}

	s, ok := sessMgr.GetSession(sessionID)
	if !ok {
		return mcp.NewToolResultText(fmt.Sprintf("Session not found: %s", sessionID)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Session: %s\nProject: %s\nMessages: %d\nCreated: %s\nUpdated: %s", 
		s.ID, s.ProjectID, len(s.Messages), s.CreatedAt, s.UpdatedAt)), nil
}

func handleListSessions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if sessMgr == nil {
		return mcp.NewToolResultError("Session manager not initialized"), nil
	}

	sessions := sessMgr.ListSessions()
	if len(sessions) == 0 {
		return mcp.NewToolResultText("No sessions found"), nil
	}

	output := fmt.Sprintf("Found %d sessions:\n\n", len(sessions))
	for _, s := range sessions {
		output += fmt.Sprintf("- %s (project: %s, messages: %d, updated: %s)\n", 
			s.ID, s.ProjectID, len(s.Messages), s.UpdatedAt)
	}
	return mcp.NewToolResultText(output), nil
}

func handleAddMessage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "session_id", "role", "content"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	sessionID := getStringArg(args, "session_id", "")
	role := getStringArg(args, "role", "")
	content := getStringArg(args, "content", "")

	if role != "user" && role != "assistant" {
		return mcp.NewToolResultError("role must be 'user' or 'assistant'"), nil
	}

	if sessMgr == nil {
		return mcp.NewToolResultError("Session manager not initialized"), nil
	}

	err := sessMgr.AddMessage(sessionID, role, content)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to add message: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Added %s message to session %s", role, sessionID)), nil
}

func handleGetConversationHistory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "session_id"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	sessionID := getStringArg(args, "session_id", "")

	if sessMgr == nil {
		return mcp.NewToolResultError("Session manager not initialized"), nil
	}

	messages := sessMgr.GetMessages(sessionID)
	if len(messages) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No messages in session %s", sessionID)), nil
	}

	output := fmt.Sprintf("Conversation history for session %s:\n\n", sessionID)
	for _, m := range messages {
		output += fmt.Sprintf("[%s] %s: %s\n", m.Time.Format("2006-01-02 15:04"), m.Role, m.Content)
	}
	return mcp.NewToolResultText(output), nil
}

func handleEstimateTokens(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "text"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	text := getStringArg(args, "text", "")

	if codeCompressor == nil {
		codeCompressor = token.NewCodeCompressor()
	}

	count := codeCompressor.EstimateTokens(text)
	return mcp.NewToolResultText(fmt.Sprintf("Estimated tokens: %d (characters: %d)", count, len(text))), nil
}

func handleOptimizeContext(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "text"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	text := getStringArg(args, "text", "")
	maxTokens := getIntArg(args, "max_tokens", 32000)

	if codeCompressor == nil {
		codeCompressor = token.NewCodeCompressor()
	}

	compressed := codeCompressor.Compress(text, "medium")
	
	if tokenOptimizer != nil {
		deduped, isDup := tokenOptimizer.Deduplicate(text)
		if isDup {
			return mcp.NewToolResultText("Content is duplicate, skipping"), nil
		}
		_ = deduped
	}
	_ = maxTokens

	origTokens := codeCompressor.EstimateTokens(text)
	optTokens := codeCompressor.EstimateTokens(compressed)
	
	savings := float64(0)
	if origTokens > 0 {
		savings = float64(origTokens-optTokens) / float64(origTokens) * 100
	}
	
	output := fmt.Sprintf("Original: ~%d tokens\nOptimized: ~%d tokens\nSavings: %.1f%%\n\n%s",
		origTokens, optTokens, savings, truncate(compressed, 1000))
	
	return mcp.NewToolResultText(output), nil
}

func handleGetTokenStats(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if idx == nil {
		return mcp.NewToolResultText("No indexed project"), nil
	}

	nodes := idx.GetGraphOverview(10000)
	totalTokens := 0
	
	for _, node := range nodes {
		if content, ok := node.Metadata["content"].(string); ok {
			if codeCompressor == nil {
				codeCompressor = token.NewCodeCompressor()
			}
			totalTokens += codeCompressor.EstimateTokens(content)
		}
	}

	return mcp.NewToolResultText(fmt.Sprintf("Token Statistics:\n- Indexed files: %d\n- Estimated total tokens: ~%d\n- Deduplication: %s",
		len(nodes), totalTokens, "enabled")), nil
}

func handleDetectMonorepo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	path := getStringArg(args, "path", ".")

	if monorepoDetector == nil {
		monorepoDetector = monorepo.NewMonorepoDetector()
	}

	projects, err := monorepoDetector.Detect(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to detect monorepo: %v", err)), nil
	}

	if len(projects) == 0 {
		return mcp.NewToolResultText("No monorepo detected - this appears to be a single project"), nil
	}

	if len(projects) == 1 {
		return mcp.NewToolResultText(fmt.Sprintf("Single project detected:\n\n- Name: %s\n- Type: %s\n- Language: %s\n- Path: %s",
			projects[0].Name, projects[0].Type, projects[0].Language, projects[0].Path)), nil
	}

	output := fmt.Sprintf("Monorepo detected with %d projects:\n\n", len(projects))
	for i, p := range projects {
		output += fmt.Sprintf("%d. %s (%s)\n   Path: %s\n", i+1, p.Name, p.Language, p.Path)
		if len(p.Modules) > 0 {
			output += fmt.Sprintf("   Modules: %d\n", len(p.Modules))
		}
	}
	
	output += fmt.Sprintf("\nMonorepo type: %s", monorepo.DetectMonorepoType(path))
	
	return mcp.NewToolResultText(output), nil
}

func handleGetProjectInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "path"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	path := getStringArg(args, "path", "")

	if monorepoDetector == nil {
		monorepoDetector = monorepo.NewMonorepoDetector()
	}

	indicators := []string{"package.json", "go.mod", "Cargo.toml", "pyproject.toml"}
	proj := monorepoDetector.DetectProject(path, indicators)
	
	if proj == nil {
		return mcp.NewToolResultText(fmt.Sprintf("No project found at: %s", path)), nil
	}

	output := fmt.Sprintf("Project: %s\nType: %s\nLanguage: %s\nPath: %s\n",
		proj.Name, proj.Type, proj.Language, proj.Path)
	
	if len(proj.Modules) > 0 {
		output += "\nModules:\n"
		for _, m := range proj.Modules {
			output += fmt.Sprintf("  - %s\n", m)
		}
	}

	return mcp.NewToolResultText(output), nil
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// Smart search handlers

var semanticSearch *smartsearch.SemanticSearch
var codeSummarizer *smartsearch.Summarizer
var patternDetector *smartsearch.PatternDetector

func init() {
	semanticSearch = smartsearch.NewSemanticSearch()
	codeSummarizer = smartsearch.NewSummarizer()
	patternDetector = smartsearch.NewPatternDetector()
}

func handleAnalyzeCode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "file_path"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	filePath := getStringArg(args, "file_path", "")
	includePatterns := getBoolArg(args, "include_patterns", true)
	includeSummary := getBoolArg(args, "include_summary", true)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
	}

	output := fmt.Sprintf("📊 Analysis for: %s\n\n", filePath)

	if includeSummary {
		summary := codeSummarizer.SummarizeFile(filePath, string(content), nil, nil)
		output += codeSummarizer.GenerateQuickSummary(summary)
		output += "\n"
	}

	if includePatterns {
		matches := patternDetector.AnalyzeFile(string(content))
		if len(matches) > 0 {
			output += "🔍 Pattern Detection:\n\n"
			
			security := patternDetector.GetSecurityIssues(matches)
			quality := patternDetector.GetQualityIssues(matches)
			
			if len(security) > 0 {
				output += "⚠️ Security Issues:\n"
				for _, m := range security {
					output += fmt.Sprintf("   [%s] Line %d: %s\n", m.Severity, m.Line, m.Pattern)
				}
				output += "\n"
			}
			
			if len(quality) > 0 {
				output += "📝 Quality Issues:\n"
				for _, m := range quality {
					output += fmt.Sprintf("   [%s] Line %d: %s\n", m.Severity, m.Line, m.Pattern)
				}
			}
		} else {
			output += "✅ No obvious patterns detected\n"
		}
	}

	return mcp.NewToolResultText(output), nil
}

func handleSemanticSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "query"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	query := getStringArg(args, "query", "")
	filePath := getStringArg(args, "file_path", "")

	// First try direct concept search
	concepts := semanticSearch.SearchByConcept(query)
	
	if len(concepts) > 0 {
		output := fmt.Sprintf("🔍 Semantic search for: '%s'\n\n", query)
		output += "Found related concepts:\n"
		for _, c := range concepts {
			category := semanticSearch.GetCategory(c)
			output += fmt.Sprintf("  • %s (category: %s)\n", c, category)
			
			related := semanticSearch.GetRelatedConcepts(c)
			if len(related) > 0 {
				output += fmt.Sprintf("    Related: %s\n", strings.Join(related, ", "))
			}
		}
		return mcp.NewToolResultText(output), nil
	}

	// If no concept match, search by file content
	if filePath != "" && idx != nil {
		results := idx.Search(query, 10)
		if len(results) > 0 {
			output := fmt.Sprintf("🔍 Search results for: '%s'\n\n", query)
			for i, r := range results {
				output += fmt.Sprintf("%d. %s\n", i+1, r.FilePath)
			}
			return mcp.NewToolResultText(output), nil
		}
	}

	return mcp.NewToolResultText(fmt.Sprintf("No results found for: %s", query)), nil
}

func handleGetCodeSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "file_path"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	filePath := getStringArg(args, "file_path", "")

	content, err := os.ReadFile(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
	}

	summary := codeSummarizer.SummarizeFile(filePath, string(content), nil, nil)
	output := codeSummarizer.GenerateQuickSummary(summary)

	return mcp.NewToolResultText(output), nil
}

func handleDetectPatterns(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	
	if err := validateRequired(args, "file_path"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	
	filePath := getStringArg(args, "file_path", "")
	category := getStringArg(args, "category", "all")

	content, err := os.ReadFile(filePath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to read file: %v", err)), nil
	}

	matches := patternDetector.AnalyzeFile(string(content))
	
	if category != "all" {
		var filtered []smartsearch.PatternMatch
		for _, m := range matches {
			if m.Category == category {
				filtered = append(filtered, m)
			}
		}
		matches = filtered
	}

	if len(matches) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No patterns detected in %s", filePath)), nil
	}

	output := fmt.Sprintf("📋 Pattern Detection Results for: %s\n\n", filePath)
	
	// Group by category
	byCategory := make(map[string][]smartsearch.PatternMatch)
	for _, m := range matches {
		byCategory[m.Category] = append(byCategory[m.Category], m)
	}
	
	for cat, items := range byCategory {
		output += fmt.Sprintf("=== %s (%d) ===\n", cat, len(items))
		for _, m := range items {
			output += fmt.Sprintf("  [%s] Line %d: %s\n", m.Severity, m.Line, m.Content)
		}
		output += "\n"
	}

return mcp.NewToolResultText(output), nil
}

func getBoolArg(args map[string]interface{}, key string, defaultVal bool) bool {
	if val, ok := args[key].(bool); ok {
		return val
	}
	return defaultVal
}

// Git handlers

func handleGetGitInfo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if gitAnalyzer == nil {
		return mcp.NewToolResultText("Git analyzer not initialized. Project may not be a git repository."), nil
	}

	info, err := gitAnalyzer.GetRepoInfo()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get git info: %v", err)), nil
	}

	output := fmt.Sprintf("Git Repository: %s\n\n", info.Root)
	output += fmt.Sprintf("Active branch: %s\n\n", info.ActiveBranch)

	if len(info.Branches) > 0 {
		output += "Branches:\n"
		for _, b := range info.Branches {
			marker := " "
			if b.IsActive {
				marker = "*"
			}
			output += fmt.Sprintf("  %s %s\n", marker, b.Name)
		}
		output += "\n"
	}

	if len(info.Remotes) > 0 {
		output += fmt.Sprintf("Remotes: %s\n", strings.Join(info.Remotes, ", "))
	}

	return mcp.NewToolResultText(output), nil
}

func handleGetCommits(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	limit := getIntArg(args, "limit", 10)

	if gitAnalyzer == nil {
		return mcp.NewToolResultText("Git analyzer not initialized. Project may not be a git repository."), nil
	}

	commits, err := gitAnalyzer.GetCommits(limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get commits: %v", err)), nil
	}

	if len(commits) == 0 {
		return mcp.NewToolResultText("No commits found"), nil
	}

	output := fmt.Sprintf("Recent %d commits:\n\n", len(commits))
	for _, c := range commits {
		output += fmt.Sprintf("%s | %s | %s\n  %s\n\n",
			c.Hash[:7], c.Author, c.Date.Format("2006-01-02"), c.Message)
	}

	return mcp.NewToolResultText(output), nil
}

func handleGetFileHistory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	if err := validateRequired(args, "file_path"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	filePath := getStringArg(args, "file_path", "")
	limit := getIntArg(args, "limit", 10)

	if gitAnalyzer == nil {
		return mcp.NewToolResultText("Git analyzer not initialized. Project may not be a git repository."), nil
	}

	commits, err := gitAnalyzer.GetFileHistory(filePath, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get file history: %v", err)), nil
	}

	if len(commits) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No history found for: %s", filePath)), nil
	}

	output := fmt.Sprintf("History for %s (%d commits):\n\n", filePath, len(commits))
	for _, c := range commits {
		output += fmt.Sprintf("%s | %s | %s\n  %s\n\n",
			c.Hash[:7], c.Author, c.Date.Format("2006-01-02"), c.Message)
	}

	return mcp.NewToolResultText(output), nil
}

func handleSearchCommits(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	if err := validateRequired(args, "query"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	query := getStringArg(args, "query", "")
	limit := getIntArg(args, "limit", 10)

	if gitAnalyzer == nil {
		return mcp.NewToolResultText("Git analyzer not initialized. Project may not be a git repository."), nil
	}

	commits, err := gitAnalyzer.SearchCommits(query, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search commits: %v", err)), nil
	}

	if len(commits) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No commits matching: %s", query)), nil
	}

	output := fmt.Sprintf("Commits matching '%s' (%d):\n\n", query, len(commits))
	for _, c := range commits {
		output += fmt.Sprintf("%s | %s | %s\n  %s\n\n",
			c.Hash[:7], c.Author, c.Date.Format("2006-01-02"), c.Message)
	}

	return mcp.NewToolResultText(output), nil
}

func handleGetContributors(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if gitAnalyzer == nil {
		return mcp.NewToolResultText("Git analyzer not initialized. Project may not be a git repository."), nil
	}

	contributors, err := gitAnalyzer.GetContributors()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get contributors: %v", err)), nil
	}

	if len(contributors) == 0 {
		return mcp.NewToolResultText("No contributors found"), nil
	}

	output := fmt.Sprintf("Contributors (%d):\n\n", len(contributors))
	for name, count := range contributors {
		output += fmt.Sprintf("  %s: %d commits\n", name, count)
	}

	return mcp.NewToolResultText(output), nil
}

func handleGetDiff(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	from := getStringArg(args, "from", "")
	to := getStringArg(args, "to", "")

	if gitAnalyzer == nil {
		return mcp.NewToolResultText("Git analyzer not initialized. Project may not be a git repository."), nil
	}

	diff, err := gitAnalyzer.GetDiff(from, to)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get diff: %v", err)), nil
	}

	if diff == "" {
		return mcp.NewToolResultText("No diff available"), nil
	}

	return mcp.NewToolResultText(diff), nil
}

func handleGetChangedFiles(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	if err := validateRequired(args, "commit_hash"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	commitHash := getStringArg(args, "commit_hash", "")

	if gitAnalyzer == nil {
		return mcp.NewToolResultText("Git analyzer not initialized. Project may not be a git repository."), nil
	}

	changes, err := gitAnalyzer.GetChangedFiles(commitHash)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get changed files: %v", err)), nil
	}

	if len(changes) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No files changed in commit: %s", commitHash)), nil
	}

	output := fmt.Sprintf("Files changed in %s:\n\n", commitHash[:7])
	for _, c := range changes {
		output += fmt.Sprintf("  %s: %s\n", c.Type, c.File)
	}

	return mcp.NewToolResultText(output), nil
}

// Cross-project session history handlers

func handleListProjects(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	limit := getIntArg(args, "limit", 10)

	if projectIndex == nil {
		return mcp.NewToolResultText("Project index not initialized"), nil
	}

	projects := projectIndex.GetRecentProjects(limit)
	if len(projects) == 0 {
		return mcp.NewToolResultText("No projects found. Work on a project to add it to the index."), nil
	}

	output := fmt.Sprintf("Recent projects (%d):\n\n", len(projects))
	for _, proj := range projects {
		output += fmt.Sprintf("- %s\n", proj.Name)
		output += fmt.Sprintf("  Path: %s\n", proj.Path)
		output += fmt.Sprintf("  Language: %s\n", proj.Language)
		output += fmt.Sprintf("  Last session: %s\n", proj.LastSession.Format("2006-01-02 15:04"))
		output += fmt.Sprintf("  Sessions: %d\n", proj.SessionCount)
		if proj.LastWork != "" {
			output += fmt.Sprintf("  Last work: %s\n", proj.LastWork)
		}
		output += "\n"
	}

	return mcp.NewToolResultText(output), nil
}

func handleSearchProjects(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	if err := validateRequired(args, "query"); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	query := getStringArg(args, "query", "")

	if projectIndex == nil {
		return mcp.NewToolResultText("Project index not initialized"), nil
	}

	projects := projectIndex.SearchProjects(query)
	if len(projects) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No projects matching: %s", query)), nil
	}

	output := fmt.Sprintf("Projects matching '%s' (%d):\n\n", query, len(projects))
	for _, proj := range projects {
		output += fmt.Sprintf("- %s\n", proj.Name)
		output += fmt.Sprintf("  Path: %s\n", proj.Path)
		output += fmt.Sprintf("  Language: %s\n", proj.Language)
		if proj.LastWork != "" {
			output += fmt.Sprintf("  Last work: %s\n", proj.LastWork)
		}
		output += "\n"
	}

	return mcp.NewToolResultText(output), nil
}

func handleRegisterProject(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments

	projectPath := getStringArg(args, "path", cfg.ProjectPath)
	language := getStringArg(args, "language", "")
	lastWork := getStringArg(args, "last_work", "")

	if projectIndex == nil {
		return mcp.NewToolResultError("Project index not initialized"), nil
	}

	if language == "" {
		language = detectLanguage(projectPath)
	}

	projectIndex.RegisterProject(projectPath, language, lastWork)

	return mcp.NewToolResultText(fmt.Sprintf("Registered project: %s\nLanguage: %s", projectPath, language)), nil
}

func handleGetProjectActivity(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.Params.Arguments
	projectPath := getStringArg(args, "path", cfg.ProjectPath)

	if projectIndex == nil {
		return mcp.NewToolResultText("Project index not initialized"), nil
	}

	proj := projectIndex.GetProject(projectPath)
	if proj == nil {
		return mcp.NewToolResultText(fmt.Sprintf("Project not in index: %s", projectPath)), nil
	}

	if sessMgr == nil {
		return mcp.NewToolResultText(fmt.Sprintf("Project: %s\nLanguage: %s\nSessions: %d\nLast: %s",
			proj.Name, proj.Language, proj.SessionCount, proj.LastSession.Format("2006-01-02"))), nil
	}

	sessions := sessMgr.ListSessions()
	output := fmt.Sprintf("Activity for %s:\n\n", proj.Name)
	output += fmt.Sprintf("Language: %s\n", proj.Language)
	output += fmt.Sprintf("Total sessions: %d\n", proj.SessionCount)
	output += fmt.Sprintf("Last session: %s\n", proj.LastSession.Format("2006-01-02 15:04"))
	if proj.LastWork != "" {
		output += fmt.Sprintf("Last work: %s\n", proj.LastWork)
	}
	output += fmt.Sprintf("\nRecent sessions:\n")
	for _, s := range sessions {
		if len(s.Messages) > 0 {
			lastMsg := s.Messages[len(s.Messages)-1]
			output += fmt.Sprintf("- %s (%d messages, last: %s)\n", s.ID, len(s.Messages), lastMsg.Time.Format("2006-01-02 15:04"))
		}
	}

	return mcp.NewToolResultText(output), nil
}

func detectLanguage(projectPath string) string {
	files := map[string]string{
		"go.mod":      "Go",
		"package.json": "JavaScript",
		"Cargo.toml":  "Rust",
		"pyproject.toml": "Python",
		"pom.xml":    "Java",
		"Gemfile":    "Ruby",
		"*.csproj":   "C#",
	}

	for file, lang := range files {
		filepath.Join(projectPath, file)
		if _, err := os.Stat(filepath.Join(projectPath, file)); err == nil {
			return lang
		}
	}

	return "Unknown"
}