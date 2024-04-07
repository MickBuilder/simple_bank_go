package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"learning.com/golang_backend/api"
	db "learning.com/golang_backend/db/sqlc/repository"
	"learning.com/golang_backend/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config [err] :", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	repository := db.NewRepository(conn)
	server := api.NewServer(repository)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
