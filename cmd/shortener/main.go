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
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/postgres"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/simpleStotage"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/router"
)

func main() {
	cfg := config.NewConfig()
	cfg.GetConfig()

	db := createDB("pgx", cfg)
	defer db.Close()

	encoder := encoder.NewEncoder(db)

	handler := handlers.NewShortURLHandler(encoder, cfg)

	log := logger.NewLogger("info")
	router := router.NewBaseRouter(handler, log.WithLogging, compress.Compress)

	err := http.ListenAndServe(cfg.GetAddress(), router.Route())
	if err != nil {
		panic(err)
	}
}

func createDB(driverName string, cfg *config.Config) simpleStotage.Storage {
	db, err := sql.Open(driverName, cfg.GetDataBaseDSN())
	if err != nil {
		return simpleStotage.NewStorage(cfg.GetFileStoragePath())
	}

	err = db.Ping()
	if err != nil {
		return simpleStotage.NewStorage(cfg.GetFileStoragePath())
	}

	postgresDB := postgres.NewDB(db)
	err = postgresDB.CreateTableURL()
	if err != nil {
		return simpleStotage.NewStorage(cfg.GetFileStoragePath())
	}
	cfg.SetDataBaseStatus(true)

	return postgresDB
}
