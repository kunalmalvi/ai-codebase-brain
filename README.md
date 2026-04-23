# AI Codebase Brain

> An intelligent context management MCP server for AI coding assistants that solves context window limitations and enables persistent, efficient code assistance.

## Overview

AI Codebase Brain is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server that provides AI coding assistants (Claude Code, OpenCode, Cursor, Windsurf, Codex, etc.) with powerful tools for working with large codebases:

- **Project Intelligence** - Dependency graphs, file relationships, codebase mapping
- **Context Management** - Session persistence, conversation history, smart context selection
- **Git Integration** - Commit history, file blame, diff analysis
- **Token Optimization** - Compression, deduplication, efficient context sizing
- **Monorepo Support** - Multi-package project detection and navigation

## Features

| Feature | Description |
|---------|-------------|
| **Dependency Graph** | Maps imports, exports, and module relationships across the codebase |
| **Smart Search** | Semantic search by concept (auth, database, API) rather than just keywords |
| **Session Persistence** | Remembers conversation context across restarts |
| **Token Optimization** | Reduces token usage by 20-40% through compression |
| **Multi-language** | Supports Go, JavaScript, TypeScript, Python, Rust |
| **Monorepo Detection** | Identifies npm workspaces, Lerna, Nx, pnpm workspaces |
| **Git Tools** | Commit history, file history, diff, contributors |

## Quick Start

```bash
# Clone the project
git clone https://github.com/your-repo/codebase-brain.git
cd codebase-brain

# Install dependencies
go mod download

# Build
go build -o codebase-brain ./cmd/server

# Run (connects to AI coder via stdio)
./codebase-brain
```

## Supported AI Coders

- Claude Code / Claude Desktop
- OpenCode
- Cursor
- Windsurf (Codeium)
- Codex
- Zed Editor
- VS Code (with MCP extension)
- JetBrains IDEs (IntelliJ, WebStorm, etc.)

See [SETUP.md](SETUP.md) for detailed configuration for each AI coder.

## MCP Tools (27 Tools)

### Project Intelligence

| Tool | Description |
|------|-------------|
| `search_codebase` | Search files by pattern or content |
| `get_related_files` | Find related files based on dependencies |
| `get_project_graph` | Get project structure overview |
| `index_project` | Trigger re-indexing |
| `get_file_context` | Get context for a specific file |

### Session Management

| Tool | Description |
|------|-------------|
| `create_session` | Create a new conversation session |
| `get_session` | Get session by ID |
| `list_sessions` | List all sessions |
| `add_message` | Add message to session |
| `get_conversation_history` | Get conversation history |

### Token Optimization

| Tool | Description |
|------|-------------|
| `estimate_tokens` | Estimate token count |
| `optimize_context` | Compress and deduplicate |
| `get_token_stats` | Get token usage statistics |

### Monorepo Support

| Tool | Description |
|------|-------------|
| `detect_monorepo` | Detect monorepo structure |
| `get_project_info` | Get sub-project information |

### Smart Search

| Tool | Description |
|------|-------------|
| `semantic_search` | Search by concept |
| `get_code_summary` | AI-powered code summary |
| `analyze_code` | Full code analysis |
| `detect_patterns` | Pattern detection |

### Git Integration

| Tool | Description |
|------|-------------|
| `get_git_info` | Repository info (branches, remotes) |
| `get_commits` | Recent commits |
| `get_file_history` | File commit history |
| `search_commits` | Search commit messages |
| `get_contributors` | Contributor statistics |
| `get_diff` | Diff between commits |
| `get_changed_files` | Files in a commit |

## Configuration

### Environment Variables

```bash
# Project to index (default: current directory)
export PROJECT_PATH="/path/to/project"

# Server port (default: 8080)
export PORT=8080

# Enable file watching
export WATCH=true

# Debug logging
export VERBOSE=true

# Parallel workers (default: 4)
export WORKERS=4
```

### Config File

Create `config.yaml` in your project:

```yaml
project_path: ./my-project
watch: true
verbose: true
workers: 8
port: 8080
```

## Architecture

```
codebase-brain/
├── cmd/server/          # Entry point
├── internal/
│   ├── config/         # Configuration
│   ├── graph/          # Dependency graph
│   ├── parser/         # Language parsers (Go, JS/TS, Python, Rust)
│   ├── storage/        # BadgerDB persistence
│   ├── indexer/        # File indexing
│   ├── session/        # Session persistence
│   ├── token/          # Token optimization
│   ├── monorepo/       # Monorepo detection
│   ├── mcp/            # MCP tools and handlers
│   ├── git/            # Git integration
│   └── ...
└── pkg/types/          # Shared types
```

## Requirements

- Go 1.23+
- Any AI coder with MCP support

## License

MIT

## See Also

- [SETUP.md](SETUP.md) - Detailed setup for each AI coder
- [API.md](API.md) - Complete MCP tool API reference
- [ARCHITECTURE.md](ARCHITECTURE.md) - Technical architecture
- [BUILD.md](BUILD.md) - Build instructions
