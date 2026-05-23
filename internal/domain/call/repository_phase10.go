package call

import "context"

// WebRTCQualityRepository manages WebRTC quality log persistence.
type WebRTCQualityRepository interface {
	Create(ctx context.Context, log *WebRTCQualityLog) error
	ListByCallID(ctx context.Context, callID int64) ([]WebRTCQualityLog, error)
	ListByAgent(ctx context.Context, tenantID, agentID int64, limit int) ([]WebRTCQualityLog, error)
}
