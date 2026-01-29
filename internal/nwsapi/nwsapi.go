package nwsapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const userAgent = "(github.com/scacner/weather-api)"

// PointResponse is a simplified subset of the NWS API /points endpoint response
type PointResponse struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

// ForecastResponse is a simplified subset of the NWS API /forecast endpoint response
type ForecastResponse struct {
	Properties struct {
		Periods []struct {
			Temperature   float64 `json:"temperature"`
			ShortForecast string  `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

// getNWSForecastURL retrieves the forecast URL for given latitude and longitude
func getNWSForecastURL(lat, lon float64) (string, error) {
	// Format coordinates (NWS supports a max of 4 decimal places of precision in coordinates)
	latStr := strconv.FormatFloat(lat, 'f', 4, 64)
	lonStr := strconv.FormatFloat(lon, 'f', 4, 64)

	pointsURL := fmt.Sprintf("https://api.weather.gov/points/%s,%s", latStr, lonStr)

	req, err := http.NewRequest("GET", pointsURL, nil)
	if err != nil {
		return "", err
	}
	// NWS requires a User-Agent
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/geo+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("points request failed: %s %s", resp.Status, string(body))
	}

	var point PointResponse
	if err := json.NewDecoder(resp.Body).Decode(&point); err != nil {
		return "", err
	}

	if point.Properties.Forecast == "" {
		return "", fmt.Errorf("no forecast URL found in points response")
	}

	// Add units=us parameter for Fahrenheit
	forecastURL := fmt.Sprintf("%s?units= us", point.Properties.Forecast)

	return forecastURL, nil
}

// GetNWSForecast retrieves the weather forecast for given latitude and longitude
func GetNWSForecast(lat, lon float64) (*ForecastResponse, error) {
	forecastURL, err := getNWSForecastURL(lat, lon)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", forecastURL, nil)
	if err != nil {
		return nil, err
	}
	// NWS requires a User-Agent
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/geo+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("forecast request failed: %s %s", resp.Status, string(body))
	}

	var forecast ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}
