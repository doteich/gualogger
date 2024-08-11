package main

import (
	"context"
	"gualogger/handlers"
)

type ExportManager struct {
	exporters []handlers.Exporter
}

// Initializes a new manager instance
func NewManager(e *Exporters) *ExportManager {
	m := new(ExportManager)
	m.exporters = make([]handlers.Exporter, 0)
	m.RegisterExporters(e)
	return m
}

// Adds exporters to the manager based on entries in the exporter config structure
func (m *ExportManager) RegisterExporters(e *Exporters) {
	for _, exp := range e.RegisteredExporters {
		switch exp {
		case "timescale-db":
			m.exporters = append(m.exporters, &e.TimeScaleDB)

		}
	}
}

// Setup exporter by calling the Initialize() function of each exporters interface
// If the initialization of one exporter fails, the first error gets returned
func (m *ExportManager) SetupPubHandlers(ctx context.Context) error {

	for _, e := range m.exporters {
		if err := e.Initialize(ctx); err != nil {
			return err
		}
	}

	return nil
}

func Publish() {

}
