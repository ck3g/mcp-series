from fastmcp import FastMCP
from fastmcp.server.auth.providers.jwt import StaticTokenVerifier

memories: dict[str, str] = {}

verifier = StaticTokenVerifier(
    tokens={
        "my-secret-token": {
            "client_id": "developer",
            "scopes": ["read:data", "write:data"],
        }
    },
    required_scopes=["read:data"],
)

mcp = FastMCP(name="memory-mcp", auth=verifier)


@mcp.tool()
def remember(key: str, value: str) -> str:
    """Store a piece of information in memory.

    Args:
        key: A short identifier for this memory (e.g., "favorite_color", "birthday")
        value: The information to remember
    """
    memories[key] = value
    return f"Remembered '{key}' = '{value}'"


@mcp.tool()
def recall(key: str) -> str:
    """Retrieve a previously stored memory.

    Args:
        key: The identifier for the memory to recall
    """
    if key in memories:
        return f"'{key}' = '{memories[key]}'"
    return f"I don't have any memory of '{key}'"


@mcp.tool()
def forget(key: str) -> str:
    """Remove a memory.

    Args:
        key: The identifier for the memory to forget
    """
    if key in memories:
        del memories[key]
        return f"Forgot '{key}'"
    return f"I don't have any memory of '{key}' to forget"


@mcp.tool()
def list_memories() -> str:
    """List all stored memories."""
    if not memories:
        return "No memories stored yet."

    lines = [f"  * {k} = {v}" for k, v in memories.items()]
    return f"Current memories ({len(memories)}):\n" + "\n".join(lines)


if __name__ == "__main__":
    mcp.run(transport="streamable-http")
