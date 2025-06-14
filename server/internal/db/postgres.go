package db

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/atomic"
	"time"
)

type Config struct {
	MigrationsTable       string
	DatabaseName          string
	SchemaName            string
	migrationsSchemaName  string
	migrationsTableName   string
	StatementTimeout      time.Duration
	MigrationsTableQuoted bool
	MultiStatementEnabled bool
	MultiStatementMaxSize int
}

type Postgres struct {
	// Locking and unlocking need to use the same connection
	conn     *sql.Conn
	Db       *sql.DB
	isLocked atomic.Bool

	// Open and WithInstance need to guarantee that config is never nil
	config *Config
}

func NewPostgres(ctx context.Context, cfg *Config, dbURI string) (*Postgres, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Use dbURI directly with pgx driver
	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set sensible connection settings
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.StatementTimeout)
	defer cancel()

	if err := db.PingContext(ctxWithTimeout); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	return &Postgres{
		conn:   conn,
		Db:     db,
		config: cfg,
	}, nil
}
