package main

import (
	"context"
	"gualogger/handlers"
	"gualogger/logging"
	"time"
)

type ExportManager struct {
	redpandaInstance *handlers.Redpanda
}

// Initializes a new manager instance
func NewManager(r *handlers.Redpanda) *ExportManager {
	m := new(ExportManager)
	m.redpandaInstance = r
	return m
}

// Setup exporter by calling the Initialize() function of each exporters interface
// If the initialization of one exporter fails, the first error gets returned
func (m *ExportManager) SetupPubHandler(ctx context.Context) error {

	if err := m.redpandaInstance.Initialize(ctx); err != nil {
		return err
	}

	logging.Logger.Info("successfully connected to redpanda brokers")
	return nil
}

func (m *ExportManager) Publish(ctx context.Context, p handlers.Payload) {
	m.redpandaInstance.Publish(ctx, p)
}

func (m *ExportManager) VerifyConnection(ctx context.Context) {
	for {

		err := m.redpandaInstance.Ping(ctx)

		if err != nil {
			logging.Logger.Warn("unable to ping")

			m.SetupPubHandler(ctx)

		}

		time.Sleep(60 * time.Second)

	}

}
