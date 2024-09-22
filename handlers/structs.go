package handlers

import (
	"context"
	"time"
)

type Payload struct {
	Value    interface{} `json:"value"`
	TS       time.Time   `json:"ts"`
	Name     string      `json:"name"`
	Id       string      `json:"id"`
	Datatype string      `json:"datatype"`
	Server   string      `json:"server"`
}

type Exporter interface {
	Initialize(ctx context.Context, callback func(context.Context) []Payload) error
	Publish(ctx context.Context, p Payload) error
	Shutdown(ctx context.Context) error
}
