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

var ch chan eventinfrastructure.Event

type elkReporter struct {
}

func (e *elkReporter) Write(event eventinfrastructure.Event) {
	ch <- event
}

func (s *elkReporter) SetOutChan(chan<- eventinfrastructure.Event) {
}

func GetReporter() common.Reporter {
	//Do whatever initialization is necessary
	ch = make(chan eventinfrastructure.Event, 100)

	go ListenAndWrite()

	return &elkReporter{}
}

func ListenAndWrite() {
	for {
		select {
		case event, ok := <-ch:
			if ok {
				log.Printf("Sending event to : %v", os.Getenv("ELASTIC_API_EVENTS"))

				b, err := json.Marshal(event)
				if err != nil {
					log.Printf("error: %v", err.Error())
				}

				resp, err := http.Post(os.Getenv("ELASTIC_API_EVENTS"),
					"application/json",
					bytes.NewBuffer(b))

				if err != nil {
					log.Printf("error: %v", err.Error())
				}

				log.Printf("Response %+v", resp)

			} else {
				log.Fatal("Write chan closed. (elk reporter)")
			}
		}
	}
}
