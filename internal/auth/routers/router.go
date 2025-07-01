package routers

import (
	"net/http"

	"voting-blockchain/internal/auth/handlers"

	"github.com/go-chi/chi/v5"
)

// NewAuthRouter создает маршруты аутентификации.
// Принимает AuthHandler и JWT секрет.
func NewAuthRouter(handler *handlers.AuthHandler, jwtSecret []byte) http.Handler {
	r := chi.NewRouter()

	// Открытые маршруты
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
	r.Post("/refresh", handler.Refresh)

	// Защищенные маршруты
	r.Group(func(r chi.Router) {
		r.Use(handlers.NewJWTMiddleware(jwtSecret))
		r.Get("/profile", handler.Profile)
	})

	return r
}

