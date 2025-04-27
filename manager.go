package main

import (
	"context"
	"fmt"
	"gualogger/handlers"
	"gualogger/logging"
)

type ExportManager struct {
	exporters   map[string]handlers.Exporter
	Meta        map[string][]handlers.Meta
	IncludeMeta bool
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
			logging.Logger.Info(fmt.Sprintf("registered exporter: %s", k), "func", "RegisterExporters")
			m.exporters[k] = h
		}
	}
}

func (m *ExportManager) BuildMeta(arr []Nodeid) {
	m.Meta = make(map[string][]handlers.Meta, 0)
	hasMeta := false

	for _, n := range arr {
		meta := make([]handlers.Meta, 0)

		if len(n.Meta) > 0 {
			hasMeta = true
		}

		for _, met := range n.Meta {
			meta = append(meta, handlers.Meta{Key: met.Key, Value: met.Value})
		}
		m.Meta[n.Id] = meta
	}
	m.IncludeMeta = hasMeta

}

// Setup exporter by calling the Initialize() function of each exporters interface
// If the initialization of one exporter fails, the first error gets returned
func (m *ExportManager) SetupPubHandlers(ctx context.Context) error {

	callback := func(c context.Context) []handlers.Payload {
		return Read(c)
	}

	for n, e := range m.exporters {
		if err := e.Initialize(ctx, callback); err != nil {
			return fmt.Errorf("error while initializing exporter %s - %s", n, err.Error())

		}
		logging.Logger.Info(fmt.Sprintf("successfully initialized exporter: %s", n), "func", "SetupPubHandlers")

	}
	return nil
}

func (m *ExportManager) Publish(ctx context.Context, p handlers.Payload) {
	if m.IncludeMeta {
		if meta, exists := m.Meta[p.Id]; exists {
			p.Meta = meta
		}
	}
	for n, e := range m.exporters {
		p.Server = conf.Opcua.Connection.Endpoint

		if err := e.Publish(ctx, p); err != nil {
			logging.Logger.Error(fmt.Sprintf("failed to publish payload for exporter %s: %s", n, err.Error()), "func", "Publish")
		}
	}
}
