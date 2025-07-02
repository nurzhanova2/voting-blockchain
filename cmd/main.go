package main

import (
    "log"
    "net/http"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v4/middleware"
    "github.com/go-chi/cors"

    "voting-blockchain/config"
    "voting-blockchain/db"

    // Auth-модуль
    authHandlers "voting-blockchain/internal/auth/handlers"
    authRepos "voting-blockchain/internal/auth/repositories"
    authRouters "voting-blockchain/internal/auth/routers"
    authServices "voting-blockchain/internal/auth/services"

    // Voting-модуль
    votingHandlers "voting-blockchain/internal/voting/handlers"
    votingRepos "voting-blockchain/internal/voting/repositories"
    votingRouters "voting-blockchain/internal/voting/routers"
    votingServices "voting-blockchain/internal/voting/services"
)

func main() {
    // Загружаем конфигурацию из .env
    cfg := config.LoadConfig()

    // Инициализируем БД
    db.InitDB(cfg)

    // Конвертация TTL для токенов
    accessTTL := time.Duration(cfg.AccessTokenTTLMin) * time.Minute
    refreshTTL := time.Duration(cfg.RefreshTokenTTLDays*24) * time.Hour

    // ===== AUTH =====
    userRepo := authRepos.NewUserRepository()
    refreshRepo := authRepos.NewRefreshTokenRepository()
    authService := authServices.NewAuthService(userRepo, refreshRepo, cfg.JWTSecret, accessTTL, refreshTTL)
    authHandler := authHandlers.NewAuthHandler(authService)

    // ===== VOTING =====
    voteRepo := votingRepos.NewVotePostgres(db.DB)
    blockchainRepo := votingRepos.NewBlockchainPostgres(db.DB)
    electionRepo := votingRepos.NewElectionPostgres(db.DB)
    choiceRepo := votingRepos.NewChoicePostgres(db.DB) // добавлено

    electionService := votingServices.NewElectionService(electionRepo, choiceRepo) // изменено
    electionHandler := votingHandlers.NewElectionHandler(electionService)

    voteService := votingServices.NewVoteService(voteRepo, blockchainRepo)
    voteHandler := votingHandlers.NewVoteHandler(voteService)


    // ===== ROUTING =====
    r := chi.NewRouter()

    // ==== MIDDLEWARE ====
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // CORS разрешаем фронту доступ
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    // (опционально) префикс API
    r.Route("/api/v1", func(api chi.Router) {
        // Auth маршруты
        api.Mount("/auth", authRouters.NewAuthRouter(authHandler, []byte(cfg.JWTSecret)))

        // Voting маршруты (с JWT)
        api.Mount("/voting",
            authHandlers.NewJWTMiddleware([]byte(cfg.JWTSecret))(
                votingRouters.NewVotingRouter(voteHandler, electionHandler),
            ),
        )
    })

    // Запуск сервера
    log.Println("Сервер запущен на :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}