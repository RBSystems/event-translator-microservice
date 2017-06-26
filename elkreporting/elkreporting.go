package elkreporting

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

type ElkEvent struct {
	eventinfrastructure.Event
	EventCauseString string `json:"event-cause-string"`
	EventTypeString  string `json:"event-type-string"`
}

func ListenAndWrite() {
	for {
		select {
		case event, ok := <-ch:
			if ok {
				//translate the enums to have string types
				newEvent := ElkEvent{event, event.Event.EventCause.String(), event.Event.Type.String()}

				b, err := json.Marshal(newEvent)
				if err != nil {
					log.Printf("[ELKReporting]error: %v", err.Error())
				}

				if len(os.Getenv("ELASTIC_API_EVENTS")) > 0 {
					log.Printf("[ELKReporting]Sending event to : %v", os.Getenv("ELASTIC_API_EVENTS"))

					resp, err := http.Post(os.Getenv("ELASTIC_API_EVENTS"),
						"application/json",
						bytes.NewBuffer(b))

					if err != nil {
						log.Printf("[ELKReporting]error: %v", err.Error())
						continue
					}

					val, err := ioutil.ReadAll(resp.Body)
					log.Printf("[ELKReporting]Response %s", val)
				}

				if len(os.Getenv("ELASTIC_API_EVENTS_DEV")) > 0 {
					log.Printf("[ELKReporting]Sending event to : %v", os.Getenv("ELASTIC_API_EVENTS_DEV"))

					resp, err := http.Post(os.Getenv("ELASTIC_API_EVENTS_DEV"),
						"application/json",
						bytes.NewBuffer(b))

					if err != nil {
						log.Printf("[ELKReporting]error: %v", err.Error())
						continue
					}

					val, err := ioutil.ReadAll(resp.Body)
					log.Printf("[ELKReporting]Response %s", val)
				}
			} else {
				log.Fatal("[ELKReporting]Write chan closed. (elk reporter)")
			}
		}
	}
}
