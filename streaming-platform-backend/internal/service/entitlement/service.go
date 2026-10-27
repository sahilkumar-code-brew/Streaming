package entitlement

import (
	"context"
	"time"

	domain "github.com/code-brew-ai/streaming-platform-backend/internal/domain/entitlement"
	"github.com/code-brew-ai/streaming-platform-backend/internal/platform/apperror"
)

type Clock func() time.Time

type Settings struct {
	OriginBaseURL       string
	SteeringURL         string
	TelemetrySampleRate float64
	RefreshLead         time.Duration
}

type Service struct {
	repository domain.Repository
	signer     TokenSigner
	clock      Clock
	settings   Settings
}

func NewService(repository domain.Repository, signer TokenSigner, clock Clock, settings Settings) *Service {
	return &Service{repository: repository, signer: signer, clock: clock, settings: settings}
}

func (s *Service) Grant(ctx context.Context, request domain.Request) (domain.Grant, error) {
	snapshot, err := s.repository.Load(ctx, request.UserID, request.StreamID)
	if err != nil {
		return domain.Grant{}, apperror.Wrap(apperror.CodeServiceUnavailable, "entitlement store unavailable", err)
	}
	now := s.clock()
	if err := checkEligibility(snapshot, now); err != nil {
		return domain.Grant{}, err
	}
	token, expiresAt, err := s.signer.Sign(request, now)
	if err != nil {
		return domain.Grant{}, apperror.Wrap(apperror.CodeInternal, "token signing failed", err)
	}
	return s.buildGrant(request.StreamID, token, expiresAt, now), nil
}

func checkEligibility(snapshot domain.Snapshot, now time.Time) error {
	if snapshot.Stream == nil {
		return apperror.New(apperror.CodeNotFound, "stream not found")
	}
	if !snapshot.Stream.Watchable() {
		return apperror.New(apperror.CodeStreamNotLive, "stream is not live")
	}
	if snapshot.Subscription == nil {
		return apperror.New(apperror.CodeServiceUnavailable, "subscription data unavailable")
	}
	if !snapshot.Subscription.ActiveAt(now) {
		return apperror.New(apperror.CodeNotSubscribed, "no active subscription")
	}
	return nil
}

func (s *Service) buildGrant(streamID string, token string, expiresAt time.Time, now time.Time) domain.Grant {
	return domain.Grant{
		PlaybackToken:       token,
		TokenExpiresAt:      expiresAt,
		RefreshAfter:        expiresAt.Sub(now) - s.settings.RefreshLead,
		ManifestURL:         s.settings.OriginBaseURL + "/live/" + streamID + "/master.m3u8",
		SteeringManifestURL: s.settings.SteeringURL,
		TelemetrySampleRate: s.settings.TelemetrySampleRate,
	}
}
