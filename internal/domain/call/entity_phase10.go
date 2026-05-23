package call

import "time"

// WebRTCQualityLog records real-time WebRTC call quality metrics.
type WebRTCQualityLog struct {
	ID              int64     `db:"id" json:"id"`
	CallID          int64     `db:"call_id" json:"call_id"`
	TenantID        int64     `db:"tenant_id" json:"tenant_id"`
	AgentID         int64     `db:"agent_id" json:"agent_id"`
	PacketLossRate  float64   `db:"packet_loss_rate" json:"packet_loss_rate"`   // 0.0-1.0
	Jitter          float64   `db:"jitter" json:"jitter"`                       // ms
	RoundTripTime   float64   `db:"round_trip_time" json:"round_trip_time"`     // ms
	MOS             float64   `db:"mos" json:"mos"`                             // 1.0-5.0
	AudioLevel      float64   `db:"audio_level" json:"audio_level"`             // dBFS
	BitrateKbps     int       `db:"bitrate_kbps" json:"bitrate_kbps"`
	CodecName       string    `db:"codec_name" json:"codec_name"`
	SampledAt       time.Time `db:"sampled_at" json:"sampled_at"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// QualityLevel returns a traffic-light level based on MOS score.
func (w *WebRTCQualityLog) QualityLevel() string {
	switch {
	case w.MOS >= 4.0:
		return "good"
	case w.MOS >= 3.0:
		return "fair"
	default:
		return "poor"
	}
}
