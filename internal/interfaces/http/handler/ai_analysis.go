package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/divord97/ccc/internal/application/aianalysis"
	"github.com/divord97/ccc/internal/interfaces/http/middleware"
	"github.com/divord97/ccc/pkg/response"
	"github.com/go-chi/chi/v5"
)

type AIAnalysisHandler struct {
	svc *aianalysis.Service
}

func NewAIAnalysisHandler(svc *aianalysis.Service) *AIAnalysisHandler {
	return &AIAnalysisHandler{svc: svc}
}

func (h *AIAnalysisHandler) Summary(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	var in struct {
		Transcript string `json:"transcript"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.GenerateSummary(r.Context(), callID, in.Transcript)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) Sentiment(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	var in struct {
		Transcript string `json:"transcript"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.AnalyzeSentiment(r.Context(), callID, in.Transcript)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) Tags(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	var in struct {
		Transcript string `json:"transcript"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.ExtractTags(r.Context(), callID, in.Transcript)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) Satisfaction(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	var in struct {
		Transcript string `json:"transcript"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.PredictSatisfaction(r.Context(), callID, in.Transcript)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) IVRAnalysis(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	var in struct {
		IVRPath string `json:"ivr_path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.AnalyzeIVRPath(r.Context(), callID, in.IVRPath)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) Completion(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	var in struct {
		Transcript string `json:"transcript"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.JudgeCompletion(r.Context(), callID, in.Transcript)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) PostCallActions(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	var in struct {
		Transcript string `json:"transcript"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.ExtractPostCallActions(r.Context(), callID, in.Transcript)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) AutoFill(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	var in struct {
		Transcript string `json:"transcript"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.AutoFillTicket(r.Context(), callID, in.Transcript)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) ScriptRecommend(w http.ResponseWriter, r *http.Request) {
	callID, _ := strconv.ParseInt(chi.URLParam(r, "callId"), 10, 64)
	var in struct {
		Transcript string   `json:"transcript"`
		Scripts    []string `json:"scripts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.RecommendScript(r.Context(), callID, in.Transcript, in.Scripts)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) BatchTags(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	var in aianalysis.BatchTagAnalysisInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	in.TenantID = tenantID
	result, err := h.svc.RunBatchTagAnalysis(r.Context(), in)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AIAnalysisHandler) HotwordAnalysis(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Transcripts []string `json:"transcripts"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.AnalyzeHotwords(r.Context(), in.Transcripts)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}
