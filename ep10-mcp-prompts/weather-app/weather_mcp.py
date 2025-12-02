from typing import Literal

import requests
from mcp.server.fastmcp import FastMCP

mcp = FastMCP("weather-mcp")


def fetch_weather(location: str) -> dict:
    """Internal helper to fetch from wttr.in"""
    url = f"https://wttr.in/{location}?format=j1"
    response = requests.get(url)
    response.raise_for_status()
    return response.json()


@mcp.tool()
def get_current_weather(location: str) -> str:
    """Get current weather conditions for a location

    Args:
        location: City name, airport code, or coordinates
    """
    data = fetch_weather(location)
    current = data["current_condition"][0]
    return f"""Weather in {location}:
Temperature: {current["temp_C"]}°C ({current["temp_F"]}°F)
Conditions: {current["weatherDesc"][0]["value"]}
Feels like: {current["FeelsLikeC"]}°C
Humidity: {current["humidity"]}%
Wind: {current["windspeedKmph"]} km/h"""


@mcp.tool()
def get_forecast(location: str, days: int = 3) -> str:
    """Get weather forecast for upcoming days

    Args:
        location: City name, airport code, or coordinates
        days: Number of days (1-3)
    """
    data = fetch_weather(location)
    forecast = data["weather"][:days]
    result = f"Forecast for {location}:\n\n"
    for day in forecast:
        result += f"{day['date']}:\n"
        result += f"  High: {day['maxtempC']}°C, Low: {day['mintempC']}°C\n"
        result += f"  Conditions: {day['hourly'][4]['weatherDesc'][0]['value']}\n"
        result += f"  Precipitation: {day['hourly'][4]['precipMM']}mm\n\n"
    return result


@mcp.tool()
def compare_weather(location1: str, location2: str) -> str:
    """Compare current weather between two locations

    Args:
        location1: First location
        location2: Second location
    """
    data1 = fetch_weather(location1)
    data2 = fetch_weather(location2)

    temp1 = int(data1["current_condition"][0]["temp_C"])
    temp2 = int(data2["current_condition"][0]["temp_C"])

    cond1 = data1["current_condition"][0]["weatherDesc"][0]["value"]
    cond2 = data2["current_condition"][0]["weatherDesc"][0]["value"]

    comparison = (
        "warmer"
        if temp1 > temp2
        else "cooler"
        if temp1 < temp2
        else "the same temperature as"
    )

    return f"""{location1}: {temp1}°C, {cond1}
{location2}: {temp2}°C, {cond2}

{location1} is {comparison} {location2} ({abs(temp1 - temp2)}°C difference)"""


@mcp.tool()
def convert_temperature(
    temp: float, from_unit: Literal["C", "F", "K"], to_unit: Literal["C", "F", "K"]
) -> str:
    """Convert temperature between Celsius, Fahrenheit, and Kelvin

    Args:
        temp: Temperature value
        from_unit: Source unit (C, F, or K)
        to_unit: Target unit (C, F, or K)
    """
    # Convert to Celsius first
    if from_unit == "F":
        celsius = (temp - 32) * 5 / 9
    elif from_unit == "K":
        celsius = temp - 273.15
    else:
        celsius = temp

    # Convert from Celsius to target
    if to_unit == "F":
        result = celsius * 9 / 5 + 32
    elif to_unit == "K":
        result = celsius + 273.15
    else:
        result = celsius

    return f"{temp}°{from_unit} = {result:.1f}°{to_unit}"


@mcp.tool()
def check_rain_probability(location: str) -> str:
    """Check if it's likely to rain today

    Args:
        location: City name or coordinates
    """
    data = fetch_weather(location)
    today = data["weather"][0]

    # Average precipitation chance across the day
    precip_chances = [int(hour["chanceofrain"]) for hour in today["hourly"]]
    avg_chance = sum(precip_chances) / len(precip_chances)

    if avg_chance > 70:
        level = "High"
        advice = "Bring an umbrella!"
    elif avg_chance > 40:
        level = "Moderate"
        advice = "Maybe bring an umbrella"
    else:
        level = "Low"
        advice = "You should be fine without an umbrella"

    return f"{level} chance of rain in {location} today ({avg_chance:.0f}%)\n{advice}"


@mcp.resource("weather://supported-formats")
def get_supported_formats() -> str:
    """List of supported location formats and examples"""
    return """Supported location formats:

1. City name: London, Paris, Tokyo, "New York"
2. Airport code: JFK, LAX, LHR
3. Coordinates: 51.5074,-0.1278 (latitude,longitude)
4. US location: NewYork,NY or Seattle,WA

Examples:
- get_current_weather("London")
- get_current_weather("JFK")
- get_current_weather("51.5074,-0.1278")
- get_current_weather("Seattle,WA")
"""


@mcp.resource("weather://weather-codes")
def get_weather_codes() -> str:
    """Explanation of weather condition codes"""
    return """Common weather conditions:

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
"""


@mcp.resource("weather://temperature-guide")
def get_temperature_guide() -> str:
    """Temperature conversion formulas and reference points"""
    return """Temperature Conversions:

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
"""


@mcp.prompt()
def weather_planning_assistant():
    """Help plan activities based on weather forecast

    Analyzes 3-day forecast, checks rain probability, and suggests
    appropriate activities and what to bring.
    """
    return """You are a weather planning assistant. When given a location:
1. Get the 3-day forecast
2. Check rain probability
3. Suggest outdoor vs indoor activities
4. Recommend what to pack (umbrella, sunscreen, jacket)
Be specific and practical."""


@mcp.prompt()
def travel_weather_comparison():
    """Compare weather between two cities for travel planning

    Uses compare_weather and provides travel recommendations.
    """
    return """You are a travel weather advisor. When comparing two locations:
1. Use compare_weather to get current conditions
2. Get forecasts for both cities
3. Recommend which has better weather for the next 3 days
4. Mention any weather alerts or concerns
Be concise but helpful."""


@mcp.prompt()
def temperature_context():
    """Explain temperatures in multiple units with context

    Converts temperatures and provides real-world context.
    """
    return """When given a temperature:
1. Convert it to Celsius, Fahrenheit, and Kelvin
2. Reference the temperature guide resource
3. Provide context (e.g., "warmer than room temperature", "freezing point")
4. Suggest appropriate clothing
Use the weather://temperature-guide resource."""


if __name__ == "__main__":
    mcp.run()
