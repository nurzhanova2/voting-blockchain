package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"voting-blockchain/internal/voting/dto"
	"voting-blockchain/internal/voting/models"
	"voting-blockchain/internal/voting/services"
	authhandlers "voting-blockchain/internal/auth/handlers"
)

type ElectionHandler struct {
	service services.ElectionService
}

func NewElectionHandler(s services.ElectionService) *ElectionHandler {
	return &ElectionHandler{service: s}
}

func (h *ElectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	role, err := authhandlers.GetUserRole(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if role != "admin" {
		http.Error(w, "forbidden: only admin can create elections", http.StatusForbidden)
		return
	}

	var req dto.CreateElectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	userID, err := authhandlers.GetUserID(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	e := &models.Election{
		Title:       req.Title,
		Description: req.Description,
		IsActive:    req.IsActive,
		CreatedBy:   userID,
	}

	if err := h.service.Create(r.Context(), e); err != nil {
		http.Error(w, "failed to create election: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(req.Choices) > 0 {
		if err := h.service.CreateChoices(r.Context(), e.ID, req.Choices); err != nil {
			http.Error(w, "failed to save choices: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (h *ElectionHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid election ID", http.StatusBadRequest)
		return
	}

	e, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func (h *ElectionHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, "failed to list elections: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *ElectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	role, err := authhandlers.GetUserRole(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if role != "admin" {
		http.Error(w, "forbidden: only admin can update elections", http.StatusForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid election ID", http.StatusBadRequest)
		return
	}

	var req dto.UpdateElectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	e := &models.Election{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		IsActive:    req.IsActive,
	}

	if err := h.service.Update(r.Context(), e); err != nil {
		http.Error(w, "failed to update election: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func (h *ElectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	role, err := authhandlers.GetUserRole(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if role != "admin" {
		http.Error(w, "forbidden: only admin can delete elections", http.StatusForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid election ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete election: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
