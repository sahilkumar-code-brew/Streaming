package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/apperror"
)

const timestampLayout = "2006-01-02T15:04:05.000Z07:00"

type Meta struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
}

type successBody[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
	Meta    Meta `json:"meta"`
}

type errorDetail struct {
	Code         apperror.Code `json:"code"`
	Message      string        `json:"message"`
	Retryable    bool          `json:"retryable"`
	RetryAfterMs int64         `json:"retryAfterMs,omitempty"`
}

type errorBody struct {
	Success bool        `json:"success"`
	Error   errorDetail `json:"error"`
	Meta    Meta        `json:"meta"`
}

func WriteData[T any](w http.ResponseWriter, r *http.Request, status int, data T) {
	writeJSON(w, status, successBody[T]{Success: true, Data: data, Meta: newMeta(r)})
}

func WriteError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	appErr := apperror.From(err)
	status := apperror.Status(appErr.Code)
	logFailure(r, logger, status, appErr)
	detail := errorDetail{Code: appErr.Code, Message: appErr.Message, Retryable: apperror.IsRetryable(appErr.Code)}
	writeJSON(w, status, errorBody{Success: false, Error: detail, Meta: newMeta(r)})
}

func newMeta(r *http.Request) Meta {
	return Meta{Timestamp: time.Now().UTC().Format(timestampLayout), RequestID: RequestID(r.Context())}
}

func logFailure(r *http.Request, logger *slog.Logger, status int, appErr *apperror.Error) {
	if status < http.StatusInternalServerError {
		return
	}
	logger.ErrorContext(r.Context(), "request failed", "requestId", RequestID(r.Context()), "code", appErr.Code, "error", appErr.Error())
}

func writeJSON[T any](w http.ResponseWriter, status int, body T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
