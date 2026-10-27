package rest

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	domain "github.com/code-brew-ai/streaming-platform-backend/internal/domain/entitlement"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/apperror"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/httpx"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/session"
	"github.com/code-brew-ai/streaming-platform-backend/pkg/validator"
)

const (
	EntitlementRoutePattern = "POST /api/v1/streams/{stream_id}/entitlement"
	headerDeviceID          = "X-Device-Id"
	pathStreamID            = "stream_id"
	timestampLayout         = "2006-01-02T15:04:05.000Z07:00"
)

type Granter interface {
	Grant(ctx context.Context, request domain.Request) (domain.Grant, error)
}

type grantResponse struct {
	PlaybackToken       string  `json:"playbackToken"`
	TokenExpiresAt      string  `json:"tokenExpiresAt"`
	RefreshAfterS       int64   `json:"refreshAfterS"`
	ManifestURL         string  `json:"manifestUrl"`
	SteeringManifestURL string  `json:"steeringManifestUrl"`
	TelemetrySampleRate float64 `json:"telemetrySampleRate"`
}

type EntitlementHandler struct {
	granter Granter
	logger  *slog.Logger
}

func NewEntitlementHandler(granter Granter, logger *slog.Logger) *EntitlementHandler {
	return &EntitlementHandler{granter: granter, logger: logger}
}

func (h *EntitlementHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request, err := parseEntitlementRequest(r)
	if err != nil {
		httpx.WriteError(w, r, h.logger, err)
		return
	}
	grant, err := h.granter.Grant(r.Context(), request)
	if err != nil {
		httpx.WriteError(w, r, h.logger, err)
		return
	}
	httpx.WriteData(w, r, http.StatusOK, toResponse(grant))
}

func parseEntitlementRequest(r *http.Request) (domain.Request, error) {
	streamID := r.PathValue(pathStreamID)
	if err := validator.ValidateStreamID(streamID); err != nil {
		return domain.Request{}, apperror.New(apperror.CodeValidation, err.Error())
	}
	deviceID := r.Header.Get(headerDeviceID)
	if err := validator.ValidateDeviceID(deviceID); err != nil {
		return domain.Request{}, apperror.New(apperror.CodeValidation, err.Error())
	}
	return domain.Request{UserID: session.UserID(r.Context()), StreamID: streamID, DeviceID: deviceID}, nil
}

func toResponse(grant domain.Grant) grantResponse {
	return grantResponse{
		PlaybackToken:       grant.PlaybackToken,
		TokenExpiresAt:      grant.TokenExpiresAt.UTC().Format(timestampLayout),
		RefreshAfterS:       int64(grant.RefreshAfter / time.Second),
		ManifestURL:         grant.ManifestURL,
		SteeringManifestURL: grant.SteeringManifestURL,
		TelemetrySampleRate: grant.TelemetrySampleRate,
	}
}
