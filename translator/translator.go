package translator

import (
	"reflect"
	"sync"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/event-translator-microservice/reporters"
)

const (
	queueSize  = 1000
	retryCount = 60
)

var (
	once         sync.Once
	reporterList []reporters.Reporter
)

// GetReporters returns a list of the reporters
func GetReporters() []reporters.Reporter {
	once.Do(func() {
		// TODO decide which ones to use based on local or in aws
		reporterList = append(reporterList, &reporters.ELKReporter{})
	})

	return reporterList
}

// StartTranslator starts the event translator
func StartTranslator(en *events.EventNode) {
	writeChan := make(chan events.Event, queueSize)

	reporterNames := []string{}
	for _, reporter := range GetReporters() {
		reporterNames = append(reporterNames, reflect.TypeOf(reporter).String())
	}

	log.L.Infof("Starting translator with reporters: %s", reporterNames)

	//Set the write channel for all of the reporters
	for _, currentReporter := range GetReporters() {
		currentReporter.SetOutChan(writeChan)
	}

	// start publisher, wait for events to come into writeChan
	go func() {
		for {
			select {
			case event, ok := <-writeChan:
				if ok {
					en.PublishEvent(events.External, event)
				} else {
					log.L.Fatal("[Publisher] Write chan closed.")
				}
			}
		}
	}()

	// listen to events and echo them out
	for {
		msg, err := en.Read()
		if err != nil {
			log.L.Errorf("Error: %v", err.Error())
			continue
		}

		log.L.Debugf("Received event, pushing to reporters")

		// write the events
		for _, reporter := range GetReporters() {
			reporter.Write(msg)
		}
	}
}
