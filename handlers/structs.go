package handlers

import (
	"context"
	"time"
)

type payload struct {
	value    interface{}
	ts       time.Time
	name     string
	id       string
	datatype string
	server   string
}

type Exporter interface {
	Initialize(ctx context.Context) error
	Publish(ctx context.Context, p payload) error
	Shutdown(ctx context.Context) error
}
