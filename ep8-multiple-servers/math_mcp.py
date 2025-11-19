from mcp.server.fastmcp import FastMCP

mcp = FastMCP("math-mcp")


@mcp.tool()
def add(a: float, b: float) -> float:
    """Add two numbers"""
    return a + b


@mcp.tool()
def subtract(a: float, b: float) -> float:
    """Subtract b from a"""
    return a - b


@mcp.tool()
def multiply(a: float, b: float) -> float:
    """Multiply two numbers"""
    return a * b


@mcp.tool()
def divide(a: float, b: float) -> float:
    """Divide a by b"""
    if b == 0:
        raise ValueError("Cannot divide by zero")
    return a / b


@mcp.resource("math://constants")
def get_constants() -> str:
    """Common mathematical constants"""
    return """
    π (pi) = 3.14159265359...
    e = 2.71828182846...
    φ (phi, golden ratio) = 1.61803398875...
    √2 = 1.41421356237...
    """


if __name__ == "__main__":
    mcp.run()
