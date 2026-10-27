package entitlement

import "time"

type StreamState string

const (
	StreamScheduled   StreamState = "scheduled"
	StreamLive        StreamState = "live"
	StreamInterrupted StreamState = "interrupted"
	StreamEnded       StreamState = "ended"
)

const SubscriptionActive = "active"

type Subscription struct {
	Status    string
	ExpiresAt time.Time
}

func (s Subscription) ActiveAt(now time.Time) bool {
	return s.Status == SubscriptionActive && now.Before(s.ExpiresAt)
}

type Stream struct {
	ID    string
	State StreamState
}

func (s Stream) Watchable() bool {
	return s.State == StreamLive || s.State == StreamInterrupted
}

type Snapshot struct {
	Subscription *Subscription
	Stream       *Stream
}

type Request struct {
	UserID   string
	StreamID string
	DeviceID string
}

type Grant struct {
	PlaybackToken       string
	TokenExpiresAt      time.Time
	RefreshAfter        time.Duration
	ManifestURL         string
	SteeringManifestURL string
	TelemetrySampleRate float64
}
