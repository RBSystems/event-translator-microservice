package elkreporting

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/byuoitav/common/events"
	"github.com/byuoitav/event-translator-microservice/common"
)

var ch chan events.Event

var devch chan events.Event
var prdch chan events.Event

type elkReporter struct {
}

func (e *elkReporter) Write(event events.Event) {
	ch <- event
}

func (s *elkReporter) SetOutChan(chan<- events.Event) {
}

func GetReporter() common.Reporter {
	//Do whatever initialization is necessary
	ch = make(chan events.Event, 1000)
	devch = make(chan events.Event, 1000)
	prdch = make(chan events.Event, 1000)

	go ListenAndWrite()
	go ListenAndWriteCh(devch, os.Getenv("ELASTIC_API_EVENTS_DEV"), 250*time.Millisecond)
	go ListenAndWriteCh(prdch, os.Getenv("ELASTIC_API_EVENTS"), 500*time.Millisecond)

	return &elkReporter{}
}

type ElkEvent struct {
	events.Event
	EventCauseString string `json:"event-cause-string"`
	EventTypeString  string `json:"event-type-string"`
}

func SendElkEvent(address string, event events.Event, timeout time.Duration) error {
	newEvent := ElkEvent{event, event.Event.EventCause.String(), event.Event.Type.String()}
	b, err := json.Marshal(newEvent)
	if err != nil {
		log.Printf("[ELKReporting]error: %v", err.Error())
	}
	var client = &http.Client{
		Timeout: timeout,
	}

	log.Printf("[ELKReporting]Sending event to : %v", address)
	resp, err := client.Post(address,
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

func ListenAndWriteCh(ch chan events.Event, addr string, timeout time.Duration) {
	for {
		event, ok := <-ch
		if ok {
			SendElkEvent(addr, event, timeout)
		} else {
			log.Fatal("[ELKReporting] propagation chan closed. (elk reporter)")
		}
	}
}

func ListenAndWrite() {
	for {
		select {
		case event, ok := <-ch:
			if ok {

				if len(os.Getenv("ELASTIC_API_EVENTS")) > 0 {
					prdch <- event
				}

				if len(os.Getenv("ELASTIC_API_EVENTS_DEV")) > 0 {
					devch <- event
				}
			} else {
				log.Fatal("[ELKReporting]Write chan closed. (elk reporter)")
			}
		}
	}
}
