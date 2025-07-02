package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "voting-blockchain/internal/voting/dto"
    "voting-blockchain/internal/voting/services"

    authhandlers "voting-blockchain/internal/auth/handlers"
    "github.com/go-chi/chi/v5"
)

type VoteHandler struct {
    voteService services.VoteService
}

func NewVoteHandler(vs services.VoteService) *VoteHandler {
    return &VoteHandler{voteService: vs}
}

// CastVoteHandler — принимает голос от пользователя
func (h *VoteHandler) CastVoteHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := authhandlers.GetUserID(r)
    if err != nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    electionIDStr := chi.URLParam(r, "id")
    electionID, err := strconv.Atoi(electionIDStr)
    if err != nil {
        http.Error(w, "invalid election ID", http.StatusBadRequest)
        return
    }

    var req dto.CastVoteRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }

    if err := h.voteService.CastVote(r.Context(), userID, electionID, req.Choice); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// GetElectionBlockchainHandler — возвращает всю цепочку блоков для голосования
func (h *VoteHandler) GetElectionBlockchainHandler(w http.ResponseWriter, r *http.Request) {
    electionIDStr := chi.URLParam(r, "id")
    electionID, err := strconv.Atoi(electionIDStr)
    if err != nil {
        http.Error(w, "invalid election ID", http.StatusBadRequest)
        return
    }

    blocks, err := h.voteService.GetBlockchain(r.Context(), electionID)
    if err != nil {
        http.Error(w, "failed to get blockchain: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(blocks); err != nil {
        http.Error(w, "failed to encode response", http.StatusInternalServerError)
    }
}

func (h *VoteHandler) GetResults(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid election ID", http.StatusBadRequest)
        return
    }

    results, err := h.voteService.GetResults(r.Context(), id)
    if err != nil {
        http.Error(w, "failed to get results: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}
