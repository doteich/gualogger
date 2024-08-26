package main

import (
	"context"
	"fmt"
	"gualogger/handlers"
)

type ExportManager struct {
	exporters map[string]handlers.Exporter
}

// Initializes a new manager instance
func NewManager(e *Exporters, emap *map[string]interface{}) *ExportManager {
	m := new(ExportManager)
	m.exporters = make(map[string]handlers.Exporter, 0)
	m.RegisterExporters(e, emap)
	return m
}

// Adds exporters to the manager based on entries in the exporter config structure
func (m *ExportManager) RegisterExporters(e *Exporters, emap *map[string]interface{}) {
	reg := e.GetExporterRegister()
	for k := range *emap {
		h, exists := reg[k]
		if exists {
			Logger.Info(fmt.Sprintf("registered exporter: %s", k), "func", "RegisterExporters")
			m.exporters[k] = h
		}
	}
}

// Setup exporter by calling the Initialize() function of each exporters interface
// If the initialization of one exporter fails, the first error gets returned
func (m *ExportManager) SetupPubHandlers(ctx context.Context) error {
	for n, e := range m.exporters {
		if err := e.Initialize(ctx); err != nil {
			return fmt.Errorf("error while initializing exporter %s - %s", n, err.Error())

		}
		Logger.Info(fmt.Sprintf("successfully initialized exporter: %s", n), "func", "SetupPubHandlers")

	}
	return nil
}

func (m *ExportManager) Publish(ctx context.Context, p handlers.Payload) {
	for n, e := range m.exporters {
		p.Server = conf.Opcua.Connection.Endpoint

		if err := e.Publish(ctx, p); err != nil {
			Logger.Error(fmt.Sprintf("failed to publish payload for exporter %s: %s", n, err.Error()), "func", "Publish")
		}
	}
}
