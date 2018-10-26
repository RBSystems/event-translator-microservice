package main

import (
	"os"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/events"
	"github.com/byuoitav/event-translator-microservice/translator"
)

func main() {
	en := events.NewEventNode("Translator", os.Getenv("EVENT_ROUTER_ADDRESS"), []string{events.Translator})
	go translator.StartTranslator(en)

	server := common.NewRouter()
	server.Start(":6998")
}
