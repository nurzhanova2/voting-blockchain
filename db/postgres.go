package db

import (
	"context"
	"log"
	"time"

	"voting-blockchain/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	DB, err = pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	err = DB.Ping(ctx)
	if err != nil {
		log.Fatalf("БД не отвечает: %v", err)
	}

	log.Println("Подключение к PostgreSQL установлено")
}
