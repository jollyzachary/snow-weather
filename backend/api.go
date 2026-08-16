package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	forecastEndpoint   = "https://api.open-meteo.com/v1/forecast"
	airQualityEndpoint = "https://air-quality-api.open-meteo.com/v1/air-quality"
	geocodingEndpoint  = "https://geocoding-api.open-meteo.com/v1/search"
	maxUpstreamBody    = 4 << 20
)

type application struct {
	cache  *responseCache
	client *http.Client
	logger *slog.Logger
}

type weatherBundle struct {
	Weather     json.RawMessage `json:"weather"`
	AirQuality  json.RawMessage `json:"air_quality"`
	Cached      bool            `json:"cached"`
	GeneratedAt time.Time       `json:"generated_at"`
}

type upstreamResult struct {
	name string
	body []byte
	err  error
}

func newApplication(logger *slog.Logger) *application {
	return &application{
		cache: newResponseCache(256),
		client: &http.Client{
			Timeout: 8 * time.Second,
		},
		logger: logger,
	}
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", app.handleHealth)
	mux.HandleFunc("GET /api/locations", app.handleLocations)
	mux.HandleFunc("GET /api/weather", app.handleWeather)
	mux.Handle("/", spaHandler())
	return app.securityHeaders(app.requestLog(mux))
}

func (app *application) handleHealth(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "available"})
}

func (app *application) handleLocations(writer http.ResponseWriter, request *http.Request) {
	query := strings.TrimSpace(request.URL.Query().Get("q"))
	if len(query) < 2 || len(query) > 80 {
		writeError(writer, http.StatusBadRequest, "search must contain between 2 and 80 characters")
		return
	}

	cacheKey := "locations:" + strings.ToLower(query)
	if body, ok := app.cache.Get(cacheKey); ok {
		writeRawJSON(writer, http.StatusOK, body)
		return
	}

	parameters := url.Values{
		"name":     {query},
		"count":    {"8"},
		"language": {"en"},
		"format":   {"json"},
	}
	body, err := app.fetch(request.Context(), geocodingEndpoint, parameters)
	if err != nil {
		app.logger.Warn("location lookup failed", "error", err)
		writeError(writer, http.StatusBadGateway, "location service is temporarily unavailable")
		return
	}

	app.cache.Set(cacheKey, body, 12*time.Hour)
	writeRawJSON(writer, http.StatusOK, body)
}

func (app *application) handleWeather(writer http.ResponseWriter, request *http.Request) {
	latitude, longitude, err := coordinates(request.URL.Query())
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	units := request.URL.Query().Get("units")
	if units == "" {
		units = "imperial"
	}
	if units != "imperial" && units != "metric" {
		writeError(writer, http.StatusBadRequest, "units must be imperial or metric")
		return
	}

	cacheKey := fmt.Sprintf("weather:%.4f:%.4f:%s", latitude, longitude, units)
	if body, ok := app.cache.Get(cacheKey); ok {
		var bundle weatherBundle
		if json.Unmarshal(body, &bundle) == nil {
			bundle.Cached = true
			writeJSON(writer, http.StatusOK, bundle)
			return
		}
	}

	requestContext, cancel := context.WithTimeout(request.Context(), 9*time.Second)
	defer cancel()

	results := make(chan upstreamResult, 2)
	go app.fetchForecast(requestContext, latitude, longitude, units, results)
	go app.fetchAirQuality(requestContext, latitude, longitude, results)

	var weatherBody, airQualityBody []byte
	for range 2 {
		select {
		case <-requestContext.Done():
			writeError(writer, http.StatusGatewayTimeout, "weather service timed out")
			return
		case result := <-results:
			if result.err != nil {
				app.logger.Warn("weather upstream failed", "source", result.name, "error", result.err)
				if result.name == "forecast" {
					writeError(writer, http.StatusBadGateway, "weather data is temporarily unavailable")
					return
				}
				airQualityBody = []byte(`{}`)
				continue
			}
			if result.name == "forecast" {
				weatherBody = result.body
			} else {
				airQualityBody = result.body
			}
		}
	}

	bundle := weatherBundle{
		Weather:     weatherBody,
		AirQuality:  airQualityBody,
		Cached:      false,
		GeneratedAt: time.Now().UTC(),
	}
	payload, err := json.Marshal(bundle)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "could not prepare weather response")
		return
	}
	app.cache.Set(cacheKey, payload, 10*time.Minute)
	writeRawJSON(writer, http.StatusOK, payload)
}

func (app *application) fetchForecast(
	ctx context.Context,
	latitude float64,
	longitude float64,
	units string,
	results chan<- upstreamResult,
) {
	parameters := url.Values{
		"latitude":      {strconv.FormatFloat(latitude, 'f', 5, 64)},
		"longitude":     {strconv.FormatFloat(longitude, 'f', 5, 64)},
		"timezone":      {"auto"},
		"forecast_days": {"10"},
		"current": {strings.Join([]string{
			"temperature_2m", "relative_humidity_2m", "apparent_temperature",
			"is_day", "precipitation", "rain", "showers", "snowfall", "weather_code",
			"cloud_cover", "surface_pressure", "wind_speed_10m", "wind_direction_10m",
			"wind_gusts_10m",
		}, ",")},
		"hourly": {strings.Join([]string{
			"temperature_2m", "apparent_temperature", "precipitation_probability",
			"weather_code", "wind_speed_10m", "uv_index", "visibility",
		}, ",")},
		"daily": {strings.Join([]string{
			"weather_code", "temperature_2m_max", "temperature_2m_min", "sunrise", "sunset",
			"uv_index_max", "precipitation_probability_max", "wind_speed_10m_max",
		}, ",")},
	}
	if units == "imperial" {
		parameters.Set("temperature_unit", "fahrenheit")
		parameters.Set("wind_speed_unit", "mph")
		parameters.Set("precipitation_unit", "inch")
	} else {
		parameters.Set("temperature_unit", "celsius")
		parameters.Set("wind_speed_unit", "kmh")
		parameters.Set("precipitation_unit", "mm")
	}

	body, err := app.fetch(ctx, forecastEndpoint, parameters)
	results <- upstreamResult{name: "forecast", body: body, err: err}
}

func (app *application) fetchAirQuality(
	ctx context.Context,
	latitude float64,
	longitude float64,
	results chan<- upstreamResult,
) {
	parameters := url.Values{
		"latitude":  {strconv.FormatFloat(latitude, 'f', 5, 64)},
		"longitude": {strconv.FormatFloat(longitude, 'f', 5, 64)},
		"timezone":  {"auto"},
		"current":   {"us_aqi,pm2_5"},
	}
	body, err := app.fetch(ctx, airQualityEndpoint, parameters)
	results <- upstreamResult{name: "air_quality", body: body, err: err}
}

func (app *application) fetch(
	ctx context.Context,
	endpoint string,
	parameters url.Values,
) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+parameters.Encode(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "snow-weather/2.0")

	response, err := app.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("upstream returned status %d", response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxUpstreamBody))
	if err != nil {
		return nil, err
	}
	if !json.Valid(body) {
		return nil, errors.New("upstream returned invalid JSON")
	}
	return body, nil
}

func coordinates(values url.Values) (float64, float64, error) {
	latitude, err := strconv.ParseFloat(values.Get("lat"), 64)
	if err != nil || latitude < -90 || latitude > 90 {
		return 0, 0, errors.New("latitude must be between -90 and 90")
	}
	longitude, err := strconv.ParseFloat(values.Get("lon"), 64)
	if err != nil || longitude < -180 || longitude > 180 {
		return 0, 0, errors.New("longitude must be between -180 and 180")
	}
	return latitude, longitude, nil
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func writeRawJSON(writer http.ResponseWriter, status int, payload []byte) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_, _ = writer.Write(payload)
}

func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}

func (app *application) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("X-Frame-Options", "DENY")
		writer.Header().Set("Permissions-Policy", "camera=(), microphone=(), payment=()")
		next.ServeHTTP(writer, request)
	})
}

func (app *application) requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		started := time.Now()
		next.ServeHTTP(writer, request)
		app.logger.Info(
			"request completed",
			"method", request.Method,
			"path", request.URL.Path,
			"duration_ms", time.Since(started).Milliseconds(),
		)
	})
}
