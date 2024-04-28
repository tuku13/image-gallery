package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var Db *sqlx.DB

func init() {
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "image-gallery-db"
	}
	usr := os.Getenv("DB_USER")
	if usr == "" {
		usr = "test"
	}
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = "test"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	uri := "user=" + usr + " dbname=" + dbname + " password=" + pass + " host=" + host + " port=" + port + " sslmode=disable"
	db, err := sqlx.Connect("postgres", uri)

	if err != nil {
		log.Fatalln(err)
	}
	Db = db
}
