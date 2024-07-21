package main

import (
	"context"
	"fmt"

	"os"
)

var conf Configuration

func init() {

	l := os.Getenv("GOPC_LOG_LEVEL")

	InitLogger(l)

	var err error

	conf, err = LoadConfig()

	if err != nil {
		Logger.Error(fmt.Sprintf("error while loading configuration: %s", err.Error()), "func", "init")
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()

	conf.Opcua.InitSuperVisor(ctx)
}
