package mcp

import (
	"github.com/codebase-brain/internal/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTools(s *server.MCPServer, cfg *config.Config) {
	searchTool := mcp.NewTool("search_codebase",
		mcp.WithDescription("Search the codebase for files matching a pattern or content query"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query: a file pattern or content to find"),
		),
		mcp.WithString("scope",
			mcp.Description("Directory scope to limit the search"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results to return"),
		),
	)

	getRelatedTool := mcp.NewTool("get_related_files",
		mcp.WithDescription("Get files related to a given file based on imports and dependencies"),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("Path to the file to find related files for"),
		),
		mcp.WithNumber("depth",
			mcp.Description("Depth of dependency traversal (default: 2)"),
		),
	)

	getProjectGraphTool := mcp.NewTool("get_project_graph",
		mcp.WithDescription("Get an overview of the project structure and dependencies"),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of nodes to return"),
		),
	)

	indexProjectTool := mcp.NewTool("index_project",
		mcp.WithDescription("Trigger a re-index of the project to update the codebase graph"),
		mcp.WithString("path",
			mcp.Description("Project path to index (defaults to configured path)"),
		),
	)

	getFileContextTool := mcp.NewTool("get_file_context",
		mcp.WithDescription("Get relevant context for a specific file including its dependencies and recent changes"),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("Path to the file to get context for"),
		),
		mcp.WithString("context_type",
			mcp.Description("Type of context: 'dependencies', 'all', 'imports', 'exports'"),
		),
	)

	// Session management tools
	createSessionTool := mcp.NewTool("create_session",
		mcp.WithDescription("Create a new conversation session for this project"),
		mcp.WithString("project_id",
			mcp.Description("Project identifier (defaults to project path)"),
		),
	)

	getSessionTool := mcp.NewTool("get_session",
		mcp.WithDescription("Get a session by ID"),
		mcp.WithString("session_id",
			mcp.Required(),
			mcp.Description("Session ID to retrieve"),
		),
	)

	listSessionsTool := mcp.NewTool("list_sessions",
		mcp.WithDescription("List all sessions for this project"),
	)

	addMessageTool := mcp.NewTool("add_message",
		mcp.WithDescription("Add a message to a session"),
		mcp.WithString("session_id",
			mcp.Required(),
			mcp.Description("Session ID to add message to"),
		),
		mcp.WithString("role",
			mcp.Required(),
			mcp.Description("Role: 'user' or 'assistant'"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Message content"),
		),
	)

	getConversationHistoryTool := mcp.NewTool("get_conversation_history",
		mcp.WithDescription("Get the conversation history for a session"),
		mcp.WithString("session_id",
			mcp.Required(),
			mcp.Description("Session ID to get history from"),
		),
	)

	s.AddTool(searchTool, handleSearch)
	s.AddTool(getRelatedTool, handleGetRelated)
	s.AddTool(getProjectGraphTool, handleGetProjectGraph)
	s.AddTool(indexProjectTool, handleIndexProject)
	s.AddTool(getFileContextTool, handleGetFileContext)
	s.AddTool(createSessionTool, handleCreateSession)
	s.AddTool(getSessionTool, handleGetSession)
	s.AddTool(listSessionsTool, handleListSessions)
	s.AddTool(addMessageTool, handleAddMessage)
	s.AddTool(getConversationHistoryTool, handleGetConversationHistory)

	// Token optimization tools
	estimateTokensTool := mcp.NewTool("estimate_tokens",
		mcp.WithDescription("Estimate token count for a given text"),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text to estimate tokens for"),
		),
	)

	optimizeContextTool := mcp.NewTool("optimize_context",
		mcp.WithDescription("Optimize context by deduplicating and compressing"),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text to optimize"),
		),
		mcp.WithNumber("max_tokens",
			mcp.Description("Maximum tokens allowed (default: 32000)"),
		),
	)

	getTokenStatsTool := mcp.NewTool("get_token_stats",
		mcp.WithDescription("Get token usage statistics for the project"),
	)

	s.AddTool(estimateTokensTool, handleEstimateTokens)
	s.AddTool(optimizeContextTool, handleOptimizeContext)
	s.AddTool(getTokenStatsTool, handleGetTokenStats)

	// Monorepo tools
	detectMonorepoTool := mcp.NewTool("detect_monorepo",
		mcp.WithDescription("Detect if the project is a monorepo and list sub-projects"),
		mcp.WithString("path",
			mcp.Description("Project path to analyze (defaults to configured path)"),
		),
	)

	getProjectInfoTool := mcp.NewTool("get_project_info",
		mcp.WithDescription("Get information about a specific sub-project in a monorepo"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Path to the sub-project"),
		),
	)

	s.AddTool(detectMonorepoTool, handleDetectMonorepo)
	s.AddTool(getProjectInfoTool, handleGetProjectInfo)

	// Smart search tools
	analyzeCodeTool := mcp.NewTool("analyze_code",
		mcp.WithDescription("Analyze code for patterns, issues, and get AI-powered insights"),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("Path to the file to analyze"),
		),
		mcp.WithBoolean("include_patterns",
			mcp.Description("Include pattern detection (security, code quality)"),
		),
		mcp.WithBoolean("include_summary",
			mcp.Description("Include code summarization"),
		),
	)

	semanticSearchTool := mcp.NewTool("semantic_search",
		mcp.WithDescription("Search code semantically by concept or feature"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query (e.g., 'authentication', 'database', 'async')"),
		),
		mcp.WithString("file_path",
			mcp.Description("Limit search to specific file"),
		),
	)

	getCodeSummaryTool := mcp.NewTool("get_code_summary",
		mcp.WithDescription("Get an AI-powered summary of a code file"),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("Path to the file to summarize"),
		),
	)

	detectPatternsTool := mcp.NewTool("detect_patterns",
		mcp.WithDescription("Detect code patterns and potential issues"),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("Path to the file to analyze"),
		),
		mcp.WithString("category",
			mcp.Description("Filter by category: security, quality, style, all"),
		),
	)

	s.AddTool(analyzeCodeTool, handleAnalyzeCode)
	s.AddTool(semanticSearchTool, handleSemanticSearch)
	s.AddTool(getCodeSummaryTool, handleGetCodeSummary)
	s.AddTool(detectPatternsTool, handleDetectPatterns)

	// Git integration tools
	getGitInfoTool := mcp.NewTool("get_git_info",
		mcp.WithDescription("Get information about the git repository (branches, remotes)"),
	)

	getCommitsTool := mcp.NewTool("get_commits",
		mcp.WithDescription("Get recent commits from the repository"),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of commits to return (default: 10)"),
		),
	)

	getFileHistoryTool := mcp.NewTool("get_file_history",
		mcp.WithDescription("Get the commit history for a specific file"),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("Path to the file to get history for"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of commits to return (default: 10)"),
		),
	)

	searchCommitsTool := mcp.NewTool("search_commits",
		mcp.WithDescription("Search commits by message text"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query for commit messages"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of commits to return (default: 10)"),
		),
	)

	getContributorsTool := mcp.NewTool("get_contributors",
		mcp.WithDescription("Get contributor statistics for the repository"),
	)

	getDiffTool := mcp.NewTool("get_diff",
		mcp.WithDescription("Get the diff between two commits or references"),
		mcp.WithString("from",
			mcp.Description("From commit or reference"),
		),
		mcp.WithString("to",
			mcp.Description("To commit or reference"),
		),
	)

	getChangedFilesTool := mcp.NewTool("get_changed_files",
		mcp.WithDescription("Get the list of files changed in a commit"),
		mcp.WithString("commit_hash",
			mcp.Required(),
			mcp.Description("The commit hash to get changed files for"),
		),
	)

	s.AddTool(getGitInfoTool, handleGetGitInfo)
	s.AddTool(getCommitsTool, handleGetCommits)
	s.AddTool(getFileHistoryTool, handleGetFileHistory)
	s.AddTool(searchCommitsTool, handleSearchCommits)
	s.AddTool(getContributorsTool, handleGetContributors)
	s.AddTool(getDiffTool, handleGetDiff)
	s.AddTool(getChangedFilesTool, handleGetChangedFiles)
}