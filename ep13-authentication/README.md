# Episode 13: Securing Your Remote MCP Server

Protect your remote MCP server with token authentication using FastMCP 2.0.

## What's Covered

- Migrating from `mcp[cli]` to FastMCP 2.0 (one import change!)
- Adding token authentication with `StaticTokenVerifier`
- Configuring Cursor with authorization headers
- Testing authentication in MCP Inspector

## The Migration

**Before (Episode 12 - mcp[cli]):**
```python
from mcp.server.fastmcp import FastMCP
```

**After (Episode 13 - FastMCP 2.0):**
```python
from fastmcp import FastMCP
```

Same code, different import - FastMCP 2.0 adds authentication support.

## Quick Start

```bash
# Install dependencies
cd ep13-authentication
uv sync

# Run the server
uv run memory_mcp.py
```

Server starts at `http://localhost:8000/mcp`

## Token

The server uses a simple token for authentication:

```
my-secret-token
```

## Cursor Configuration

```json
{
  "mcpServers": {
    "memory-mcp": {
      "url": "http://localhost:8000/mcp",
      "headers": {
        "Authorization": "Bearer my-secret-token"
      }
    }
  }
}
```


## What's Next?

For production deployments, FastMCP 2.0 also supports:
- OAuth 2.0 providers (GitHub, Google, Azure, Auth0)
- JWT verification
- Custom authentication providers

See [FastMCP Authentication Docs](https://gofastmcp.com/servers/auth/authentication) for more options.
