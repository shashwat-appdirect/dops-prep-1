package handlers

import (
	"appdirect-workshop-backend/internal/config"
	"appdirect-workshop-backend/internal/database"
)

type Handlers struct {
	db  database.DatabaseInterface
	cfg *config.Config
}

func New(db database.DatabaseInterface, cfg *config.Config) *Handlers {
	return &Handlers{
		db:  db,
		cfg: cfg,
	}
}

