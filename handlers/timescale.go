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
	Pool     *pgxpool.Pool
}

func (t *TimeScaleDB) Initialize(ctx context.Context) error {
	cs := fmt.Sprintf("postgresql://%s:%s@%s:%d", t.Username, t.Password, t.Host, t.Port)

	var err error

	t.Pool, err = pgxpool.New(ctx, cs)

	if err != nil {
		return err
	}

	fmt.Println("successfully connected")

	return nil
}

// func SetupTable() error {

// }

func (t *TimeScaleDB) Publish(ctx context.Context, p payload) error {
	return nil
}

func (t *TimeScaleDB) Shutdown(ctx context.Context) error {
	t.Pool.Close()
	return nil
}
