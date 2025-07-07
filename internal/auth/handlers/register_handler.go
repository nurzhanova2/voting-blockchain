package handlers

import (
	"encoding/json"
	"net/http"

	"voting-blockchain/internal/auth/dto"
	"voting-blockchain/internal/auth/services"
)

type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler получает AuthService, чтобы внутри вызывать бизнес-логику
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	// Декодируем JSON-запрос
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Не удалось отправить ответ", http.StatusInternalServerError)
	}
}
