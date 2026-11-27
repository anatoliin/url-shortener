package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func NewDBConnection(string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
}
