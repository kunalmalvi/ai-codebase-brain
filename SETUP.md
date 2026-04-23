# Codebase Brain - Setup Guide

This guide covers setting up Codebase Brain with various AI coding assistants.

## Prerequisites

```bash
# Install Go 1.23+
go version

# Clone and build the server
git clone https://github.com/your-repo/codebase-brain.git
cd codebase-brain
go build -o codebase-brain ./cmd/server
```

---

## AI Coder Setup

### Claude Code / Claude Desktop

#### Option 1: CLI (Recommended)
```bash
# macOS/Linux
claude mcp add --transport stdio /path/to/codebase-brain

# Windows
claude mcp add --transport stdio C:\path\to\codebase-brain.exe
```

#### Option 2: Configuration File
Add to `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "codebase-brain": {
      "command": "/path/to/codebase-brain",
      "transport": "stdio"
    }
  }
}
```

#### Option 3: Project Config
Add to `.claude.json` in your project:

```json
{
  "mcpServers": {
    "codebase-brain": {
      "command": "./codebase-brain"
    }
  }
}
```

---

### OpenCode

Add to `~/.opencode/mcp.json`:

```json
{
  "mcpServers": {
    "codebase-brain": {
      "command": "./codebase-brain",
      "transport": "stdio"
    }
  }
}
```

Or use the CLI:
```bash
opencode mcp add codebase-brain ./codebase-brain
```

---

### Cursor

1. Open Settings (Cmd+, / Ctrl+,)
2. Navigate to: Features → MCP
3. Click "Add new server"
4. Enter:
   - Name: `codebase-brain`
   - Command: `/path/to/codebase-brain`

---

### Windsurf (Codeium)

```bash
# Using windsurf CLI
windsurf mcp add --transport stdio /path/to/codebase-brain

# Or via config
# Add to ~/.windsurf/config.json
{
  "mcpServers": {
    "codebase-brain": {
      "command": "/path/to/codebase-brain",
      "transport": "stdio"
    }
  }
}
```

---

### Codex

```bash
codex mcp add codebase-brain /path/to/codebase-brain
```

---

### Zed Editor

Add to `~/.zed/settings.json`:

```json
{
  "mcp_servers": {
    "codebase-brain": {
      "command": "/path/to/codebase-brain"
    }
  }
}
```

---

### VS Code (with Copilot Chat)

> VS Code requires the "MCP extension" to use MCP servers. Several options:
> - [Continue Extension](https://continue.dev) - Full AI coding assistant
> - [MCP Quick Access](https://marketplace.visualstudio.com/items?itemName=some-publisher.mcp-quick-access)

#### With Continue Extension
1. Install Continue extension
2. Add to `.continue/config.json`:

```json
{
  "models": [{
    "provider": "openai",
    "model": "gpt-4"
  }],
  "mcpServers": {
    "codebase-brain": {
      "command": "/path/to/codebase-brain",
      "args": []
    }
  }
}
```

---

### JetBrains IDEs (IntelliJ, WebStorm, GoLand, etc.)

1. Install the **MCP** plugin from JetBrains Marketplace
2. Go to Settings → Tools → MCP Server
3. Click "+" to add new server:
   - Name: `codebase-brain`
   - Command: `/path/to/codebase-brain`

---

### Vapor / Devin

```bash
# Check specific CLI docs for your AI coder
vapor mcp add codebase-brain /path/to/codebase-brain
```

---

## Configuration

### Set Project Path

```bash
# Environment variable
export PROJECT_PATH=/your/project

# Windows
$env:PROJECT_PATH = "C:\your\project"
```

### Enable Features

```bash
# File watching (auto-reindex on changes)
export WATCH=true

# Debug logging
export VERBOSE=true

# Workers for parallel processing
export WORKERS=8
```

### Config File

Create `config.yaml` in the same directory as the server:

```yaml
project_path: /path/to/your/project
watch: true
verbose: false
workers: 4
port: 8080
```

---

## Verify Installation

### Test the Server

```bash
# Run directly to verify it starts
./codebase-brain

# Or with debug output
./codebase-brain --verbose
```

### Test from AI Coder

Ask your AI coder:

```
> What files are in my project?
> Show me the project structure
> Find files related to authentication
> Create a session and remember I'm working on the login feature
```

---

## Troubleshooting

### Connection Issues

| Error | Solution |
|-------|----------|
| Server not responding | Verify the binary is executable: `chmod +x codebase-brain` |
| Path not found | Use absolute path instead of relative path |
| Permission denied | Check file permissions |

### Tool Issues

| Error | Solution |
|-------|----------|
| No results found | Run `index_project` first to build the index |
| Session not found | Check session ID is correct |
| Git tools not working | Verify project is a git repository |

### Performance Issues

| Issue | Solution |
|-------|----------|
| Slow indexing | Increase workers: `export WORKERS=8` |
| High memory | Project may be too large for indexing |
| Slow search | Re-index: `index_project` |

---

## Next Steps

- See [README.md](README.md) for feature overview
- See [API.md](API.md) for complete tool documentation
- See [BUILD.md](BUILD.md) for build instructions