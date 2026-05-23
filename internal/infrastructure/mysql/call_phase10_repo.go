package mysql

import (
	"context"

	"github.com/divord97/ccc/internal/domain/call"
	"github.com/jmoiron/sqlx"
)

// --- WebRTCQualityRepo ---

type WebRTCQualityRepo struct{ db *sqlx.DB }

func NewWebRTCQualityRepo(db *sqlx.DB) *WebRTCQualityRepo {
	return &WebRTCQualityRepo{db: db}
}

func (r *WebRTCQualityRepo) Create(ctx context.Context, l *call.WebRTCQualityLog) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO webrtc_quality_logs (id, call_id, tenant_id, agent_id, packet_loss_rate, jitter, round_trip_time, mos, audio_level, bitrate_kbps, codec_name, sampled_at, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		l.ID, l.CallID, l.TenantID, l.AgentID, l.PacketLossRate, l.Jitter, l.RoundTripTime, l.MOS, l.AudioLevel, l.BitrateKbps, l.CodecName, l.SampledAt, l.CreatedAt)
	return err
}

func (r *WebRTCQualityRepo) ListByCallID(ctx context.Context, callID int64) ([]call.WebRTCQualityLog, error) {
	var out []call.WebRTCQualityLog
	if err := r.db.SelectContext(ctx, &out, "SELECT * FROM webrtc_quality_logs WHERE call_id = ? ORDER BY sampled_at", callID); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *WebRTCQualityRepo) ListByAgent(ctx context.Context, tenantID, agentID int64, limit int) ([]call.WebRTCQualityLog, error) {
	var out []call.WebRTCQualityLog
	if err := r.db.SelectContext(ctx, &out, "SELECT * FROM webrtc_quality_logs WHERE tenant_id = ? AND agent_id = ? ORDER BY sampled_at DESC LIMIT ?", tenantID, agentID, limit); err != nil {
		return nil, err
	}
	return out, nil
}
