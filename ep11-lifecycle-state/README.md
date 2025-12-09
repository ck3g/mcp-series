# Episode 11: MCP Server Lifecycle & State

MCP servers using stdio transport spawn a **separate process for each client connection**.
This means in-memory state doesn't persist across chat sessions.

## Two Implementations

### `naive/` - In-Memory Storage
Uses a Python dictionary. Works within a single chat, but **loses data** when you open a new chat (new process = empty dictionary).

### `redis/` - External Storage
Uses Redis for persistence. Data survives across chat sessions and server restarts.

## Redis Setup

```bash
# macOS
brew install redis
brew services start redis

# Docker
docker run -d -p 6379:6379 redis
```

## Cursor Configuration

```json
{
  "mcpServers": {
    "memory-mcp": {
      "command": "/path/to/uv",
      "args": [
        "--directory",
        "/path/to/ep11-lifecycle-state/naive/",
        "run",
        "memory_mcp.py"
      ]
    }
  }
}
```

Switch `naive/` to `redis/` to use persistent storage.

