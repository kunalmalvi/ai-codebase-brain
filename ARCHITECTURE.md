# Codebase Brain - Architecture

> Technical documentation for developers who want to understand or extend the server.

## Overview

Codebase Brain is built using:
- **Go 1.23+** - Primary language
- **mcp-go** - MCP server library
- **BadgerDB** - Embedded key-value storage

## System Architecture

```
┌─────────────────────────────────────────────────┐
│           Codebase Brain                    │
│                                             │
│  ┌─────────────┐    ┌─────────────────────┐  │
│  │   Config   │    │   MCP Protocol      │  │
│  │  Manager  │    │    Handler         │  │
│  └─────────┘    └─────────────────────┘  │
│                                             │
│  ┌─────────────────────────────────────┐   │
│  │         Tools (27)                   │   │
│  │  ┌─────┐ ┌──────┐ ┌───────┐ ┌────┐ │   │
│  │  │Graph│ │Search│ │Session│ │Git │ │   │
│  │  └─────┘ └──────┘ └───────┘ └────┘ │   │
│  └─────────────────────────────────────┘   │
│                                             │
│  ┌─────────────────────────────────────┐   │
│  │        Storage Layer                  │   │
│  │  ┌────────┐  ┌────────────┐        │   │
│  │  │BadgerDB│  │  File     │        │   │
│  │  │       │  │  Indexer  │        │   │
│  │  └────────┘  └────────────┘        │   │
│  └─────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

---

## Core Components

### 1. Config (`internal/config/`)

Manages configuration from:
- Environment variables
- Config file (`config.yaml`)
- CLI arguments

```go
type Config struct {
    ProjectPath string
    Port       int
    Watch      bool
    Verbose    bool
    Workers   int
}
```

### 2. Graph (`internal/graph/`)

Represents the codebase as a directed graph:
- **Nodes**: Files, functions, types
- **Edges**: Imports, exports, references

```go
type Graph struct {
    Nodes map[string]*Node
    Edges map[string][]string
}

type Node struct {
    Name     string
    Type    string  // file, function, type
    Imports []string
    Exports []string
    Metadata map[string]interface{}
}
```

### 3. Parser (`internal/parser/`)

Language-specific parsers extract:
- Imports/exports
- Function signatures
- Type definitions

Supported:
- `go_parser.go` - Go files
- `js_parser.go` - JavaScript/TypeScript
- `python_parser.go` - Python
- `rust_parser.go` - Rust

### 4. Indexer (`internal/indexer/`)

Builds and maintains the search index:
- File discovery
- Content extraction
- Graph updating

### 5. Storage (`internal/storage/`)

BadgerDB-based persistence:
- Graph storage
- Session storage
- Configuration cache

### 6. MCP Handler (`internal/mcp/`)

Maps MCP tools to Go functions:
- `tools.go` - Tool definitions
- `handlers.go` - Handler implementations

### 7. Session (`internal/session/`)

Conversation persistence:
- Session creation
- Message storage
- History retrieval

### 8. Git (`internal/git/`)

Git integration:
- Commit history
- File history
- Diff generation

---

## Data Flows

### Indexing Flow

```
User → index_project → Indexer → Parser → Graph → Storage
```

1. `index_project` called
2. `Indexer` discovers files
3. `Parser` extracts structure
4. `Graph` updated
5. Stored in BadgerDB

### Search Flow

```
User → search_codebase → Indexer.Search → Results
```

1. Query received
2. `Indexer.Search()` searches
3. Results returned

### Session Flow

```
User → create_session → SessionManager → Storage
```

1. `create_session` called
2. Session ID generated
3. Session saved to disk

---

## Module Dependencies

```
cmd/server
  └─ internal/mcp
       ├─ internal/config
       ├─ internal/graph
       ├─ internal/indexer
       ├─ internal/storage
       ├─ internal/session
       ├─ internal/token
       ├─ internal/monorepo
       ├─ internal/logger
       ├─ internal/errors
       └─ internal/git
            └─ internal/logger
```

---

## Key Interfaces

### Storage Interface

```go
type Storage interface {
    SaveGraph(projectID string, g *Graph) error
    LoadGraph(projectID string) (*Graph, error)
    Close() error
}
```

### SessionManager Interface

```go
type SessionManager interface {
    CreateSession(projectID string) *Session
    GetSession(id string) (*Session, bool)
    AddMessage(sessionID, role, content string) error
    ListSessions() []*Session
}
```

---

## Adding New Tools

1. Define tool in `internal/mcp/tools.go`:
```go
myTool := mcp.NewTool("my_tool",
    mcp.WithDescription("Description"),
    mcp.WithString("param", ...),
)
s.AddTool(myTool, handleMyTool)
```

2. Implement handler in `internal/mcp/handlers.go`:
```go
func handleMyTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Implementation
    return mcp.NewToolResultText("result"), nil
}
```

---

## Testing

```bash
# Run all tests
go test ./...

# Run specific module
go test ./internal/graph/...

# With coverage
go test -cover ./...
```

---

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PROJECT_PATH` | `.` | Project to index |
| `PORT` | `8080` | Server port |
| `WATCH` | `false` | Enable file watching |
| `VERBOSE` | `false` | Debug logging |
| `WORKERS` | `4` | Parallel workers |

---

## Performance Considerations

### Threading

- Indexer uses worker pool for parallel parsing
- Adjust `WORKERS` for CPU count

### Memory

- BadgerDB manages memory automatically
- Large projects may need more memory

### Caching

- Graph cached in memory
- Sessions cached on disk

---

## Extension Points

### Custom Parsers

Implement `Parser` interface:
```go
type Parser interface {
    Parse(filePath string) (*Node, error)
    GetImports(content string) []string
    GetExports(content string) []string
}
```

### Custom Storage

Implement `Storage` interface for different backends.

### Custom Tools

Add new tools following the pattern in `internal/mcp/`.

---

## Security

- Server runs in stdio mode only (no network by default)
- No user input execution
- File access limited to project path

---

## License

MIT - See LICENSE file