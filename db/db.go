package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func InitDB() {
	var err error
	connStr := "postgres://owner:123@localhost:5432/somemarket"
	DB, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
}
