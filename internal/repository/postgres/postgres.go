package postgres

import (
	"fmt"
	"github.com/Verce11o/yata-comments/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func NewPostgres(cfg *config.Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Name))

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	return db
}
