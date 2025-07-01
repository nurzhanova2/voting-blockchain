package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"voting-blockchain/internal/voting/dto"
	"voting-blockchain/internal/voting/services"
)

// BlockchainHandler — обработчик блокчейна
type BlockchainHandler struct {
	blockchain services.BlockchainService
}

// NewBlockchainHandler — конструктор
func NewBlockchainHandler(svc services.BlockchainService) *BlockchainHandler {
	return &BlockchainHandler{blockchain: svc}
}

// POST /elections/{id}/vote
func (h *BlockchainHandler) Vote(w http.ResponseWriter, r *http.Request) {
	electionIDStr := chi.URLParam(r, "id")
	electionID, err := strconv.Atoi(electionIDStr)
	if err != nil {
		http.Error(w, "invalid election id", http.StatusBadRequest)
		return
	}

	var req dto.AddBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.VoteHash == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	block, err := h.blockchain.AddBlock(r.Context(), electionID, req.VoteHash)
	if err != nil {
		http.Error(w, "failed to add block: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(block)
}

// GET /elections/{id}/blocks
func (h *BlockchainHandler) GetChain(w http.ResponseWriter, r *http.Request) {
	electionIDStr := chi.URLParam(r, "id")
	electionID, err := strconv.Atoi(electionIDStr)
	if err != nil {
		http.Error(w, "invalid election id", http.StatusBadRequest)
		return
	}

	chain, err := h.blockchain.GetChain(r.Context(), electionID)
	if err != nil {
		http.Error(w, "failed to load chain: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chain)
}
