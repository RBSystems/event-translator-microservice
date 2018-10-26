package reporters

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
)

var ch chan events.Event

var devch chan events.Event
var prdch chan events.Event

// ELKReporter .
type ELKReporter struct {
}

// Write .
func (e *ELKReporter) Write(event events.Event) {
	ch <- event
}

// SetOutChan .
func (e *ELKReporter) SetOutChan(chan<- events.Event) {
}

/*
// GetELKReporter
func GetELKReporter() Reporter {
	//Do whatever initialization is necessary
	ch = make(chan events.Event, 1000)
	devch = make(chan events.Event, 1000)
	prdch = make(chan events.Event, 1000)

	go ListenAndWrite()
	go ListenAndWriteCh(devch, os.Getenv("ELASTIC_API_EVENTS_DEV"), 250*time.Millisecond)
	go ListenAndWriteCh(prdch, os.Getenv("ELASTIC_API_EVENTS"), 500*time.Millisecond)

	return &elkReporter{}
}
*/

// SendElkEvent .
func SendElkEvent(address string, event events.Event, timeout time.Duration) error {
	b, err := json.Marshal(event)
	if err != nil {
		log.L.Infof("[ELKReporting]error: %v", err.Error())
	}
	var client = &http.Client{
		Timeout: timeout,
	}

	log.L.Infof("[ELKReporting]Sending event to : %v", address)
	resp, err := client.Post(address,
		"application/json",
		bytes.NewBuffer(b))

	if err != nil {
		log.L.Infof("[ELKReporting]error: %v", err.Error())
		return err
	}

	defer resp.Body.Close()

	val, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.L.Infof("[ELKReporting]error: %v", err.Error())
		return err
	}
	log.L.Infof("[ELKReporting]Response %s", val)
	return nil
}

// ListenAndWriteCh .
func ListenAndWriteCh(ch chan events.Event, addr string, timeout time.Duration) {
	for {
		event, ok := <-ch
		if ok {
			SendElkEvent(addr, event, timeout)
		} else {
			log.L.Fatal("[ELKReporting] propagation chan closed. (elk reporter)")
		}
	}
}

// ListenAndWrite .
func ListenAndWrite() {
	for {
		select {
		case event, ok := <-ch:
			log.L.Infof("[ELKReporting] Event received, distributing to distribution channels")
			if ok {

				if len(os.Getenv("ELASTIC_API_EVENTS")) > 0 {
					prdch <- event
				}

				if len(os.Getenv("ELASTIC_API_EVENTS_DEV")) > 0 {
					devch <- event
				}
			} else {
				log.L.Fatal("[ELKReporting]Write chan closed. (elk reporter)")
			}
		}
	}
}
