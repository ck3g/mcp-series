from mcp.server.fastmcp import FastMCP

import redis

mcp = FastMCP("memory-mcp")

REDIS_PREFIX = "mcp:memory:"
r = redis.Redis(host="localhost", port=6379, decode_responses=True)


def _key(name: str) -> str:
    """Add prefix to key for namespacing."""
    return f"{REDIS_PREFIX}{name}"


@mcp.tool()
def remember(key: str, value: str) -> str:
    """Store a piece of information in memory.

    Args:
        key: A short identifier for this memory (e.g., "favorite_color", "birthday")
        value: The information to remember
    """
    r.set(_key(key), value)
    return f"Remembered '{key}' = '{value}'"


@mcp.tool()
def recall(key: str) -> str:
    """Retrieve a previously stored memory.

    Args:
        key: The identifier for the memory to recall
    """
    value = r.get(_key(key))
    if value is not None:
        return f"'{key}' = '{value}'"
    return f"I don't have any memory of '{key}'"


@mcp.tool()
def forget(key: str) -> str:
    """Remove a memory.

    Args:
        key: The identifier for the memory to forget
    """
    deleted = r.delete(_key(key))
    if deleted:
        return f"Forgot '{key}'"
    return f"I don't have any memory of '{key}' to forget"


@mcp.tool()
def list_memories() -> str:
    """List all stored memories."""
    # Get all keys with our prefix
    keys: list[str] = r.keys(f"{REDIS_PREFIX}*")  # type: ignore[assignment]

    if not keys:
        return "No memories stored yet."

    lines = []
    for full_key in keys:
        short_key = full_key.replace(REDIS_PREFIX, "")
        value = r.get(full_key)
        lines.append(f"  * {short_key} = {value}")

    return f"Current memories ({len(keys)}):\n" + "\n".join(lines)


if __name__ == "__main__":
    mcp.run()
