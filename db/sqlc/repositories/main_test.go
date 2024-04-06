package repositories

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testRepository Repository

const source = "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable"

func TestMain(m *testing.M) {
	conn, err := pgxpool.New(context.Background(), source)
	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}
	testRepository = NewRepository(conn)

	os.Exit(m.Run())
}
