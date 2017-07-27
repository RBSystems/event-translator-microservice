package main

import (
	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-translator-microservice/translator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	sub := eventinfrastructure.NewSubscriber([]string{eventinfrastructure.Translator})
	port := "7002"
	go translator.StartTranslator(port, sub)

	server := echo.New()
	server.Pre(middleware.RemoveTrailingSlash())
	server.Use(middleware.CORS())

	server.POST("/subscribe", sub.HandleSubscriptionRequest)
	server.Start(":6998")
}
