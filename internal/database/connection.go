package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func ConnectPG() *pgxpool.Pool {
	DSN := viper.GetString("DATABASE_URL")
	if DSN == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// context with 60 seconds timeout
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, DSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database")

	return conn
}
