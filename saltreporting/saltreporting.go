package saltreporting

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-translator-microservice/common"
)

type saltReporter struct {
	eventBuffer chan eventinfrastructure.Event
}

func (s *saltReporter) Write(event eventinfrastructure.Event) {
	s.eventBuffer <- event
	return
}

//no events will be received over the salt bus
func (s *saltReporter) SetOutChan(chan<- eventinfrastructure.Event) {
	return
}

func GetReporter() common.Reporter {
	buf = make(chan eventBuffer, 1000)
	reporter := saltReporter{buf}

	go reporter.startWriter("http://localhost:7010")
	return &saltReporter{}
}

func (s *saltReporter) startWriter(saltEventAddr string) {
	for {
		event, ok := <-s.eventBuffer
		if ok {
			log.Printf("[SaltReporting] Writing event")

			addr := saltEventAddr + "/event/" + event.Event.Type.String() + "/" + event.Event.EventCause.String()

			b, err := json.Marshal(event)
			if err != nil {
				log.Printf("[SaltReporting] error masrhalling event %v to JSON. ERROR: %v", event, err.Error())
				continue
			}

			_, err := http.Get(addr, "application/json", bytes.NewBuffer(b))
			if err != nil {
				log.Printf("[SaltReporting] Error sending event %v. ERROR: %v", event, err.Error())
			}

		} else {
			log.Printf("[SaltReporting] Event queue closed.")
			return
		}
	}
}
