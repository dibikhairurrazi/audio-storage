package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"

	"cloud.google.com/go/cloudsqlconn/postgres/pgxv4"
	"github.com/dibikhairurrazi/audio-storage/config"
)

//go:embed migration/*.sql
var migrationsFS embed.FS

type DB struct {
	MasterConn  *sql.DB
	ReplicaConn *sql.DB
}

func Initialize(cfg *config.DatabaseConfig) (*DB, error) {
	driverName := "postgres"
	if cfg.UseCloudSQL {
		driverName = "cloudsql-postgres"
		_, err := pgxv4.RegisterDriver(
			driverName,
		)

		if err != nil {
			return nil, err
		}
	}

	masterDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Master.Host, cfg.Master.Port, cfg.Master.User, cfg.Master.Password, cfg.Master.DBName)
	masterConn, err := sql.Open(
		driverName,
		masterDSN,
	)

	slog.Info("connection to DB with DSN", "dsn", masterDSN)

	if err != nil {
		slog.Error("error connecting to master DB", "err", err.Error())
		return nil, err
	}

	replicaDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Replica.Host, cfg.Replica.Port, cfg.Replica.User, cfg.Replica.Password, cfg.Replica.DBName)
	replicaConn, err := sql.Open(
		driverName,
		replicaDSN,
	)

	slog.Info("connection to DB with DSN", "dsn", replicaDSN)

	if err != nil {
		slog.Error("error connecting to replica DB", "err", err.Error())
		return nil, err
	}

	err = masterConn.Ping()
	if err != nil {
		slog.Error("error ping-ing master DB", "err", err.Error())
		return nil, err
	}

	err = replicaConn.Ping()
	if err != nil {
		slog.Error("error ping-ing replica DB", "err", err.Error())
		return nil, err
	}

	return &DB{
		MasterConn:  masterConn,
		ReplicaConn: replicaConn,
	}, nil
}

func MigrateUp(client *sql.DB, migrationFolder, dbName string) error {
	m, err := prepareMigration(client, migrationFolder, dbName)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func MigrateDown(client *sql.DB, migrationFolder, dbName string) error {
	m, err := prepareMigration(client, migrationFolder, dbName)
	if err != nil {
		return err
	}

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func MigrateDrop(client *sql.DB, migrationFolder, dbName string) error {
	m, err := prepareMigration(client, migrationFolder, dbName)
	if err != nil {
		return err
	}

	err = m.Drop()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func MigrateSteps(client *sql.DB, migrationFolder, dbName string, steps int) error {
	m, err := prepareMigration(client, migrationFolder, dbName)
	if err != nil {
		return err
	}

	err = m.Steps(steps)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func prepareMigration(client *sql.DB, migrationFolder, dbName string) (*migrate.Migrate, error) {
	dir, err := iofs.New(migrationsFS, migrationFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded FS: %w", err)
	}

	driver, err := postgres.WithInstance(client, &postgres.Config{DatabaseName: dbName})
	if err != nil {
		return nil, fmt.Errorf("failed to get migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", dir, dbName, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

func SeedDB(client *sql.DB) error {
	seedFiles, err := os.ReadDir("db/seed/")
	if err != nil {
		return err
	}

	for _, f := range seedFiles {
		c, err := os.ReadFile("db/seed/" + f.Name())
		if err != nil {
			return err
		}

		sqlCode := string(c)

		_, err = client.Exec(sqlCode)
		if err != nil {
			return err
		}
	}

	return nil
}
