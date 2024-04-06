package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"learning.com/golang_backend/api"
	db "learning.com/golang_backend/db/sqlc/repository"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	repository := db.NewRepository(conn)
	server := api.NewServer(repository)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
