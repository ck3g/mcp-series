from mcp.server.fastmcp import FastMCP

mcp = FastMCP("hello-mcp")


@mcp.tool()
def get_greetings(name: str) -> str:
    """Get a personalized greeting

    Args:
        name: Name to greet
    """
    return f"Hello, {name}! Welcome to MCP!"
