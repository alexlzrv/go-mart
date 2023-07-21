package sql

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationsDir embed.FS

type Postgres struct {
	*sql.DB
}

func NewPostgresStorage(dsn string, log *zap.SugaredLogger) (*Postgres, error) {
	err := runMigrations(dsn, log)
	if err != nil {
		log.Errorf("attempt to establish connection failed: %s", err)
		return nil, err
	}

	db, err := newDBSession(dsn, log)
	if err != nil {
		log.Errorf("during initializing of new db session, error occurred: %s", err)
		return nil, err
	}

	return &Postgres{db}, nil
}

func runMigrations(dsn string, log *zap.SugaredLogger) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		log.Errorf("failed to return an iofs driver: %s", err)
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		log.Errorf("failed to get a new migrate instance: %s", err)
		return err
	}
	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Errorf("failed to apply migrations to the DB: %s", err)
			return err
		}
	}
	return nil
}

func newDBSession(dsn string, log *zap.SugaredLogger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Errorf("attempt to establish connection failed: %s", err)
		return nil, err
	}
	return db, nil
}
