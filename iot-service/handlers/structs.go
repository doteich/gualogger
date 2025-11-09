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
	Meta     []Meta      `json:"meta"`
}

type Meta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Exporter interface {
	Initialize(ctx context.Context) error
	Publish(ctx context.Context, p Payload) error
	Shutdown(ctx context.Context) error
}
