package call

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebRTCQuality_QualityLevel(t *testing.T) {
	tests := []struct {
		mos    float64
		expect string
	}{
		{4.5, "good"},
		{4.0, "good"},
		{3.5, "fair"},
		{3.0, "fair"},
		{2.5, "poor"},
		{1.0, "poor"},
	}
	for _, tt := range tests {
		log := &WebRTCQualityLog{MOS: tt.mos}
		assert.Equal(t, tt.expect, log.QualityLevel(), "MOS %.1f", tt.mos)
	}
}

func TestWebRTCQuality_CreateAndList(t *testing.T) {
	repo := NewMockWebRTCQualityRepo()
	ctx := context.Background()

	log1 := &WebRTCQualityLog{
		CallID: 100, TenantID: 1, AgentID: 10,
		PacketLossRate: 0.02, Jitter: 15.5, RoundTripTime: 80,
		MOS: 4.2, AudioLevel: -30, BitrateKbps: 64,
		CodecName: "opus", SampledAt: time.Now(),
	}
	require.NoError(t, repo.Create(ctx, log1))
	assert.NotZero(t, log1.ID)

	log2 := &WebRTCQualityLog{
		CallID: 100, TenantID: 1, AgentID: 10,
		PacketLossRate: 0.05, Jitter: 25, RoundTripTime: 120,
		MOS: 3.5, SampledAt: time.Now(),
	}
	require.NoError(t, repo.Create(ctx, log2))

	logs, err := repo.ListByCallID(ctx, 100)
	require.NoError(t, err)
	assert.Len(t, logs, 2)

	logs, err = repo.ListByAgent(ctx, 1, 10, 10)
	require.NoError(t, err)
	assert.Len(t, logs, 2)
}
