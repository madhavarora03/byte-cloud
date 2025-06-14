package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/madhavarora03/byte-cloud/internal/config"
)

type Postgres struct {
	Conn *pgx.Conn
}

func Init(cfg *config.Config) (*Postgres, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DbUri)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		Conn: conn,
	}, nil
}
