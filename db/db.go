package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func MustInitPostgres(url string) {
	var err error
	pool, err = pgxpool.New(context.Background(), url)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}
	log.Printf("Intialized db %q", url)
}

func Close() {
	pool.Close()
	log.Println("Closed db")
}
