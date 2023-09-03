package main

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/encoder"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/handlers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/compress"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/logger"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/router"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/storage"
)

func main() {
	cfg := config.NewConfig()
	cfg.GetConfig()

	_ = runDB("pgx", cfg)

	storage := storage.NewStorage(cfg.GetFileStoragePath())

	encoder := encoder.NewEncoder(storage)

	handler := handlers.NewShortURLHandler(encoder, cfg)

	log := logger.NewLogger("info")
	router := router.NewBaseRouter(handler, log.WithLogging, compress.Compress)

	err := http.ListenAndServe(cfg.GetAddress(), router.Route())
	if err != nil {
		panic(err)
	}
}

func runDB(driverName string, cfg *config.Config) *sql.DB {
	db, err := sql.Open(driverName, cfg.GetDataBaseDSN())
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	cfg.SetDataBaseStatus(true)

	defer db.Close()
	return db
}
