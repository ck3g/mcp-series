# frozen_string_literal: true
# rubocop:disable all

require 'mcp'
require 'net/http'
require 'json'
require 'uri'

# ============================================================================
# Helper method to fetch weather data
# ============================================================================

def fetch_weather(location)
  uri = URI("https://wttr.in/#{URI.encode_www_form_component(location)}?format=j1")
  response = Net::HTTP.get_response(uri)

  raise "Weather API error: #{response.code}" unless response.is_a?(Net::HTTPSuccess)

  JSON.parse(response.body)
end

# ============================================================================
# Tool 1: Get Current Weather
# ============================================================================

class GetCurrentWeatherTool < MCP::Tool
  tool_name 'get_current_weather'
  description 'Get current weather conditions for a location'

  input_schema(
    properties: {
      location: {
        type: 'string',
        description: 'City name, airport code, or coordinates'
      }
    },
    required: ['location']
  )

  def self.call(location:, server_context:)
    data = fetch_weather(location)
    current = data['current_condition'][0]

    text = <<~WEATHER
      Weather in #{location}:
      Temperature: #{current['temp_C']}°C (#{current['temp_F']}°F)
      Conditions: #{current['weatherDesc'][0]['value']}
      Feels like: #{current['FeelsLikeC']}°C
      Humidity: #{current['humidity']}%
      Wind: #{current['windspeedKmph']} km/h
    WEATHER

    MCP::Tool::Response.new([{ type: 'text', text: text.strip }])
  end
end

# ============================================================================
# Tool 2: Get Forecast
# ============================================================================

class GetForecastTool < MCP::Tool
  tool_name 'get_forecast'
  description 'Get weather forecast for upcoming days'

  input_schema(
    properties: {
      location: {
        type: 'string',
        description: 'City name, airport code, or coordinates'
      },
      days: {
        type: 'integer',
        description: 'Number of days for forecast (1-3)',
        default: 3,
        minimum: 1,
        maximum: 3
      }
    },
    required: ['location']
  )

  def self.call(location:, server_context:, days: 3)
    data = fetch_weather(location)
    forecast = data['weather'][0...days]

    result = "Forecast for #{location}:\n\n"

    forecast.each do |day|
      result += "#{day['date']}:\n"
      result += "  High: #{day['maxtempC']}°C, Low: #{day['mintempC']}°C\n"
      result += "  Conditions: #{day['hourly'][4]['weatherDesc'][0]['value']}\n"
      result += "  Precipitation: #{day['hourly'][4]['precipMM']}mm\n\n"
    end

    MCP::Tool::Response.new([{ type: 'text', text: result }])
  end
end

# ============================================================================
# Tool 3: Compare Weather
# ============================================================================

class CompareWeatherTool < MCP::Tool
  tool_name 'compare_weather'
  description 'Compare current weather between two locations'

  input_schema(
    properties: {
      location1: {
        type: 'string',
        description: 'First location to compare'
      },
      location2: {
        type: 'string',
        description: 'Second location to compare'
      }
    },
    required: %w[location1 location2]
  )

  def self.call(location1:, location2:, server_context:)
    data1 = fetch_weather(location1)
    data2 = fetch_weather(location2)

    temp1 = data1['current_condition'][0]['temp_C'].to_i
    temp2 = data2['current_condition'][0]['temp_C'].to_i

    cond1 = data1['current_condition'][0]['weatherDesc'][0]['value']
    cond2 = data2['current_condition'][0]['weatherDesc'][0]['value']

    comparison = if temp1 > temp2
                   'warmer'
                 elsif temp1 < temp2
                   'cooler'
                 else
                   'the same temperature as'
                 end

    diff = (temp1 - temp2).abs

    text = <<~COMPARISON
      #{location1}: #{temp1}°C, #{cond1}
      #{location2}: #{temp2}°C, #{cond2}

      #{location1} is #{comparison} #{location2} (#{diff}°C difference)
    COMPARISON

    MCP::Tool::Response.new([{
      type: 'text',
      text: text.strip
    }])
  end
end

# ============================================================================
# Tool 4: Convert Temperature
# ============================================================================

class ConvertTemperatureTool < MCP::Tool
  tool_name 'convert_temperature'
  description 'Convert temperature between Celsius, Fahrenheit, and Kelvin'

  input_schema(
    properties: {
      temp: {
        type: 'number',
        description: 'Temperature value to convert'
      },
      from_unit: {
        type: 'string',
        description: 'Source unit: C, F, or K',
        enum: %w[C F K]
      },
      to_unit: {
        type: 'string',
        description: 'Target unit: C, F, or K',
        enum: %w[C F K]
      }
    },
    required: %w[temp from_unit to_unit]
  )

  def self.call(temp:, from_unit:, to_unit:, server_context:)
    # Convert to Celsius first
    celsius = case from_unit
              when 'F'
                (temp - 32) * 5.0 / 9.0
              when 'K'
                temp - 273.15
              else
                temp
              end

    # Convert from Celsius to target
    result = case to_unit
             when 'F'
               celsius * 9.0 / 5.0 + 32
             when 'K'
               celsius + 273.15
             else
               celsius
             end

    text = "#{temp}°#{from_unit} = #{result.round(1)}°#{to_unit}"

    MCP::Tool::Response.new([{
      type: 'text',
      text: text
    }])
  end
end

# ============================================================================
# Tool 5: Check Rain Probability
# ============================================================================

class CheckRainProbabilityTool < MCP::Tool
  tool_name 'check_rain_probability'
  description "Check if it's likely to rain today"

  input_schema(
    properties: {
      location: {
        type: 'string',
        description: 'City name or coordinates'
      }
    },
    required: ['location']
  )

  def self.call(location:, server_context:)
    data = fetch_weather(location)
    today = data['weather'][0]

    # Average precipitation chance across the day
    precip_chances = today['hourly'].map { |hour| hour['chanceofrain'].to_i }
    avg_chance = precip_chances.sum / precip_chances.size.to_f

    level, advice = if avg_chance > 70
                      ['High', 'Bring an umbrella!']
                    elsif avg_chance > 40
                      ['Moderate', 'Maybe bring an umbrella']
                    else
                      ['Low', 'You should be fine without an umbrella']
                    end

    text = "#{level} chance of rain in #{location} today (#{avg_chance.round}%)\n#{advice}"

    MCP::Tool::Response.new([{
      type: 'text',
      text: text
    }])
  end
end

# ============================================================================
# Resources
# ============================================================================

SUPPORTED_FORMATS = <<~TEXT
  Supported location formats:

  1. City name: London, Paris, Tokyo, "New York"
  2. Airport code: JFK, LAX, LHR
  3. Coordinates: 51.5074,-0.1278 (latitude,longitude)
  4. US location: NewYork,NY or Seattle,WA

  Examples:
  - get_current_weather("London")
  - get_current_weather("JFK")
  - get_current_weather("51.5074,-0.1278")
  - get_current_weather("Seattle,WA")
TEXT

WEATHER_CODES = <<~TEXT
  Common weather conditions:

  Clear/Sunny - No clouds, clear skies
  Partly Cloudy - Some clouds, mostly clear
  Cloudy - Overcast skies
  Rain - Precipitation falling
  Thunderstorm - Rain with lightning/thunder
  Snow - Frozen precipitation
  Fog/Mist - Low visibility
  Windy - Strong winds

  Temperature "feels like" accounts for:
  - Wind chill (makes it feel colder)
  - Humidity (makes it feel warmer/muggier)
TEXT

TEMPERATURE_GUIDE = <<~TEXT
  Temperature Conversions:

  Formulas:
  - Celsius to Fahrenheit: (C × 9/5) + 32
  - Fahrenheit to Celsius: (F - 32) × 5/9
  - Celsius to Kelvin: C + 273.15
  - Kelvin to Celsius: K - 273.15

  Reference Points:
  - -40°C = -40°F (same!)
  - 0°C = 32°F = 273.15K (water freezes)
  - 20°C = 68°F = 293.15K (room temperature)
  - 37°C = 98.6°F = 310.15K (body temperature)
  - 100°C = 212°F = 373.15K (water boils)

  Quick estimates:
  - To convert C to F: double it and add 30
  - To convert F to C: subtract 30 and halve it
TEXT

# ============================================================================
# Create and start server
# ============================================================================

server = MCP::Server.new(
  name: 'weather-mcp-ruby',
  version: '1.0.0',
  tools: [
    GetCurrentWeatherTool,
    GetForecastTool,
    CompareWeatherTool,
    ConvertTemperatureTool,
    CheckRainProbabilityTool
  ],
  resources: [
    MCP::Resource.new(
      uri: 'weather://supported-formats',
      name: 'supported-formats',
      description: 'List of supported location formats and examples',
      mime_type: 'text/plain'
    ),
    MCP::Resource.new(
      uri: 'weather://weather-codes',
      name: 'weather-codes',
      description: 'Explanation of weather condition codes',
      mime_type: 'text/plain'
    ),
    MCP::Resource.new(
      uri: 'weather://temperature-guide',
      name: 'temperature-guide',
      description: 'Temperature conversion formulas and reference points',
      mime_type: 'text/plain'
    )
  ]
)

# Register resource handler
server.resources_read_handler do |params|
  content = case params[:uri]
            when 'weather://supported-formats'
              SUPPORTED_FORMATS
            when 'weather://weather-codes'
              WEATHER_CODES
            when 'weather://temperature-guide'
              TEMPERATURE_GUIDE
            else
              'Resource not found'
            end

  [{
    uri: params[:uri],
    mimeType: 'text/plain',
    text: content
  }]
end

# Start stdio transport
transport = MCP::Server::Transports::StdioTransport.new(server)
transport.open
