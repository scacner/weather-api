# Weather API
An HTTP server that serves the current weather for a provided set of latitude and longitude coordinates. Uses the [National Weather Service API Web Service](https://www.weather.gov/documentation/services-web-api) as a data source. The server returns the current short forecast for the provided area for Today (e.g. “Partly Cloudy”) and a characterization of whether the temperature is cold (less than 40F), "moderate" (at or above 40F and less than 80F), or "hot" (at or above 80F)

This HTTP server is built using the [Swaggest REST](https://github.com/swaggest/rest) library to provide OpenAPI documentation for the API. I enjoy using Swaggest for building APIs in Go as it provides a great developer experience and makes it easy to keep API documentation in sync with the code.

When running this app, an API Spec can be viewed at http://localhost:8080/docs

## Running API Locally
To run locally, `go` needs to be installed on the local machine. The app can be ran with the following command:

```shell
go run main.go
```

A Makefile has also been provided. If the local machine supports the `make` command, simply run the following command instead:

```shell
make run
```

## Example Call
An example call to the API for the Rochester, NY area is shown below:

```shell
$ curl "http://localhost:8080/current?latitude=43.1567&longitude=-77.6148"
{"current_short_forecast":"Snow","current_temperature_range":"cold"}
```
