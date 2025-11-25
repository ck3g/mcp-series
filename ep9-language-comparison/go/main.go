package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ============================================================================
// Input types for tools
// ============================================================================

type LocationInput struct {
	Location string `json:"location" jsonschema:"City name, airport code, or coordinates"`
}

type ForecastInput struct {
	Location string `json:"location" jsonschema:"City name, airport code, or coordinates"`
	Days     int    `json:"days,omitempty" jsonschema:"Number of days for forecast (1-3)"`
}

type CompareInput struct {
	Location1 string `json:"location1" jsonschema:"First location to compare"`
	Location2 string `json:"location2" jsonschema:"Second location to compare"`
}

type ConvertInput struct {
	Temp     float64 `json:"temp" jsonschema:"Temperature value to convert"`
	FromUnit string  `json:"from_unit" jsonschema:"Source unit: C, F, or K"`
	ToUnit   string  `json:"to_unit" jsonschema:"Target unit: C, F, or K"`
}

// ============================================================================
// Output types for tools
// ============================================================================

type WeatherOutput struct {
	Result string `json:"result" jsonschema:"Weather information"`
}

type ForecastOutput struct {
	Result string `json:"result" jsonschema:"Weather forecast"`
}

type CompareOutput struct {
	Result string `json:"result" jsonschema:"Weather comparison"`
}

type ConvertOutput struct {
	Result string `json:"result" jsonschema:"Temperature conversion result"`
}

type RainOutput struct {
	Result string `json:"result" jsonschema:"Rain probability information"`
}

// ============================================================================
// Weather API response structures
// ============================================================================

type WeatherData struct {
	CurrentCondition []struct {
		TempC         string `json:"temp_C"`
		TempF         string `json:"temp_F"`
		FeelsLikeC    string `json:"FeelsLikeC"`
		Humidity      string `json:"humidity"`
		WindspeedKmph string `json:"windspeedKmph"`
		WeatherDesc   []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
	} `json:"current_condition"`
	Weather []struct {
		Date     string `json:"date"`
		MaxtempC string `json:"maxtempC"`
		MintempC string `json:"mintempC"`
		Hourly   []struct {
			WeatherDesc []struct {
				Value string `json:"value"`
			} `json:"weatherDesc"`
			PrecipMM     string `json:"precipMM"`
			ChanceOfRain string `json:"chanceofrain"`
		} `json:"hourly"`
	} `json:"weather"`
}

// ============================================================================
// Helper function to fetch weather data
// ============================================================================

func fetchWeather(location string) (*WeatherData, error) {
	url := fmt.Sprintf("https://wttr.in/%s?format=j1", location)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data WeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// ============================================================================
// Tool 1: Get Current Weather
// ============================================================================

func GetCurrentWeather(ctx context.Context, req *mcp.CallToolRequest, input LocationInput) (*mcp.CallToolResult, WeatherOutput, error) {
	data, err := fetchWeather(input.Location)
	if err != nil {
		return nil, WeatherOutput{}, err
	}

	current := data.CurrentCondition[0]
	result := fmt.Sprintf(`Weather in %s:
Temperature: %s°C (%s°F)
Conditions: %s
Feels like: %s°C
Humidity: %s%%
Wind: %s km/h`,
		input.Location,
		current.TempC,
		current.TempF,
		current.WeatherDesc[0].Value,
		current.FeelsLikeC,
		current.Humidity,
		current.WindspeedKmph,
	)

	return nil, WeatherOutput{Result: result}, nil
}

// ============================================================================
// Tool 2: Get Forecast
// ============================================================================

func GetForecast(ctx context.Context, req *mcp.CallToolRequest, input ForecastInput) (*mcp.CallToolResult, ForecastOutput, error) {
	if input.Days == 0 {
		input.Days = 3
	}
	if input.Days > 3 {
		input.Days = 3
	}

	data, err := fetchWeather(input.Location)
	if err != nil {
		return nil, ForecastOutput{}, err
	}

	result := fmt.Sprintf("Forecast for %s:\n\n", input.Location)

	maxDays := min(input.Days, len(data.Weather))

	for i := range maxDays {
		day := data.Weather[i]
		result += fmt.Sprintf("%s:\n", day.Date)
		result += fmt.Sprintf("  High: %s°C, Low: %s°C\n", day.MaxtempC, day.MintempC)

		if len(day.Hourly) > 4 {
			result += fmt.Sprintf("  Conditions: %s\n", day.Hourly[4].WeatherDesc[0].Value)
			result += fmt.Sprintf("  Precipitation: %smm\n\n", day.Hourly[4].PrecipMM)
		}
	}

	return nil, ForecastOutput{Result: result}, nil
}

// ============================================================================
// Tool 3: Compare Weather
// ============================================================================

func CompareWeather(ctx context.Context, req *mcp.CallToolRequest, input CompareInput) (*mcp.CallToolResult, CompareOutput, error) {
	data1, err := fetchWeather(input.Location1)
	if err != nil {
		return nil, CompareOutput{}, err
	}

	data2, err := fetchWeather(input.Location2)
	if err != nil {
		return nil, CompareOutput{}, err
	}

	temp1 := data1.CurrentCondition[0].TempC
	temp2 := data2.CurrentCondition[0].TempC
	cond1 := data1.CurrentCondition[0].WeatherDesc[0].Value
	cond2 := data2.CurrentCondition[0].WeatherDesc[0].Value

	var temp1Int, temp2Int int
	fmt.Sscanf(temp1, "%d", &temp1Int)
	fmt.Sscanf(temp2, "%d", &temp2Int)

	diff := temp1Int - temp2Int
	comparison := "the same temperature as"
	if diff > 0 {
		comparison = "warmer"
	} else if diff < 0 {
		comparison = "cooler"
		diff = -diff
	} else {
		diff = 0
	}

	result := fmt.Sprintf(`%s: %s°C, %s
%s: %s°C, %s

%s is %s %s (%d°C difference)`,
		input.Location1, temp1, cond1,
		input.Location2, temp2, cond2,
		input.Location1, comparison, input.Location2, diff,
	)

	return nil, CompareOutput{Result: result}, nil
}

// ============================================================================
// Tool 4: Convert Temperature
// ============================================================================

func ConvertTemperature(ctx context.Context, req *mcp.CallToolRequest, input ConvertInput) (*mcp.CallToolResult, ConvertOutput, error) {
	// Convert to Celsius first
	var celsius float64
	switch input.FromUnit {
	case "F":
		celsius = (input.Temp - 32) * 5 / 9
	case "K":
		celsius = input.Temp - 273.15
	default:
		celsius = input.Temp
	}

	// Convert from Celsius to target
	var result float64
	switch input.ToUnit {
	case "F":
		result = celsius*9/5 + 32
	case "K":
		result = celsius + 273.15
	default:
		result = celsius
	}

	output := fmt.Sprintf("%.1f°%s = %.1f°%s", input.Temp, input.FromUnit, result, input.ToUnit)
	return nil, ConvertOutput{Result: output}, nil
}

// ============================================================================
// Tool 5: Check Rain Probability
// ============================================================================
func CheckRainProbability(ctx context.Context, req *mcp.CallToolRequest, input LocationInput) (*mcp.CallToolResult, RainOutput, error) {
	data, err := fetchWeather(input.Location)
	if err != nil {
		return nil, RainOutput{}, err
	}

	if len(data.Weather) == 0 {
		return nil, RainOutput{}, fmt.Errorf("no weather data available")
	}

	today := data.Weather[0]

	// Calculate average rain chance
	var total int
	count := 0
	for _, hour := range today.Hourly {
		var chance int
		fmt.Sscanf(hour.ChanceOfRain, "%d", &chance)
		total += chance
		count++
	}

	avgChance := 0
	if count > 0 {
		avgChance = total / count
	}

	level := "Low"
	advice := "You should be fine without an umbrella"

	if avgChance > 70 {
		level = "High"
		advice = "Bring an umbrella!"
	} else if avgChance > 40 {
		level = "Moderate"
		advice = "Maybe bring an umbrella"
	}

	result := fmt.Sprintf("%s chance of rain in %s today (%d%%)\n%s", level, input.Location, avgChance, advice)
	return nil, RainOutput{Result: result}, nil
}

// ============================================================================
// Resource handlers
// ============================================================================

func getSupportedFormats(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	content := `Supported location formats:

1. City name: London, Paris, Tokyo, "New York"
2. Airport code: JFK, LAX, LHR
3. Coordinates: 51.5074,-0.1278 (latitude,longitude)
4. US location: NewYork,NY or Seattle,WA

Examples:
- get_current_weather("London")
- get_current_weather("JFK")
- get_current_weather("51.5074,-0.1278")
- get_current_weather("Seattle,WA")`

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "weather://supported-formats",
				MIMEType: "text/plain",
				Text:     content,
			},
		},
	}, nil
}

func getWeatherCodes(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	content := `Common weather conditions:

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
- Humidity (makes it feel warmer/muggier)`

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "weather://weather-codes",
				MIMEType: "text/plain",
				Text:     content,
			},
		},
	}, nil
}

func getTemperatureGuide(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	content := `Temperature Conversions:

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
- To convert F to C: subtract 30 and halve it`

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "weather://temperature-guide",
				MIMEType: "text/plain",
				Text:     content,
			},
		},
	}, nil
}

// ============================================================================
// Main function
// ============================================================================

func main() {
	// Create server
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "weather-mcp-go",
			Version: "1.0.0",
		},
		nil,
	)

	// Register tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_current_weather",
		Description: "Get current weather conditions for a location",
	}, GetCurrentWeather)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_forecast",
		Description: "Get weather forecast for upcoming days",
	}, GetForecast)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "compare_weather",
		Description: "Compare current weather between two locations",
	}, CompareWeather)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "convert_temperature",
		Description: "Convert temperature between Celsius, Fahrenheit, and Kelvin",
	}, ConvertTemperature)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_rain_probability",
		Description: "Check if it's likely to rain today",
	}, CheckRainProbability)

	// Register resources
	server.AddResource(&mcp.Resource{
		URI:         "weather://supported-formats",
		Name:        "supported-formats",
		Description: "List of supported location formats and examples",
		MIMEType:    "text/plain",
	}, getSupportedFormats)

	server.AddResource(&mcp.Resource{
		URI:         "weather://weather-codes",
		Name:        "weather-codes",
		Description: "Explanation of weather condition codes",
		MIMEType:    "text/plain",
	}, getWeatherCodes)

	server.AddResource(&mcp.Resource{
		URI:         "weather://temperature-guide",
		Name:        "temperature-guide",
		Description: "Temperature conversion formulas and reference points",
		MIMEType:    "text/plain",
	}, getTemperatureGuide)

	// Run server over stdio
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
