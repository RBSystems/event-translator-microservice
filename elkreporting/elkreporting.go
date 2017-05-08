package elkreporting

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-translator-microservice/common"
)

type elkReporter struct{}

func (e *elkReporter) Write(event eventinfrastructure.Event) {
	log.Printf("Sending event to : %v", os.Getenv("ELASTIC_API_EVENTS"))

	b, err := json.Marshal(event)
	if err != nil {
		log.Printf("error: %v", err.Error())
	}

	_, err = http.Post(os.Getenv("ELASTIC_API_EVENTS"),
		"application/json",
		bytes.NewBuffer(b))

	if err != nil {
		log.Printf("error: %v", err.Error())
	}
}

func (s *elkReporter) SetOutChan(chan<- eventinfrastructure.Event) {
}

func GetReporter() common.Reporter {
	//Do whatever initialization is necessary

	return &elkReporter{}
}
