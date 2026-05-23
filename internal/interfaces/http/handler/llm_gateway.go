package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/divord97/ccc/internal/domain/ai"
	"github.com/divord97/ccc/internal/interfaces/http/middleware"
	"github.com/divord97/ccc/pkg/response"
	"github.com/go-chi/chi/v5"
)

type LLMGatewayHandler struct {
	svc *ai.LLMGatewayService
}

func NewLLMGatewayHandler(svc *ai.LLMGatewayService) *LLMGatewayHandler {
	return &LLMGatewayHandler{svc: svc}
}

func (h *LLMGatewayHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	var in ai.CreateLLMModelConfigInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	in.TenantID = tenantID
	config, err := h.svc.CreateConfig(r.Context(), in)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, config)
}

func (h *LLMGatewayHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	config, err := h.svc.GetConfig(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, config)
}

func (h *LLMGatewayHandler) ListConfigs(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	configs, err := h.svc.ListConfigs(r.Context(), tenantID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, configs)
}

func (h *LLMGatewayHandler) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	config, err := h.svc.GetConfig(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	var in struct {
		Name         string `json:"name"`
		ProviderType string `json:"provider_type"`
		Endpoint     string `json:"endpoint"`
		APIKey       string `json:"api_key"`
		ModelName    string `json:"model_name"`
		IsDefault    *bool  `json:"is_default"`
		IsActive     *bool  `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if in.Name != "" {
		config.Name = in.Name
	}
	if in.ProviderType != "" {
		config.ProviderType = in.ProviderType
	}
	if in.Endpoint != "" {
		config.Endpoint = in.Endpoint
	}
	if in.APIKey != "" {
		config.APIKey = in.APIKey
	}
	if in.ModelName != "" {
		config.ModelName = in.ModelName
	}
	if in.IsDefault != nil {
		config.IsDefault = *in.IsDefault
	}
	if in.IsActive != nil {
		config.IsActive = *in.IsActive
	}
	if err := h.svc.UpdateConfig(r.Context(), config); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, config)
}

func (h *LLMGatewayHandler) DeleteConfig(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err := h.svc.DeleteConfig(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}

func (h *LLMGatewayHandler) GetDefault(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	config, err := h.svc.GetDefault(r.Context(), tenantID)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, config)
}
