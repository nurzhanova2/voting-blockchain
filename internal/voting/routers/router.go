package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"voting-blockchain/internal/voting/handlers"
)

// NewVotingRouter создает роутер для голосования и управления выборами.
// Принимает VoteHandler и ElectionHandler.
func NewVotingRouter(voteHandler *handlers.VoteHandler, electionHandler *handlers.ElectionHandler) http.Handler {
	r := chi.NewRouter()

	// Эндпоинты голосования
	r.Post("/elections/{id}/vote", voteHandler.CastVoteHandler)
	r.Get("/elections/{id}/blocks", voteHandler.GetElectionBlockchainHandler)

	// CRUD выборов
	r.Route("/elections", func(r chi.Router) {
		r.Post("/", electionHandler.Create)
		r.Get("/", electionHandler.List)
		r.Get("/{id}", electionHandler.Get)
		r.Put("/{id}", electionHandler.Update)
		r.Delete("/{id}", electionHandler.Delete)
	})

	return r
}
