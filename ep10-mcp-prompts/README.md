# Episode 10: MCP Prompts - Reusable AI Instructions

A weather MCP server demonstrating **tools, resources, AND prompts** using FastMCP.
This extends Episode 7's weather server by adding 3 prompts that teach Cursor how to use the tools intelligently.

## What's New in Episode 10

**3 Prompts added:**
- `weather_planning_assistant` - Analyzes forecast and suggests activities + packing list
- `travel_weather_comparison` - Compares two cities and recommends which to visit
- `temperature_context` - Explains temperatures in multiple units with real-world context

These prompts combine multiple tools and resources into pre-programmed workflows.

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

## Prompts (3) NEW!
- `weather_planning_assistant` - Plan activities based on 3-day forecast
- `travel_weather_comparison` - Compare cities for travel decisions
- `temperature_context` - Explain temperatures with context and clothing advice

## Setup
```bash
uv add "mcp[cli]" requests
```

## Run
```bash
# Test with MCP Inspector
uv run mcp dev weather_mcp.py

# Configure in Cursor
```

## Cursor Configuration
```json
{
  "mcpServers": {
    "weather-mcp": {
      "command": "/absolute/path/to/uv",
      "args": [
        "--directory",
        "/absolute/path/to/ep10-mcp-prompts/weather-app/",
        "run",
        "weather_mcp.py"
      ]
    }
  }
}
```

## Try These Prompts in Cursor
1. Use the "weather_planning_assistant" prompt with "Berlin"
2. Use "travel_weather_comparison" with "Tokyo" and "San Francisco"
3. Use "temperature_context" with "25°C"

Watch how Cursor automatically uses the right tools and resources!

## Note
Uses [wttr.in](https://wttr.in) - a free weather API that requires no API key.
