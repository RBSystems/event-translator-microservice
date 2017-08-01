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

	server.POST("/subscribe", Subscribe, BindSubscriber(sub))
	server.Start(":6998")
}

func Subscribe(context echo.Context) error {
	return eventinfrastructure.HandleSubscriptionRequest(context)
}

func BindSubscriber(s *eventinfrastructure.Subscriber) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(eventinfrastructure.ContextSubscriber, s)
			return next(c)
		}
	}
}
