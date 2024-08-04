package main

type manager struct {
	exporters []func()
}

func (e Exporters) setupPubHandlers() {

	e.TimeScaleDB.CreateConPool()

}

func Publish() {

}
