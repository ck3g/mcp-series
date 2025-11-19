# Episode 8: Connecting Multiple MCP Servers

This episode demonstrates how to run multiple MCP servers simultaneously in Cursor.

## Servers

### Math MCP
- **Tools:** `add`, `subtract`, `multiply`, `divide`
- **Resources:** math://constants

### Time MCP
- **Tools:** `current_time`, `current_date`, `days_until`
- **Resources:** time://timezones

### Weather MCP (from Episode 7)
- See [episode-07-weather-server/](../ep7-weather-server/) for code

## Running

1. Update paths in `cursor_config.json`
2. Add configuration to Cursor Settings > MCP Servers
3. Restart MCP servers
4. Try the demo queries!

## Demo Queries

- "What's the weather in Tokyo?"
- "What's 25 times 17?"
- "How many days until Christmas?"
- "Temperature in Berlin and days until New Year?"
- "Show me the math constants"

[MCP Series YouTube Playlist](https://www.youtube.com/playlist?list=PLdyaFFMLPMAmnacJLbWaOh6li8pN-1YoK)
