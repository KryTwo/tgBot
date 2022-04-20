package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func NewDB() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "dbname=tg_bot_db user=root password=123456 sslmode=disable")
	if err != nil {
		log.Fatalf("cant connect to db %v", err)
	}
	return db
}
