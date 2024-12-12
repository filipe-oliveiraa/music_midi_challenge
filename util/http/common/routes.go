package common

import (
	"crossjoin.com/gorxestra/util/http"
)

// Routes are routes that are common for all versions
func (c *CommonApi) Routes() http.Routes {
	return http.Routes{
		http.Route{
			Method:      "OPTIONS",
			HandlerFunc: c.optionsHandler,
			Path:        "",
			Name:        "",
		},

		http.Route{
			Name:        "info",
			Method:      "GET",
			Path:        "/info",
			HandlerFunc: c.InfoHandler,
		},

		http.Route{
			Name:        "healthcheck",
			Method:      "GET",
			Path:        "/health",
			HandlerFunc: c.HealthCheck,
		},

		http.Route{
			Name:        "ready",
			Method:      "GET",
			Path:        "/ready",
			HandlerFunc: c.Ready,
		},

		http.Route{
			Name:        "swagger",
			Method:      "GET",
			Path:        "/swagger",
			HandlerFunc: c.Swagger,
		},

		http.Route{
			Name:        "startup",
			Method:      "GET",
			Path:        "/startup",
			HandlerFunc: c.StartupCheck,
		},

		http.Route{
			Name:        "metrics",
			Method:      "GET",
			Path:        "/metrics",
			HandlerFunc: c.Metrics,
		},
	}
}
