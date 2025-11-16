package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	// Blank import needed to register the postgres driver.
	_ "github.com/lib/pq"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/config"
)

func Open(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName,
			cfg.Database.SSLMode),
	)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Close(db *sqlx.DB) {
	if err := db.Close(); err != nil {
		panic(err)
	}
}
