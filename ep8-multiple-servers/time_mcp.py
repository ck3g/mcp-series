from datetime import datetime

from mcp.server.fastmcp import FastMCP

mcp = FastMCP("time-mcp")


@mcp.tool()
def current_time() -> str:
    """Get the current time in HH:MM:SS format"""
    return datetime.now().strftime("%H:%M:%S")


@mcp.tool()
def current_date() -> str:
    """Get the current date in YYYY-MM-DD format"""
    return datetime.now().strftime("%Y-%m-%d")


@mcp.tool()
def days_until(target_date: str) -> int:
    """Calculate days until a target date

    Args:
        target_date: Target date in YYYY-MM-DD format
    """
    target = datetime.strptime(target_date, "%Y-%m-%d")
    today = datetime.now().replace(hour=0, minute=0, second=0, microsecond=0)
    delta = target - today
    return delta.days


@mcp.resource("time://timezones")
def get_timezones() -> str:
    """Common timezone abbreviations and UTC offsets"""
    return """Common Timezones:

    UTC - Coordinated Universal Time (UTC±0)
    EST - Eastern Standard Time (UTC-5)
    EDT - Eastern Daylight Time (UTC-4)
    PST - Pacific Standard Time (UTC-8)
    PDT - Pacific Daylight Time (UTC-7)
    CET - Central European Time (UTC+1)
    CEST - Central European Summer Time (UTC+2)
    JST - Japan Standard Time (UTC+9)
    AEST - Australian Eastern Standard Time (UTC+10)
    """


if __name__ == "__main__":
    mcp.run()
