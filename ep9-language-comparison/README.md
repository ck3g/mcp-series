# Episode 9: MCP in Python vs Go vs Ruby

Build the same weather MCP server in three languages to compare approaches.

## Implementations

### Python (FastMCP)
- **Location:** `../ep7-weather-server/weather-app/weather_mcp.py`
- **Run:** `uv run --directory /path/to/ep7-weather-server/weather-app weather_mcp.py`

### Go (SDK v1.1.0)
- **Location:** `go/`
- **Build:** `go build -o bin/weather-mcp-go main.go`
- **Run:** `./bin/weather-mcp-go`

### Ruby (SDK v0.4.0)
- **Location:** `ruby/`
- **Run:** `bundle exec ruby weather_mcp.rb`


## All Servers Provide

- 5 tools (weather queries + temperature conversion)
- 3 resources (format guide, weather codes, temp guide)
- Same functionality via [wttr.in](https://wttr.in) API

## Run All Three in Cursor

Configure all three servers in `mcp.json` to compare them side-by-side. See individual READMEs for configuration details.

