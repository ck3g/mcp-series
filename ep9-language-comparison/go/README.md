# Weather MCP Server - Go Implementation

Weather MCP server built with Go 1.25 and the official MCP Go SDK v1.1.0.

## Features

**5 Tools:**
- `get_current_weather` - Current conditions for any location
- `get_forecast` - Multi-day weather forecast (1-3 days)
- `compare_weather` - Compare weather between two locations
- `convert_temperature` - Convert between Celsius, Fahrenheit, Kelvin
- `check_rain_probability` - Rain forecast with umbrella advice

**3 Resources:**
- `weather://supported-formats` - Location format examples
- `weather://weather-codes` - Weather condition explanations
- `weather://temperature-guide` - Temperature conversion formulas

## Prerequisites

- Go 1.25 or later ([install](https://go.dev/doc/install))
- Internet connection (uses wttr.in API)

## Installation
```bash
# Navigate to the Go directory
cd go/

# Initialize Go module
go mod init weather-mcp-go

# Install dependencies
go get github.com/modelcontextprotocol/go-sdk@v1.1.0

# Build the binary
go build -o bin/weather-mcp-go main.go
```

## Testing with MCP Inspector
```bash
# Run the server
./bin/weather-mcp-go

# The server will wait for JSON-RPC input on stdin
# You can also use the MCP Inspector web tool at:
# https://github.com/modelcontextprotocol/inspector
```

## Configuration

### Cursor IDE

Add to Cursor settings (Cmd/Ctrl + Shift + P → "Preferences: Open User Settings (JSON)"):
```json
{
  "mcpServers": {
    "weather-mcp-go": {
      "command": "/absolute/path/to/weather-mcp-go"
    }
  }
}
```

**Finding your absolute path:**
```bash
cd /path/to/episode-09-language-comparison/go
pwd
# Output: /Users/yourname/projects/mcp-series/episode-09-language-comparison/go
# Full command path: /Users/yourname/projects/mcp-series/episode-09-language-comparison/go/bin/weather-mcp-go
```


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


## Development

**Rebuild after changes:**
```bash
go build -o bin/weather-mcp-go main.go
```

## Testing with MCP Inspector

The [MCP Inspector](https://github.com/modelcontextprotocol/inspector) is a web-based debugging tool for MCP servers.

### Using MCP Inspector

1. **Install and run MCP Inspector:**
```bash
   npx @modelcontextprotocol/inspector
```

2. **In the Inspector UI:**
   - Select **"STDIO"** as the transport type
   - For **Command**, enter the full path to your binary:
```
     /absolute/path/to/weather-mcp-go
```
   - Click **"Connect"**

3. **Test your server:**
   - Go to **"Tools"** tab → Click **"List"** to see all 5 tools
   - Go to **"Resources"** tab → Click **"List"** to see all 3 resources
   - Click any tool to test it with example arguments

4. **Example test:**
   - Select `get_current_weather` tool
   - Enter `{"location": "Berlin"}` as arguments
   - Click **"Run Tool"**
   - See the weather result

