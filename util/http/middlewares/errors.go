package middlewares

import (
	"net/http"

	log "crossjoin.com/gorxestra/logging"
	"crossjoin.com/gorxestra/util/errors"
	"crossjoin.com/gorxestra/util/http/common"
	"github.com/labstack/echo/v4"
)

const (
	InternalServerErrorMessage = "internal server error"
)

type AppError struct {
	StatusCode   int
	ErrorMessage string
	ShowMessage  bool
}

type (
	ErrorMapper func(err error) (AppError, bool)
)

var echoErrors = map[error]AppError{
	echo.ErrNotFound: {
		StatusCode:   http.StatusNotFound,
		ErrorMessage: "not found",
		ShowMessage:  true,
	},
	echo.ErrMethodNotAllowed: {
		StatusCode:   http.StatusMethodNotAllowed,
		ErrorMessage: "method not allowed",
		ShowMessage:  true,
	},
}

// LoggerMiddleware provides some extra state to the logger middleware
type ErrorMiddleware struct {
	log         log.Logger
	errorMapper ErrorMapper
}

// MakeLogger initializes the logger middleware function
func MakeError(log log.Logger, mapper ErrorMapper) echo.MiddlewareFunc {
	error := ErrorMiddleware{
		log:         log,
		errorMapper: toInternalErrorMapper(mapper),
	}
	return error.handler
}

func toInternalErrorMapper(mapper ErrorMapper) ErrorMapper {
	return func(err error) (AppError, bool) {
		if !errors.IsHashable(err) {
			//nolint
			return AppError{}, false
		}

		mappedErr, ok := echoErrors[err]
		if ok {
			return mappedErr, true
		}

		return mapper(err)
	}
}

// Logger is an echo middleware to add log to the API
func (logger *ErrorMiddleware) handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)

		if err == nil {
			return nil
		}

		mappedError, ok := logger.errorMapper(err)
		if ok {
			msg := "Error occurred"
			if mappedError.ShowMessage {
				msg = mappedError.ErrorMessage
			}
			return ctx.JSON(mappedError.StatusCode, common.Error{
				Error: msg,
			})
		}

		logger.log.ErrorStack(err.Error())

		return ctx.JSON(http.StatusInternalServerError, common.Error{
			Error: InternalServerErrorMessage,
		})
	}
}
