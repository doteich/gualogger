package handlers

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TimeScaleDB struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Table    string `mapstructure:"table"`
	Pool     *pgxpool.Pool
}

func (t *TimeScaleDB) Initialize(ctx context.Context, cb func(context.Context) []Payload) error {
	cs := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", t.Username, t.Password, t.Host, t.Port, t.Database)

	var err error

	t.Pool, err = pgxpool.New(ctx, cs)

	if err != nil {
		return err
	}

	if err := t.Pool.Ping(ctx); err != nil {
		return err
	}

	sql := `CREATE TABLE IF NOT EXISTS ` + t.Table

	sql += ` (
		value    TEXT NOT NULL,
		ts       TIMESTAMPTZ NOT NULL,
		name     TEXT NOT NULL,
		id       TEXT NOT NULL,
		datatype TEXT NOT NULL,
		server   TEXT NOT NULL,
		meta 	JSONB NOT NULL DEFAULT '{}'::JSONB
		);`

	_, err = t.Pool.Exec(ctx, sql)

	if err != nil {
		fmt.Println("Error creating tabel:", err)
		return err
	}

	sql = fmt.Sprintf("SELECT create_hypertable('%s', by_range('ts'), if_not_exists => TRUE)", t.Table)

	_, err = t.Pool.Exec(ctx, sql)

	if err != nil {
		fmt.Println("Error creating hypertable:", err)
		return err
	}

	return nil
}

// func SetupTable() error {

// }

func (t *TimeScaleDB) Publish(ctx context.Context, p Payload) error {

	sql := fmt.Sprintf("INSERT INTO %s (value, ts, name, id, datatype, server, meta) VALUES ($1, $2, $3, $4, $5, $6, $7)", t.Table)

	args := []any{fmt.Sprint(p.Value), p.TS, p.Name, p.Id, p.Datatype, p.Server, p.Meta}

	_, err := t.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (t *TimeScaleDB) Shutdown(ctx context.Context) error {
	t.Pool.Close()
	return nil
}
