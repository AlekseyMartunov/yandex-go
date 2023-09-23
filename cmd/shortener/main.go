package main

import (
	"database/sql"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/encoder"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/handlers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/authentication"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/compress"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/logger"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/migrations"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/urlpostgres"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/users/simpleusers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/users/userspostgres"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/router"
)

func main() {
	cfg := config.NewConfig()
	cfg.GetConfig()

	URLDB, err := createStorageURL("pgx", cfg)
	if err != nil {
		panic(err)
	}

	defer URLDB.Close()

	dbUser, err := createStorageUser("pgx", cfg)
	if err != nil {
		panic(err)
	}

	encoder := encoder.NewEncoder(URLDB)
	handler := handlers.NewShortURLHandler(encoder, cfg)

	tokenController := authentication.NewTokenController(dbUser)
	log := logger.NewLogger("info")

	router := router.NewBaseRouter(
		handler,
		log.WithLogging,
		compress.Compress,
		tokenController.CheckToken,
	)

	err = http.ListenAndServe(cfg.GetAddress(), router.Route())
	if err != nil {
		panic(err)
	}
}

func createStorageURL(driverName string, cfg *config.Config) (simpleurl.Storage, error) {
	if cfg.GetDataBaseDSN() != "" {
		conn, err := sql.Open(driverName, cfg.GetDataBaseDSN())
		if err != nil {
			return nil, err
		}
		err = migrations.MakeMigration(conn)
		if err != nil {
			return nil, err
		}
		postgresDB := urlpostgres.NewDB(conn)
		return postgresDB, nil
	}

	if cfg.GetFileStoragePath() != "" {
		fileStorage, err := simpleurl.NewFileStorage(cfg.GetFileStoragePath())
		if err != nil {
			return nil, err
		}
		return fileStorage, nil
	}

	mapStorage, err := simpleurl.NewMapStorage()
	if err != nil {
		return nil, err
	}
	return mapStorage, nil
}

func createStorageUser(driverName string, cfg *config.Config) (simpleusers.Users, error) {
	if cfg.GetDataBaseDSN() != "" {
		conn, err := sql.Open(driverName, cfg.GetDataBaseDSN())
		if err != nil {
			return simpleusers.NewUser(), nil
		}
		err = migrations.MakeMigration(conn)
		if err != nil {
			return nil, err
		}
		return userspostgres.NewUserModel(conn), nil
	}
	return simpleusers.NewUser(), nil

}
