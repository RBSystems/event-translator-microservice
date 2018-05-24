package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/byuoitav/common/events"
	"github.com/byuoitav/device-monitoring-microservice/statusinfrastructure"
	"github.com/byuoitav/event-translator-microservice/translator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	en := events.NewEventNode("Translator", []string{events.Translator}, os.Getenv("EVENT_ROUTER_ADDRESS"))
	go translator.StartTranslator(en)

	server := echo.New()
	server.Pre(middleware.RemoveTrailingSlash())
	server.Use(middleware.CORS())

	server.GET("/mstatus", GetStatus)
	server.Start(":6998")
}

func BindEventNode(en *events.EventNode) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(events.ContextEventNode, en)
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
