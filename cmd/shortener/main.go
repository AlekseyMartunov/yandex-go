package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/compress"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/encoder"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/handlers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/authentication"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/middleware/logger"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/migrations"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/urlpostgres"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/users/simpleusers"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/users/userspostgres"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/router"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	ctx, stopApp := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopApp()

	go func() {
		startServer(ctx)
	}()

	<-ctx.Done()
}

func startServer(ctx context.Context) {
	cfg := config.NewConfig()
	cfg.GetConfig()

	conn, err := getConn("pgx", cfg)
	if err != nil {
		panic(err)
	}
	if conn != nil {
		defer conn.Close()
	}

	URLDB, err := createStorageURL(conn, cfg)
	if err != nil {
		panic(err)
	}

	dbUser, err := createStorageUser(conn, cfg)
	if err != nil {
		panic(err)
	}

	encoder := encoder.NewEncoder(URLDB)
	handler := handlers.NewShortURLHandler(encoder, cfg, ctx)

	tokenController := authentication.NewTokenController(dbUser)
	log := logger.NewLogger("info")

	router := router.NewBaseRouter(
		handler,
		log.WithLogging,
		compress.Compress,
		tokenController.CheckToken,
	)

	r := router.Route()
	r.Mount("/debug", middleware.Profiler())

	greet()

	err = http.ListenAndServe(cfg.GetAddress(), r)
	if err != nil {
		panic(err)
	}
}

func createStorageURL(conn *sql.DB, cfg *config.Config) (simpleurl.Storage, error) {
	if conn != nil && cfg.GetDataBaseDSN() != "" {
		return urlpostgres.NewDB(conn), nil
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

func createStorageUser(conn *sql.DB, cfg *config.Config) (simpleusers.Users, error) {
	if conn != nil && cfg.GetDataBaseDSN() != "" {
		return userspostgres.NewUserModel(conn), nil
	}
	return simpleusers.NewUser(), nil
}

func getConn(driverName string, cfg *config.Config) (*sql.DB, error) {
	if cfg.GetDataBaseDSN() != "" {
		conn, err := sql.Open(driverName, cfg.GetDataBaseDSN())
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

func greet() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}

	if buildDate == "" {
		buildDate = "N/A"
	}

	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
