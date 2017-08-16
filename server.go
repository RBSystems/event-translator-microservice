package main

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/device-monitoring-microservice/statusinfrastructure"
	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-translator-microservice/translator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	en := eventinfrastructure.NewEventNode("Translator", "7002", []string{eventinfrastructure.Translator})
	go translator.StartTranslator(en)

	server := echo.New()
	server.Pre(middleware.RemoveTrailingSlash())
	server.Use(middleware.CORS())

	server.GET("/mstatus", GetStatus)
	server.POST("/subscribe", Subscribe, BindEventNode(en))
	server.Start(":6998")
}

func Subscribe(context echo.Context) error {
	var cr eventinfrastructure.ConnectionRequest
	context.Bind(&cr)

	e := context.Get(eventinfrastructure.ContextEventNode)
	if en, ok := e.(*eventinfrastructure.EventNode); ok {
		err := eventinfrastructure.HandleSubscriptionRequest(cr, en)
		if err != nil {
			return context.JSON(http.StatusBadRequest, err.Error())
		}
	}
	return context.JSON(http.StatusOK, nil)
}

func BindEventNode(en *eventinfrastructure.EventNode) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(eventinfrastructure.ContextEventNode, en)
			return next(c)
		}
	}
}

func GetStatus(context echo.Context) error {
	var s statusinfrastructure.Status
	var err error

	s.Version, err = statusinfrastructure.GetVersion("version.txt")
	if err != nil {
		s.Version = "missing"
		s.Status = statusinfrastructure.StatusSick
		s.StatusInfo = fmt.Sprintf("Error: %s", err.Error())
	} else {
		s.Status = statusinfrastructure.StatusOK
		s.StatusInfo = ""
	}

	return context.JSON(http.StatusOK, s)
}
