package main

import (
	"os"

	"github.com/byuoitav/event-translator-microservice/handlers"
	"github.com/byuoitav/event-translator-microservice/translator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := "7002"
	routerAddr := ""
	//Get the router address
	if len(os.Getenv("LOCAL_ENVIRONMENT")) > 0 {
		routerAddr = "localhost:7000"
	}
	go translator.StartTranslator(routerAddr, port)

	server := echo.New()
	server.Pre(middleware.RemoveTrailingSlash())
	server.Use(middleware.CORS())

	server.POST("/subscribe", handlers.Subscribe)
	server.Start(":6998")
}
