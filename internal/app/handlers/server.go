package handlers

import (
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/config"
	"github.com/AlekseyMartunov/yandex-go.git/internal/app/storage"
)

type Server struct {
	db  *storage.Storage
	Cfg *config.NetAddress
}

func NewServer(db *storage.Storage, cfg *config.NetAddress) *Server {
	return &Server{db: db, Cfg: cfg}
}
