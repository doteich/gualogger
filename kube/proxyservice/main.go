package main

import (
	"flag"
	"proxyservice/routes"

	"github.com/labstack/echo/v4"
)

var (
	address *string
)

func init() {
	address = flag.String("a", ":3001", "address & port for the server")
	flag.Parse()

}

func main() {

	e := echo.New()

	e.GET("/crds", routes.GetCRDs)

	e.Logger.Fatal(e.Start(*address))

}
