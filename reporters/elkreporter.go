package reporters

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/v2/events"
)

// elkReporter .
type elkReporter struct {
	SendTimeout time.Duration
	ELKAddress  string

	writeChan chan events.Event
}

// GetELKReporter .
func GetELKReporter() Reporter {
	address := os.Getenv("STATE_PARSER_ADDRESS")

	reporter := &elkReporter{
		SendTimeout: 500 * time.Millisecond,
		ELKAddress:  address,
		writeChan:   make(chan events.Event, 1000),
	}

	go func() {
		for event := range reporter.writeChan {
			err := SendElkEvent(reporter.ELKAddress, event, reporter.SendTimeout)
			if err != nil {
				log.L.Warnf("unable to send event to elk: %v. Event: %+v", err.Error(), event)
			}
		}

		log.L.Fatalf("elk write chan closed unexpectedly")
	}()

	return reporter
}

// Write .
func (e *elkReporter) Write(event events.Event) {
	e.writeChan <- event
}

// SetOutChan .
func (e *elkReporter) SetOutChan(chan<- events.Event) {
}

// SendElkEvent .
func SendElkEvent(address string, event events.Event, timeout time.Duration) *nerr.E {
	b, err := json.Marshal(event)
	if err != nil {
		return nerr.Translate(err).Addf("unable to send event to elk")
	}

	log.L.Infof("[ELKReporter] Sending event to %v", address)
	var client = &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Post(address, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nerr.Translate(err).Addf("unable to send event to elk")
	}
	defer resp.Body.Close()

	val, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nerr.Translate(err).Addf("unable to send event to elk")
	}

	log.L.Debugf("[ELKReporting] Response %s", val)
	return nil
}
