//package main
//
//import (
//	"database/sql"
//	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
//	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/migrations"
//	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simplestotage"
//	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/urlpostgres"
//	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/users/simpleusers"
//	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/users/userspostgres"
//)
//
//func getConnectionPool(cfg *config.Config) (*sql.DB, error) {
//	if cfg.GetDataBaseDSN() != "" {
//		conn, err := sql.Open("pgx", cfg.GetDataBaseDSN())
//		if err != nil {
//			return nil, err
//		}
//		err = migrations.MakeMigration(conn)
//		if err != nil {
//			return nil, err
//		}
//		return conn, nil
//	}
//	return nil, nil
//}
//
//func createURLStorage(conn *sql.DB, cfg *config.Config) (simplestotage.Storage, error) {
//	if conn != nil {
//		db := urlpostgres.NewDB(conn)
//		return db, nil
//	}
//
//	if cfg.GetFileStoragePath() != "" {
//		fileStorage, err := simplestotage.NewFileStorage(cfg.GetFileStoragePath())
//		if err != nil {
//			return nil, err
//		}
//		return fileStorage, nil
//	}
//
//	mapStorage, err := simplestotage.NewMapStorage()
//	if err != nil {
//		return nil, err
//	}
//	return mapStorage, nil
//
//}
//
//func createUserStorage(conn *sql.DB) (simpleusers.Users, error) {
//	if conn != nil {
//		db := userspostgres.NewUserModel(conn)
//		return db, nil
//	}
//	db := simpleusers.NewUser()
//	return db, nil
//}
