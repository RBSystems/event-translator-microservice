package saltreporting

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/byuoitav/common/events"
	"github.com/byuoitav/event-translator-microservice/common"
)

var eventBuffer chan events.Event

type saltReporter struct {
}

func (s *saltReporter) Write(event events.Event) {
	return
}

//no events will be received over the salt bus
func (s *saltReporter) SetOutChan(chan<- events.Event) {
	return
}

func GetReporter() common.Reporter {
	return &saltReporter{}
}

func (s *saltReporter) startWriter(saltEventAddr string) {
	log.Printf("[SaltReporting] Waiting for events.")
	for {
		event, ok := <-eventBuffer
		if ok {
			log.Printf("[SaltReporting] Writing event")

			addr := saltEventAddr + "/event/" + event.Event.Type.String() + "/" + event.Event.EventCause.String()

			b, err := json.Marshal(event)
			if err != nil {
				log.Printf("[SaltReporting] error masrhalling event %v to JSON. ERROR: %v", event, err.Error())
				continue
			}

			_, err = http.Post(addr, "application/json", bytes.NewBuffer(b))
			if err != nil {
				log.Printf("[SaltReporting] Error sending event %v. ERROR: %v", event, err.Error())
			}

		} else {
			log.Printf("[SaltReporting] Event queue closed.")
			return
		}
	}
}
