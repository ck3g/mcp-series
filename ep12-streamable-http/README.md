# Episode 12: Building a Remote MCP Server

This episode demonstrates running an MCP server with Streamable HTTP transport instead of stdio.

## Key Difference from Episode 11

| Aspect | stdio (Ep 11) | Streamable HTTP (Ep 12) |
|--------|---------------|-------------------------|
| `mcp.run()` | `mcp.run()` | `mcp.run(transport="streamable-http")` |
| Server lifecycle | Spawned per client | Runs continuously |
| State sharing | Isolated | Shared across clients |
| Endpoint | stdin/stdout | `http://localhost:8000/mcp` |

## Running the Server

```bash
# Install dependencies
uv sync

# Run the server
uv run memory_mcp.py
```

Server starts at: `http://localhost:8000/mcp`

The server stays running (unlike stdio which exits after each session)!

## Testing with MCP Inspector

1. Run `mcp dev` in another terminal
2. Open http://localhost:6274
3. Select "Streamable HTTP" transport
4. Enter URL: `http://localhost:8000/mcp`
5. Click Connect

Try opening TWO inspector windows - both see the same memories!

## Cursor Configuration

```json
{
  "mcpServers": {
    "memory-mcp-http": {
      "url": "http://localhost:8000/mcp"
    }
  }
}
```

Compare to stdio config:
```json
{
  "mcpServers": {
    "memory-mcp": {
      "command": "uv",
      "args": ["--directory", "/path/to/project", "run", "memory_mcp.py"]
    }
  }
}
```

## Tools

- `remember(key, value)` - Store a value
- `recall(key)` - Retrieve a value  
- `forget(key)` - Delete a value
- `list_memories()` - Show all stored values

## Important Notes

- Multiple clients share the same in-memory dictionary
- State is lost when the server restarts
- For persistent state across restarts, use Redis (Episode 11)

## When to Use HTTP vs stdio

| Use stdio for... | Use HTTP for... |
|-----------------|-----------------|
| Personal local tools | Team/shared tools |
| Single user | Multiple users |
| Maximum performance | Remote access |
| Spawned on demand | Always-on service |
