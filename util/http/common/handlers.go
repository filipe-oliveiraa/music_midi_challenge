package common

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"crossjoin.com/gorxestra/util/conf"
	lib "crossjoin.com/gorxestra/util/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CommonApi struct {
	ApiVersions     []string
	SwaggerSpecYAML string
}

// SwaggerJSON is an httpHandler for route GET /swagger.json
func (c *CommonApi) Swagger(ctx lib.ReqContext, context echo.Context) {
	w := context.Response().Writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(c.SwaggerSpecYAML))
}

// HealthCheck is an httpHandler for route GET /health
func (c *CommonApi) HealthCheck(ctx lib.ReqContext, context echo.Context) {
	w := context.Response().Writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(nil)
}

// StartupCheck is an httpHandler for route GET /startup
func (c *CommonApi) StartupCheck(ctx lib.ReqContext, context echo.Context) {
	w := context.Response().Writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(nil)
}

// Ready is a httpHandler for route GET /ready
func (c *CommonApi) Ready(ctx lib.ReqContext, context echo.Context) {
	w := context.Response().Writer
	w.Header().Set("Content-Type", "application/json")

	err := ctx.Node.Status()
	code := http.StatusOK

	if err != nil {
		code = http.StatusInternalServerError
		ctx.Log.Error(err)
	}

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(nil)
}

// InfoHandler is an httpHandler for route GET /info
func (c *CommonApi) InfoHandler(ctx lib.ReqContext, context echo.Context) {
	w := context.Response().Writer

	currentVersion := conf.GetCurrentVersion()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := InfoResponse{
		Body: Info{
			Versions: c.ApiVersions,
			Build: BuildVersion{
				Major:       currentVersion.Major,
				Minor:       currentVersion.Minor,
				BuildNumber: currentVersion.BuildNumber,
				CommitHash:  currentVersion.CommitHash,
				Branch:      currentVersion.Branch,
				Channel:     currentVersion.Channel,
			},
		},
	}

	_ = json.NewEncoder(w).Encode(response.Body)
}

// Metrics returns data collected by prometheus
func (c *CommonApi) Metrics(ctx lib.ReqContext, context echo.Context) {
	promhttp.Handler().ServeHTTP(context.Response().Writer, context.Request())
}

// CORS
func (c *CommonApi) optionsHandler(ctx lib.ReqContext, context echo.Context) {
	context.Response().Writer.WriteHeader(http.StatusOK)
}
