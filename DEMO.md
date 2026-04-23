# Codebase Brain - Demo & Examples

Interactive examples showing what you can do with Codebase Brain.

---

## Quick Demo

### Try These Commands

After connecting to your AI coder, try:

```
> Search for files containing "handleError"
```

```
> Show me the project structure
```

```
> Find files related to authentication
```

```
> Create a session and remember I'm working on the login feature
```

```
> Get recent commits in this project
```

---

## Sample Responses

### Search Results
```
Found 3 results for 'handleError':

1. src/utils/errors.go
   Error handling utilities
2. src/api/middleware.go  
   API error middleware
3. src/db/connection.go
   Database error handling
```

### Project Graph
```
Project structure (23 files):

- cmd/server/main.go (Go)
- internal/config/ (Go)
- internal/mcp/ (Go)
- internal/graph/ (Go)
- internal/parser/ (Go)
- package.json (JavaScript)
- src/components/ (TypeScript)
```

### Semantic Search
```
Semantic search for: 'authentication'

Found related concepts:
  • authentication (category: security)
    Related: authorization, jwt, security
  • jwt (category: security)
  • security (category: security)
```

### Session Management
```
Created session: 1742345234-a1b2c3d4
Project: my-project
Created at: 2024-03-18 10:34:56
```

```
Added user message to session 1742345234-a1b2c3d4
```

### Git History
```
Recent 5 commits:

a1b2c3d | alice | 2024-03-18
  Add authentication feature

d4e5f6a | bob | 2024-03-17
  Fix login redirect

7g8h9i0 | alice | 2024-03-17
  Update user model

...
```

### Token Optimization
```
Original: ~4500 tokens
Optimized: ~3200 tokens
Savings: 28.9%

[Optimized code...]
```

---

## Real-World Scenarios

### Scenario 1: Understanding a New Codebase

```
> What files are in this project?
> Show me the dependency graph
> Find files related to database
> Summarize src/db/connection.go
```

### Scenario 2: Continuing Previous Work

```
> Create a session
> What was I working on last time?
> Show conversation history
```

### Scenario 3: Finding Code

```
> Find files related to authentication
> Show me test files for auth
> Get the related files for src/login.go
```

### Scenario 4: Git Investigation

```
> Show me recent commits
> What changed in main.go?
> Who wrote the most code?
```

---

## Before vs After

| Task | Without Server | With Server |
|------|-----------------|-------------|
| Find auth files | `grep -r` manual search | `semantic_search("auth")` |
| Understand codebase | Read files manually | `get_project_graph` |
| Remember context | Can't | `create_session` |
| Token usage | All files sent | Optimized -30% |
| Find related tests | Manual discovery | `get_related_files` |
| Git history | Check terminal separately | `get_commits` |
| Monorepo | Manual exploration | `detect_monorepo` |

---

## Testing the Server

```bash
# Build
go build -o codebase-brain ./cmd/server

# Test stdio mode (should just start and wait for input)
./codebase-brain

# Test with index
PROJECT_PATH=/your/project ./codebase-brain

# Verbose
VERBOSE=true ./codebase-brain
```

---

## Debugging

### Check Index Status
```
> Index this project
> Show me the project structure
```

### Check Sessions
```
> List all sessions
> Get conversation history for [session-id]
```

### Check Git
```
> Get git info
> Show recent commits
```

---

## Common Issues

| Issue | Solution |
|-------|----------|
| No search results | Run `index_project` first |
| Session not found | Check session ID spelling |
| Git not working | Verify .git directory exists |
| Slow performance | Increase WORKERS |