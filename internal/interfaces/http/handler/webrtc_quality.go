package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/divord97/ccc/internal/domain/call"
	"github.com/divord97/ccc/internal/interfaces/http/middleware"
	"github.com/divord97/ccc/pkg/response"
	"github.com/divord97/ccc/pkg/snowflake"
	"github.com/go-chi/chi/v5"
)

type WebRTCQualityHandler struct {
	repo call.WebRTCQualityRepository
}

func NewWebRTCQualityHandler(repo call.WebRTCQualityRepository) *WebRTCQualityHandler {
	return &WebRTCQualityHandler{repo: repo}
}

func (h *WebRTCQualityHandler) Save(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	var in struct {
		CallID         int64   `json:"call_id"`
		AgentID        int64   `json:"agent_id"`
		PacketLossRate float64 `json:"packet_loss_rate"`
		Jitter         float64 `json:"jitter"`
		RoundTripTime  float64 `json:"round_trip_time"`
		MOS            float64 `json:"mos"`
		AudioLevel     float64 `json:"audio_level"`
		BitrateKbps    int     `json:"bitrate_kbps"`
		CodecName      string  `json:"codec_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	log := &call.WebRTCQualityLog{
		ID:             snowflake.NextID(),
		CallID:         in.CallID,
		TenantID:       tenantID,
		AgentID:        in.AgentID,
		PacketLossRate: in.PacketLossRate,
		Jitter:         in.Jitter,
		RoundTripTime:  in.RoundTripTime,
		MOS:            in.MOS,
		AudioLevel:     in.AudioLevel,
		BitrateKbps:    in.BitrateKbps,
		CodecName:      in.CodecName,
		SampledAt:      time.Now(),
		CreatedAt:      time.Now(),
	}
	if err := h.repo.Create(r.Context(), log); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, log)
}

func (h *WebRTCQualityHandler) ListByCall(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	logs, err := h.repo.ListByCallID(r.Context(), callID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, logs)
}

func (h *WebRTCQualityHandler) ListByAgent(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	agentID, _ := strconv.ParseInt(chi.URLParam(r, "agentId"), 10, 64)
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	logs, err := h.repo.ListByAgent(r.Context(), tenantID, agentID, limit)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, logs)
}
