package data

import (
	"errors"
	"net/http"

	"crossjoin.com/gorxestra/util/http/middlewares"
)

var (
	ErrInternalError        = errors.New("internal error")
	MusicAlreadyBeingPlayed = errors.New("music already being played")
)

type AppError struct {
	StatusCode   int
	ErrorMessage string
	ShowMessage  bool
}

var errorMap = map[error]middlewares.AppError{
	ErrInternalError: {
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: ErrInternalError.Error(),
		ShowMessage:  false,
	},
	MusicAlreadyBeingPlayed: {
		StatusCode:   http.StatusConflict,
		ErrorMessage: MusicAlreadyBeingPlayed.Error(),
		ShowMessage:  true,
	},
}

var strErrorMapper = map[string]error{
	ErrInternalError.Error(): ErrInternalError,
}

func MiddlewareErrorMap(err error) (middlewares.AppError, bool) {
	mappedError, ok := errorMap[err]
	if ok {
		return mappedError, true
	}

	return mappedError, false
}

func StrErrorMapper(errStr string) (error, bool) {
	err, ok := strErrorMapper[errStr]
	return err, ok
}
