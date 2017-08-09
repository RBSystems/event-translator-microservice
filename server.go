package main

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/device-monitoring-microservice/microservicestatus"
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

	server.GET("/mstatus", GetStatus)
	server.POST("/subscribe", Subscribe, BindSubscriber(sub))
	server.Start(":6998")
}

func Subscribe(context echo.Context) error {
	var cr eventinfrastructure.ConnectionRequest
	context.Bind(&cr)

	s := context.Get(eventinfrastructure.ContextSubscriber)
	if sub, ok := s.(*eventinfrastructure.Subscriber); ok {
		err := eventinfrastructure.HandleSubscriptionRequest(cr, sub)
		if err != nil {
			return context.JSON(http.StatusBadRequest, err.Error())
		}
	}
	return context.JSON(http.StatusOK, nil)

}

func BindSubscriber(s *eventinfrastructure.Subscriber) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(eventinfrastructure.ContextSubscriber, s)
			return next(c)
		}
	}
}

func GetStatus(context echo.Context) error {
	var s microservicestatus.Status
	var err error
	s.Version, err = microservicestatus.GetVersion("version.txt")
	if err != nil {
		s.Version = "missing"
		s.Status = microservicestatus.StatusSick
		s.StatusInfo = fmt.Sprintf("Error: %s", err.Error())
	} else {
		s.Status = microservicestatus.StatusOK
		s.StatusInfo = ""
	}

	return context.JSON(http.StatusOK, s)
}
