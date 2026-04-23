# Codebase Brain - API Reference

Complete reference for all 27 MCP tools.

---

## Project Intelligence Tools

### search_codebase

Search the codebase for files matching a pattern or content.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `query` | string | Yes | Search query (file pattern or content) |
| `scope` | string | No | Directory scope to limit search |
| `limit` | number | No | Max results (default: 10) |

**Example:**
```
search_codebase(query="handleError", limit=5)
```

**Response:**
```
Found 3 results for 'handleError':

1. src/utils/errors.go
   Error handling utilities
2. src/api/middleware.go
   API error middleware
3. src/db/connection.go
   Database error handling
```

---

### get_related_files

Get files related to a given file based on imports and dependencies.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `file_path` | string | Yes | Path to the file |
| `depth` | number | No | Dependency depth (default: 2) |

**Example:**
```
get_related_files(file_path="src/main.go", depth=2)
```

---

### get_project_graph

Get an overview of the project structure and dependencies.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `limit` | number | No | Max nodes (default: 100) |

**Example:**
```
get_project_graph(limit=50)
```

---

### index_project

Trigger a re-index of the project to update the codebase graph.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | No | Project path (defaults to configured) |

**Example:**
```
index_project()
```

---

### get_file_context

Get relevant context for a specific file.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `file_path` | string | Yes | Path to the file |
| `context_type` | string | No | Type: 'dependencies', 'all', 'imports', 'exports' |

**Example:**
```
get_file_context(file_path="src/auth/login.go", context_type="all")
```

---

## Session Management Tools

### create_session

Create a new conversation session for this project.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project_id` | string | No | Project identifier |

**Example:**
```
create_session(project_id="my-project")
```

**Response:**
```
Created session: 1742345234-a1b2c3d4
Project: my-project
Created at: 2024-03-18 10:34:56
```

---

### get_session

Get a session by ID.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `session_id` | string | Yes | Session ID |

**Example:**
```
get_session(session_id="1742345234-a1b2c3d4")
```

---

### list_sessions

List all sessions for this project.

**Parameters:** None

**Example:**
```
list_sessions()
```

---

### add_message

Add a message to a session.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `session_id` | string | Yes | Session ID |
| `role` | string | Yes | 'user' or 'assistant' |
| `content` | string | Yes | Message content |

**Example:**
```
add_message(session_id="1742345234-a1b2c3d4", role="user", content="I'm working on the login feature")
```

---

### get_conversation_history

Get conversation history for a session.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `session_id` | string | Yes | Session ID |

**Example:**
```
get_conversation_history(session_id="1742345234-a1b2c3d4")
```

---

## Token Optimization Tools

### estimate_tokens

Estimate token count for a given text.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `text` | string | Yes | Text to estimate |

**Example:**
```
estimate_tokens(text="function main() { console.log('hello'); }")
```

**Response:**
```
Estimated tokens: 12 (characters: 45)
```

---

### optimize_context

Optimize context by deduplicating and compressing.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `text` | string | Yes | Text to optimize |
| `max_tokens` | number | No | Max tokens (default: 32000) |

**Example:**
```
optimize_context(text="...long code...", max_tokens=32000)
```

---

### get_token_stats

Get token usage statistics for the project.

**Parameters:** None

**Example:**
```
get_token_stats()
```

---

## Monorepo Tools

### detect_monorepo

Detect if the project is a monorepo and list sub-projects.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | No | Project path |

**Example:**
```
detect_monorepo()
```

**Response:**
```
Monorepo detected with 3 projects:

1. packages/ui (Type: npm workspace)
   Path: packages/ui
   Modules: 12
2. packages/api
3. packages/shared
```

---

### get_project_info

Get information about a specific sub-project.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to sub-project |

**Example:**
```
get_project_info(path="packages/ui")
```

---

## Smart Search Tools

### semantic_search

Search code semantically by concept or feature.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `query` | string | Yes | Search concept (e.g., 'auth', 'database') |
| `file_path` | string | No | Limit to specific file |

**Example:**
```
semantic_search(query="authentication")
```

**Response:**
```
Semantic search for: 'authentication'

Found related concepts:
  • authentication (category: security)
    Related: authorization, security, jwt
  • jwt (category: security)
  • security (category: security)
```

---

### get_code_summary

Get an AI-powered summary of a code file.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `file_path` | string | Yes | Path to file |

**Example:**
```
get_code_summary(file_path="src/auth/login.go")
```

**Response:**
```
src/auth/login.go

Language: Go | Lines: 156 | Complexity: medium

Purpose: Handle user authentication

Features:
  - 3 function(s): Login, Logout, ValidateToken
  - HTTP handlers for /auth/login, /auth/logout
  - JWT token generation

Dependencies:
  - github.com/golang-jwt/jwt/v5
  - golang.org/x/crypto
```

---

### analyze_code

Analyze code for patterns, issues, and insights.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `file_path` | string | Yes | Path to file |
| `include_patterns` | boolean | No | Include pattern detection |
| `include_summary` | boolean | No | Include code summary |

**Example:**
```
analyze_code(file_path="src/main.go", include_patterns=true, include_summary=true)
```

---

### detect_patterns

Detect code patterns and potential issues.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `file_path` | string | Yes | Path to file |
| `category` | string | No | Filter: 'security', 'quality', 'style', 'all' |

**Example:**
```
detect_patterns(file_path="src/utils/auth.go", category="security")
```

---

## Git Integration Tools

### get_git_info

Get information about the git repository.

**Parameters:** None

**Example:**
```
get_git_info()
```

**Response:**
```
Git Repository: /path/to/project

Active branch: main

Branches:
  * main
  feature/auth
  bugfix/login

Remotes: origin
```

---

### get_commits

Get recent commits from the repository.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `limit` | number | No | Max commits (default: 10) |

**Example:**
```
get_commits(limit=5)
```

---

### get_file_history

Get commit history for a specific file.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `file_path` | string | Yes | Path to file |
| `limit` | number | No | Max commits (default: 10) |

**Example:**
```
get_file_history(file_path="src/main.go", limit=10)
```

---

### search_commits

Search commits by message text.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `query` | string | Yes | Search query |
| `limit` | number | No | Max commits (default: 10) |

**Example:**
```
search_commits(query="auth", limit=5)
```

---

### get_contributors

Get contributor statistics.

**Parameters:** None

**Example:**
```
get_contributors()
```

---

### get_diff

Get diff between two commits or references.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `from` | string | No | From commit/ref |
| `to` | string | No | To commit/ref |

**Example:**
```
get_diff(from="HEAD~5", to="HEAD")
```

---

### get_changed_files

Get files changed in a commit.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `commit_hash` | string | Yes | Commit hash |

**Example:**
```
get_changed_files(commit_hash="abc123...")
```

---

## Error Responses

All tools return errors in this format:

```json
{
  "error": "Error message description"
}
```

**Common errors:**
| Error | Cause |
|-------|-------|
| `Session not found` | Invalid session ID |
| `Git analyzer not initialized` | Project not a git repository |
| `Indexer not initialized` | Run index_project first |
| `Project may not be a git repository` | No .git directory |