package utils

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	// _ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	pgx "github.com/jackc/pgx/v4"

	log "github.com/sirupsen/logrus"
)

// create connection with postgres db
func NewPostgresConnection(config PostgresConfig) *pgx.Conn {

	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.User, config.Pass, config.Host, config.Port, config.DB)
	conn, err := pgx.Connect(context.Background(), postgresUrl)

	if err != nil {
		log.Error("Cannot connect to postgres")
	}

	log.Info("Successfully connected!")
	// return the connection
	return conn
}

func MigratePostgresUp(config PostgresConfig) {
	migrateUrl := fmt.Sprintf("pgx://%s:%s@%s:%s/%s", config.User, config.Pass, config.Host, config.Port, config.DB)
	m, err := migrate.New("file://migrations/", migrateUrl)
	if err != nil {
		log.Fatalf("New migrate: %s", err)
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Infof("Migrate up: %s", err)
		} else {
			log.Fatalf("Migrate up: %s", err)
		}
	}
}
