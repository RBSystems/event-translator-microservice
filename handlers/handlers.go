package handlers

import (
	"log"

	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-router-microservice/subscription"
	"github.com/byuoitav/event-translator-microservice/translator"
	"github.com/labstack/echo"
)

func Subscribe(context echo.Context) {
	var sr subscription.SubscribeRequest
	context.Bind(&sr)
	log.Printf("Subscribing to %s", sr.Address)
	translator.Sub.Subscribe(sr.Address, []string{eventinfrastructure.Translator})
}
