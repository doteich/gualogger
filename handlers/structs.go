package handlers

import (
	"context"
	"time"
)

type Payload struct {
	Value    interface{}
	TS       time.Time
	Name     string
	Id       string
	Datatype string
	Server   string
}

type Exporter interface {
	Initialize(ctx context.Context) error
	Publish(ctx context.Context, p Payload) error
	Shutdown(ctx context.Context) error
}
