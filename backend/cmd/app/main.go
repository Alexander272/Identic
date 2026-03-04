package main

import (
	"context"
	"log"
	"os"

	"github.com/Alexander272/Identic/backend/internal/config"
	"github.com/Alexander272/Identic/backend/internal/migrate"
	"github.com/Alexander272/Identic/backend/pkg/database/postgres"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/subosito/gotenv"
)

func main() {
	//* Init config
	if os.Getenv("APP_ENV") == "" {
		if err := gotenv.Load(".env"); err != nil {
			log.Fatalf("error loading env variables: %s", err.Error())
		}
	}

	conf, err := config.Init("configs/config.yaml")
	if err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	logger.NewLogger(logger.WithLevel(conf.LogLevel), logger.WithAddSource(conf.LogSource))

	//* Dependencies
	db, err := postgres.NewPostgresDB(context.Background(), &postgres.Config{
		Host:     conf.Postgres.Host,
		Port:     conf.Postgres.Port,
		Username: conf.Postgres.Username,
		Password: conf.Postgres.Password,
		DBName:   conf.Postgres.DbName,
		SSLMode:  conf.Postgres.SSLMode,
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	if err := migrate.Migrate(db); err != nil {
		log.Fatalf("failed to migrate: %s", err.Error())
	}
}
