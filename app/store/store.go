package store

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Store struct {
	config *Config
	Db     *sql.DB
}

func Open(config *Config) *Store {
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("connected to the database")
	return &Store{
		config: config,
		Db:     db,
	}
}
