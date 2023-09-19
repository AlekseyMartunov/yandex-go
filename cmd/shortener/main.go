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
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simplestotage"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/urlpostgres"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/users/simpleusers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/users/userspostgres"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/router"
)

func main() {
	cfg := config.NewConfig()
	cfg.GetConfig()

	conn, err := getConnectionPool(cfg)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	db, err := createURLStorage(conn, cfg)
	if err != nil {
		panic(nil)
	}

	dbUser, err := createUserStorage(conn)
	if err != nil {
		panic(nil)
	}

	encoder := encoder.NewEncoder(db)
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

func getConnectionPool(cfg *config.Config) (*sql.DB, error) {
	if cfg.GetDataBaseDSN() != "" {
		conn, err := sql.Open("pgx", cfg.GetDataBaseDSN())
		if err != nil {
			return nil, err
		}
		err = migrations.MakeMigration(conn)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	return nil, nil
}

func createURLStorage(conn *sql.DB, cfg *config.Config) (simplestotage.Storage, error) {
	if conn != nil {
		db := urlpostgres.NewDB(conn)
		return db, nil
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

func createUserStorage(conn *sql.DB) (simpleusers.Users, error) {
	if conn != nil {
		db := userspostgres.NewUserModel(conn)
		return db, nil
	}
	db := simpleusers.NewUser()
	return db, nil
}
