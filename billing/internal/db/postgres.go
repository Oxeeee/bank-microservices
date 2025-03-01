package db

import (
	"fmt"
	"log"

	"github.com/Oxeeee/bank-microservices/billing/internal/config"
	"github.com/jmoiron/sqlx"
)

func InitDB(cfg *config.Config) *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=%v",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Address,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode)
	DB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("error while connecting DB: %v", err)
	}

	return DB
}
