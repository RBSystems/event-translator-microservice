package main

import (
	"os"

	"github.com/byuoitav/event-translator-microservice/translator"
)

func main() {
	port := "7002"
	routerAddr := ""
	//Get the router address
	if len(os.Getenv("LOCAL_ENVIRONMENT")) > 0 {
		routerAddr = "localhost:7000"
	}
	translator.StartTranslator(routerAddr, port)
}
