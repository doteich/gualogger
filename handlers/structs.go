package handlers

import "time"

type payload struct {
	value    interface{}
	ts       time.Time
	name     string
	id       string
	datatype string
	server   string
}
