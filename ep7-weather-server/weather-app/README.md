# Episode 7: Multiple Tools and Resources

A weather MCP server demonstrating multiple tools and resources using FastMCP.

## Tools (5)
- `get_current_weather(location)` - Get current weather conditions
- `get_forecast(location, days)` - Get weather forecast for upcoming days
- `compare_weather(location1, location2)` - Compare weather between two locations
- `convert_temperature(temp, from_unit, to_unit)` - Convert between C, F, and K
- `check_rain_probability(location)` - Check if it's likely to rain today

## Resources (3)
- `weather://supported-formats` - List of valid location format examples
- `weather://weather-codes` - Explanation of weather condition meanings
- `weather://temperature-guide` - Temperature conversion formulas and reference points

## Setup
```bash
uv add "mcp[cli]" requests
```

## Run
```bash
# Test with MCP Inspector
uv run mcp dev weather_mcp.py

# Or configure in Cursor (see Episode 7 video for configuration details)
```

## Note
Uses [wttr.in](https://wttr.in) - a free weather API that requires no API key.
