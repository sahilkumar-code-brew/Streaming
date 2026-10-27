package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	domain "github.com/code-brew-ai/streaming-platform-backend/internal/domain/entitlement"
)

const (
	subscriptionKeyPrefix = "sub:"
	streamMetaKeyPrefix   = "stream:meta:"
	fieldStatus           = "status"
	fieldExpiresAt        = "expires_at"
	fieldState            = "state"
)

type EntitlementRepository struct {
	client *goredis.Client
}

func NewEntitlementRepository(client *goredis.Client) *EntitlementRepository {
	return &EntitlementRepository{client: client}
}

func (r *EntitlementRepository) Load(ctx context.Context, userID string, streamID string) (domain.Snapshot, error) {
	pipe := r.client.Pipeline()
	subscriptionCmd := pipe.HGetAll(ctx, subscriptionKeyPrefix+userID)
	streamCmd := pipe.HGetAll(ctx, streamMetaKeyPrefix+streamID)
	if _, err := pipe.Exec(ctx); err != nil {
		return domain.Snapshot{}, fmt.Errorf("load entitlement snapshot: %w", err)
	}
	subscription, err := toSubscription(subscriptionCmd.Val())
	if err != nil {
		return domain.Snapshot{}, err
	}
	return domain.Snapshot{Subscription: subscription, Stream: toStream(streamID, streamCmd.Val())}, nil
}

func toSubscription(fields map[string]string) (*domain.Subscription, error) {
	if len(fields) == 0 {
		return nil, nil
	}
	expiresAt, err := time.Parse(time.RFC3339, fields[fieldExpiresAt])
	if err != nil {
		return nil, fmt.Errorf("parse subscription %s: %w", fieldExpiresAt, err)
	}
	return &domain.Subscription{Status: fields[fieldStatus], ExpiresAt: expiresAt}, nil
}

func toStream(streamID string, fields map[string]string) *domain.Stream {
	if len(fields) == 0 {
		return nil
	}
	return &domain.Stream{ID: streamID, State: domain.StreamState(fields[fieldState])}
}
