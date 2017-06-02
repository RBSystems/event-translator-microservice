package handlers

import (
	"log"
	"net/http"

	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-router-microservice/subscription"
	"github.com/byuoitav/event-translator-microservice/translator"
	"github.com/labstack/echo"
)

func Subscribe(context echo.Context) error {
	var sr subscription.SubscribeRequest
	err := context.Bind(&sr)
	if err != nil {
		log.Printf("[error] %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	log.Printf("Subscribing to %s", sr.Address)
	err = translator.Sub.Subscribe(sr.Address, []string{eventinfrastructure.Translator})
	if err != nil {
		log.Printf("[error] %s", err.Error())
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, context)
}
