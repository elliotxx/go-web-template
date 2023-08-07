package errcode

import (
	"net/http"
	"strings"

	"github.com/elliotxx/errors"
)

func NewErrorCode(code string, message string) errors.ErrorCode {
	mustSetCodeIfNotPresent(code)
	return errors.NewErrorCode(code, message)
}

// Scope returns the error's scope
func Scope(err error) string {
	if err != nil {
		if e, ok := err.(errors.DetailError); ok && len(e.GetCode()) == 5 {
			return e.GetCode()[0:3]
		}
		if e, ok := err.(errors.ErrorCode); ok && len(e.GetCode()) == 5 {
			return e.GetCode()[0:3]
		}
		return InvalidScope
	}
	return ""
}

func InvalidErrorScope(code string) bool {
	c := strings.TrimSpace(code)
	return c == InvalidScope || c == ""
}

// StatusCode returns the http status code of the error
func StatusCode(err error) int {
	switch Scope(err) {
	case Scope(Success):
		return http.StatusOK
	case Scope(NotFound):
		return http.StatusNotFound
	case Scope(ServerError), Scope(InternalError):
		return http.StatusInternalServerError
	case Scope(InvalidParams):
		return http.StatusBadRequest
	case Scope(AccessPermissionError):
		return http.StatusUnauthorized
	case Scope(TooManyRequests):
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
