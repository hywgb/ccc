package handler

import (
	"encoding/json"
	"net/http"

	"github.com/divord97/ccc/internal/domain/ai"
	"github.com/divord97/ccc/internal/interfaces/http/middleware"
	"github.com/divord97/ccc/pkg/response"
)

type PerformanceHandler struct {
	svc *ai.PerformanceScorecardService
}

func NewPerformanceHandler(svc *ai.PerformanceScorecardService) *PerformanceHandler {
	return &PerformanceHandler{svc: svc}
}

func (h *PerformanceHandler) Generate(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	var in ai.GenerateScorecardInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	in.TenantID = tenantID
	sc, err := h.svc.Generate(r.Context(), in)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, sc)
}

func (h *PerformanceHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	period := r.URL.Query().Get("period")
	list, err := h.svc.List(r.Context(), tenantID, period)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, list)
}
