package apperror

import (
	"errors"
	"net/http"
)

type Code string

const (
	CodeValidation         Code = "VALIDATION_ERROR"
	CodeUnauthorized       Code = "UNAUTHORIZED"
	CodeSessionExpired     Code = "SESSION_EXPIRED"
	CodeNotSubscribed      Code = "ENTITLEMENT_NOT_SUBSCRIBED"
	CodeNotFound           Code = "NOT_FOUND"
	CodeStreamNotLive      Code = "STREAM_NOT_LIVE"
	CodeInternal           Code = "INTERNAL_ERROR"
	CodeServiceUnavailable Code = "SERVICE_UNAVAILABLE"
	CodeRateLimited        Code = "RATE_LIMITED"
)

var statusByCode = map[Code]int{
	CodeValidation:         http.StatusBadRequest,
	CodeUnauthorized:       http.StatusUnauthorized,
	CodeSessionExpired:     http.StatusUnauthorized,
	CodeNotSubscribed:      http.StatusForbidden,
	CodeNotFound:           http.StatusNotFound,
	CodeStreamNotLive:      http.StatusConflict,
	CodeInternal:           http.StatusInternalServerError,
	CodeServiceUnavailable: http.StatusServiceUnavailable,
	CodeRateLimited:        http.StatusTooManyRequests,
}

var retryableCodes = map[Code]bool{
	CodeInternal:           true,
	CodeServiceUnavailable: true,
	CodeRateLimited:        true,
}

type Error struct {
	Code         Code
	Message      string
	Cause        error
	RetryAfterMs int64
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return string(e.Code) + ": " + e.Message
	}
	return string(e.Code) + ": " + e.Message + ": " + e.Cause.Error()
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(code Code, message string, cause error) *Error {
	return &Error{Code: code, Message: message, Cause: cause}
}

func NewRetryAfter(code Code, message string, retryAfterMs int64) *Error {
	return &Error{Code: code, Message: message, RetryAfterMs: retryAfterMs}
}

func Status(code Code) int {
	status, known := statusByCode[code]
	if !known {
		return http.StatusInternalServerError
	}
	return status
}

func IsRetryable(code Code) bool {
	return retryableCodes[code]
}

func From(err error) *Error {
	if appErr, ok := errors.AsType[*Error](err); ok {
		return appErr
	}
	return Wrap(CodeInternal, "unexpected error", err)
}
