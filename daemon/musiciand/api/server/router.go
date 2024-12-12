// Package server Conductord REST API.
package server

import (
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"crossjoin.com/gorxestra/config"
	"crossjoin.com/gorxestra/daemon/musiciand/api"
	v1 "crossjoin.com/gorxestra/daemon/musiciand/api/server/v1"
	"crossjoin.com/gorxestra/daemon/musiciand/api/server/v1/openapi/generated/server"
	"crossjoin.com/gorxestra/data"
	httpUtils "crossjoin.com/gorxestra/util/http"
	"crossjoin.com/gorxestra/util/http/middlewares"
	"crossjoin.com/gorxestra/util/http/pprof"

	"crossjoin.com/gorxestra/util/http/common"

	openapi "crossjoin.com/gorxestra/daemon/musiciand/api/server/v1/openapi"
	"crossjoin.com/gorxestra/logging"
)

// APINodeInterface describes all the node methods required by common
type APINodeInterface interface {
	httpUtils.NodeInterface
	api.NodeInterface
	Config() config.MusicianConf
}

const (
	// TokenHeader is the header where we put the token.
	TokenHeader = "X-API-Token" //nolint: all
	// MaxRequestBodyBytes is the maximum request body size that we allow in our APIs.
	MaxRequestBodyBytes = "10MB"
)

// NewHttpRouter builds and returns a new router with our REST handlers registered.
func NewHttpRouter(
	logger logging.Logger,
	node APINodeInterface,
	shutdown <-chan struct{},
	listener net.Listener,
	numConnectionsLimit uint64,
) *echo.Echo {
	publicMiddleware := []echo.MiddlewareFunc{
		middleware.BodyLimit(MaxRequestBodyBytes),
	}

	e := echo.New()
	e.HidePort = true
	e.Listener = listener
	e.HideBanner = true

	e.Pre(
		middlewares.MakeConnectionLimiter(numConnectionsLimit),
		middleware.RemoveTrailingSlash())

	e.Use(
		middlewares.MakeLogger(logger),
		middlewares.MakeCORS(TokenHeader),
		middlewares.MakeError(logger, data.MiddlewareErrorMap),
	)

	// Request Context
	ctx := httpUtils.ReqContext{Node: node, Log: logger, Shutdown: shutdown}

	// Pprof
	pprof.Wrap(e)

	// Registering common routes (no auth)
	common := common.CommonApi{
		SwaggerSpecYAML: openapi.SwaggerSpecYAMLEmbed,
		ApiVersions:     []string{"v1"},
	}
	httpUtils.RegisterHandlers(e, "", common.Routes(), ctx)

	// register v1 handlers
	v1 := v1.Handlers{
		Node: node,
		Log:  logger,
	}

	server.RegisterHandlers(e, &v1, publicMiddleware...)

	return e
}
