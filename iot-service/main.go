package main

import (
	"context"
	"fmt"
	"gualogger/logging"
	"os"
)

var (
	conf *Configuration
	mgr  *ExportManager
)

func init() {

	l := os.Getenv("GOPC_LOG_LEVEL")

	logging.InitLogger(l)

	var err error

	conf, err = LoadConfig()

	if err != nil {
		logging.Logger.Error(fmt.Sprintf("error while loading configuration: %s", err.Error()), "func", "init")
		os.Exit(1)
	}

}

func main() {
	ctx := context.Background()

	mgr = NewManager(&conf.Exporters, &conf.ExpMap)

	mgr.BuildMeta(conf.Opcua.Subscription.Nodeids)

	if err := mgr.SetupPubHandlers(ctx); err != nil {
		logging.Logger.Error(err.Error(), "func", "main")
		return
	}

	conf.Opcua.InitSuperVisor(ctx)
}
