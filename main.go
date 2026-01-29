package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/scacner/weather-api/internal/nwsapi"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

// getTemperatureCategory categorizes temperature into "cold" (<40F), "moderate" (40-80F), or "hot" (>=80F).
func getTemperatureCategory(temp float64) (string, error) {
	switch {
	case temp < 40:
		return "cold", nil
	case temp >= 40 && temp < 80:
		return "moderate", nil
	case temp >= 80:
		return "hot", nil
	default:
		return "nil", errors.New("unable to determine temperature category")
	}
}

// getCurrentWeather defines a use case interactor to get current weather information.
func getCurrentWeather() usecase.Interactor {
	// Declare input type
	type currentInput struct {
		Latitude  float64 `query:"latitude" description:"Latitude of the coordinate being queried. Positive values are north of the equator, negative values are south."`
		Longitude float64 `query:"longitude" description:"Longitude of the coordinate being queried. Positive values are east of the Prime Meridian, negative values are west."`
	}

	// Declare output type
	type currentOutput struct {
		CurrentShortForecast    string `json:"current_short_forecast"`
		CurrentTemperatureRange string `json:"current_temperature_range"`
	}

	// Create use case interactor with references to input/output types and interaction function.
	u := usecase.NewInteractor(func(ctx context.Context, input currentInput, output *currentOutput) error {
		// Example: Rochester, NY approx coords (43.1567, -77.6148)
		forecast, err := nwsapi.GetNWSForecast(input.Latitude, input.Longitude)
		if err != nil {
			statusCode := status.Internal

			if errors.Is(err, nwsapi.NotFoundErr) {
				statusCode = status.NotFound
			}

			return status.Wrap(err, statusCode)
		}

		// Print today's forecast as example
		if len(forecast.Properties.Periods) > 0 {
			period := forecast.Properties.Periods[0]

			tempRange, err := getTemperatureCategory(period.Temperature)
			if err != nil {
				return status.Wrap(err, status.Internal)
			}

			output.CurrentShortForecast = period.ShortForecast
			output.CurrentTemperatureRange = tempRange
		}

		return nil
	})

	// Describe use case interactor
	u.SetTitle("Current Weather")
	u.SetDescription("Get current weather information for given latitude and longitude.")

	u.SetExpectedErrors(status.NotFound, status.Internal)

	return u
}

// Main function to set up and start the Weather API service.
func main() {
	s := web.NewService(openapi3.NewReflector())

	// Init API documentation schema
	s.OpenAPISchema().SetTitle("Weather API")
	s.OpenAPISchema().SetDescription("Uses the National Weather Service API Web Service to provide weather information.")
	s.OpenAPISchema().SetVersion("v0.0.1")

	// Add use case handler to router
	s.Get("/current", getCurrentWeather())

	// Swagger UI endpoint at /docs
	s.Docs("/docs", swgui.New)

	// Start server
	log.Println("API Spec: http://localhost:8080/docs")
	if err := http.ListenAndServe("localhost:8080", s); err != nil {
		log.Fatal(err)
	}
}
