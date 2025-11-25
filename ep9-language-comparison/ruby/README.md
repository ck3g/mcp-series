# Weather MCP Server - Ruby

Weather information tools for MCP using Ruby SDK v0.4.0.

## Quick Start

```bash
# Install dependencies
mise install
bundle install

# Run server
bundle exec ruby weather_mcp.rb
```

## Tools (5)

- `get_current_weather(location)` - Current conditions
- `get_forecast(location, days)` - Multi-day forecast (1-3 days)
- `compare_weather(location1, location2)` - Compare two locations
- `convert_temperature(temp, from_unit, to_unit)` - C/F/K conversion
- `check_rain_probability(location)` - Rain chance + umbrella advice

## Resources (3)

- `weather://supported-formats` - Location format examples
- `weather://weather-codes` - Weather condition codes
- `weather://temperature-guide` - Conversion formulas

## Cursor Configuration

**Important:** Must use `sh -c` with `cd` so bundler can find gems.

```json
{
  "mcpServers": {
    "weather-mcp-ruby": {
      "command": "sh",
      "args": [
        "-c",
        "cd /absolute/path/to/ruby && bundle exec ruby weather_mcp.rb"
      ]
    }
  }
}
```

Replace `/absolute/path/to/ruby` with your actual path.

## Usage Examples

Once connected to Cursor:

**Current weather:**
```
What's the weather in Berlin right now?
```

**Forecast:**
```
Give me a 3-day forecast for Tokyo
```

**Compare locations:**
```
Compare the weather in London and New York
```

**Temperature conversion:**
```
Convert 72°F to Celsius
```

**Rain probability:**
```
Should I bring an umbrella in Seattle today?
```

**Access resources:**
```
What location formats do you support?
```

## Data Source

Uses [wttr.in](https://wttr.in) weather API (no key required).

