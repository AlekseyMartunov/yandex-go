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
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/simplestotage"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/router"
)

func main() {
	cfg := config.NewConfig()
	cfg.GetConfig()

	db, err := createDB("pgx", cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	encoder := encoder.NewEncoder(db)

	handler := handlers.NewShortURLHandler(encoder, cfg)

	log := logger.NewLogger("info")
	router := router.NewBaseRouter(handler, log.WithLogging, compress.Compress)

	err = http.ListenAndServe(cfg.GetAddress(), router.Route())
	if err != nil {
		panic(err)
	}
}

func createDB(driverName string, cfg *config.Config) (simplestotage.Storage, error) {
	if cfg.GetDataBaseDSN() != "" {
		db, err := sql.Open(driverName, cfg.GetDataBaseDSN())
		if err != nil {
			return nil, err
		}
		postgresDB := postgres.NewDB(db)
		err = postgresDB.CreateTableURL()
		if err != nil {
			return nil, err
		}
		return postgresDB, nil
	}

	if cfg.GetFileStoragePath() != "" {
		fileStorage, err := simplestotage.NewFileStorage(cfg.GetFileStoragePath())
		if err != nil {
			return nil, err
		}
		return fileStorage, nil
	}

	mapStorage, err := simplestotage.NewMapStorage()
	if err != nil {
		return nil, err
	}
	return mapStorage, nil

}
