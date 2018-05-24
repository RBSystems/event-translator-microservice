package elkreporting

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/byuoitav/common/events"
	"github.com/byuoitav/event-translator-microservice/common"
)

var ch chan events.Event

type elkReporter struct {
}

func (e *elkReporter) Write(event events.Event) {
	ch <- event
}

func (s *elkReporter) SetOutChan(chan<- events.Event) {
}

func GetReporter() common.Reporter {
	//Do whatever initialization is necessary
	ch = make(chan events.Event, 100)

	go ListenAndWrite()

	return &elkReporter{}
}

type ElkEvent struct {
	events.Event
	EventCauseString string `json:"event-cause-string"`
	EventTypeString  string `json:"event-type-string"`
}

func SendElkEvent(address string, event events.Event) error {
	newEvent := ElkEvent{event, event.Event.EventCause.String(), event.Event.Type.String()}
	b, err := json.Marshal(newEvent)
	if err != nil {
		log.Printf("[ELKReporting]error: %v", err.Error())
	}

	log.Printf("[ELKReporting]Sending event to : %v", address)
	resp, err := http.Post(address,
		"application/json",
		bytes.NewBuffer(b))

	if err != nil {
		log.Printf("[ELKReporting]error: %v", err.Error())
		return err
	}

	defer resp.Body.Close()

	val, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ELKReporting]error: %v", err.Error())
		return err
	}
	log.Printf("[ELKReporting]Response %s", val)
	return nil
}

func ListenAndWrite() {
	for {
		select {
		case event, ok := <-ch:
			if ok {

				if len(os.Getenv("ELASTIC_API_EVENTS")) > 0 {
					SendElkEvent(os.Getenv("ELASTIC_API_EVENTS"), event)
				}

				if len(os.Getenv("ELASTIC_API_EVENTS_DEV")) > 0 {
					SendElkEvent(os.Getenv("ELASTIC_API_EVENTS_DEV"), event)
				}
			} else {
				log.Fatal("[ELKReporting]Write chan closed. (elk reporter)")
			}
		}
	}
}
