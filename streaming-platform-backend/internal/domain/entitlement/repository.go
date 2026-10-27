package entitlement

import "context"

type Repository interface {
	Load(ctx context.Context, userID string, streamID string) (Snapshot, error)
}
