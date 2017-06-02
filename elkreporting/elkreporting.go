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
				log.Printf("[ELKReporting]Sending event to : %v", os.Getenv("ELASTIC_API_EVENTS"))

				b, err := json.Marshal(newEvent)
				if err != nil {
					log.Printf("[ELKReporting]error: %v", err.Error())
				}

				resp, err := http.Post(os.Getenv("ELASTIC_API_EVENTS"),
					"application/json",
					bytes.NewBuffer(b))

				if err != nil {
					continue
					log.Printf("[ELKReporting]error: %v", err.Error())
				}

				val, err := ioutil.ReadAll(resp.Body)

				log.Printf("[ELKReporting]Response %s", val)

				//---------------------------------------------------------------
				//TEMP
				resp, err = http.Post("http://dev-elk-shipper0.byu.edu:5543",
					"application/json",
					bytes.NewBuffer(b))

				if err != nil {
					log.Printf("[ELKReporting]error: %v", err.Error())
					continue
				}

				val, err = ioutil.ReadAll(resp.Body)

				log.Printf("[ELKReporting]Response %s", val)
				//END TEMP
				//---------------------------------------------------------------

			} else {
				log.Fatal("[ELKReporting]Write chan closed. (elk reporter)")
			}
		}
	}
}
