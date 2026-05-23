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

type DigitalEmployeeHandler struct {
	svc *ai.DigitalEmployeeService
}

func NewDigitalEmployeeHandler(svc *ai.DigitalEmployeeService) *DigitalEmployeeHandler {
	return &DigitalEmployeeHandler{svc: svc}
}

func (h *DigitalEmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	var in ai.CreateDigitalEmployeeInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	in.TenantID = tenantID
	de, err := h.svc.Create(r.Context(), in)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, de)
}

func (h *DigitalEmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	de, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, de)
}

func (h *DigitalEmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	de, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	var in struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		AvatarURL   *string `json:"avatar_url"`
		IsActive    *bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if in.Name != nil {
		de.Name = *in.Name
	}
	if in.Description != nil {
		de.Description = *in.Description
	}
	if in.AvatarURL != nil {
		de.AvatarURL = *in.AvatarURL
	}
	if in.IsActive != nil {
		de.IsActive = *in.IsActive
	}
	if err := h.svc.Update(r.Context(), de); err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, de)
}

func (h *DigitalEmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	list, err := h.svc.List(r.Context(), tenantID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, list)
}

func (h *DigitalEmployeeHandler) CreateScene(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantIDFromCtx(r.Context())
	deID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	var in ai.CreateSceneInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	in.DigitalEmployeeID = deID
	in.TenantID = tenantID
	scene, err := h.svc.CreateScene(r.Context(), in)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, scene)
}

func (h *DigitalEmployeeHandler) ListScenes(w http.ResponseWriter, r *http.Request) {
	deID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	list, err := h.svc.ListScenes(r.Context(), deID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, list)
}

func (h *DigitalEmployeeHandler) PublishScene(w http.ResponseWriter, r *http.Request) {
	sceneID, _ := strconv.ParseInt(chi.URLParam(r, "sceneId"), 10, 64)
	scene, err := h.svc.PublishScene(r.Context(), sceneID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, scene)
}

func (h *DigitalEmployeeHandler) TestIntent(w http.ResponseWriter, r *http.Request) {
	sceneID, _ := strconv.ParseInt(chi.URLParam(r, "sceneId"), 10, 64)
	var in struct {
		Input string `json:"input"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.svc.MatchIntent(r.Context(), sceneID, in.Input)
	if err != nil {
		response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, result)
}
