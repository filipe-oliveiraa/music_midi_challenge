package http

import (
	"net/http"

	"crossjoin.com/gorxestra/logging"
	"github.com/labstack/echo/v4"
)

const (
	ContentTypeJson = "application/json"
)

type ErrorFunc func(status int, err error) error

func SetContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

// NodeInterface defines the node's methods required by the common APIs
type NodeInterface interface {
	Status() error
}

// HandlerFunc defines a wrapper for http.HandlerFunc that includes a context
type HandlerFunc func(ReqContext, echo.Context)

// Route type description
type Route struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc HandlerFunc
}

// Routes contains all routes
type Routes []Route

// ReqContext is passed to each of the handlers below via wrapCtx, allowing
// handlers to interact with the node
type ReqContext struct {
	Node     NodeInterface
	Log      logging.Logger
	Shutdown <-chan struct{}
}

// wrapCtx passes a common context to each request without a global variable.
func WrapCtx(ctx ReqContext, handler func(ReqContext, echo.Context)) echo.HandlerFunc {
	return func(context echo.Context) error {
		handler(ctx, context)
		return nil
	}
}

// registerHandler registers a set of Routes to the given router.
func RegisterHandlers(
	router *echo.Echo,
	prefix string,
	routes Routes,
	ctx ReqContext,
	m ...echo.MiddlewareFunc,
) {
	for _, route := range routes {
		r := router.Add(route.Method, prefix+route.Path, WrapCtx(ctx, route.HandlerFunc), m...)
		r.Name = route.Name
	}
}

func WrapHandler(f func(w http.ResponseWriter, r *http.Request)) echo.HandlerFunc {
	return func(c echo.Context) error {
		f(c.Response(), c.Request())
		return nil
	}
}
