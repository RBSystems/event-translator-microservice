package translator

import (
	"github.com/byuoitav/common/events"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/event-translator-microservice/awsshadowreporting"
	"github.com/byuoitav/event-translator-microservice/common"
	"github.com/byuoitav/event-translator-microservice/elkreporting"
	"github.com/byuoitav/event-translator-microservice/saltreporting"
)

var initialized = false

var reporters = []common.Reporter{}

const queueSize = 1000

var retryCount = 60

/*
Get Reporters returns a list of the reporters
*/
func GetReporters() []common.Reporter {
	if !initialized {
		reporters = append(reporters, awsshadowreporting.GetReporter())
		reporters = append(reporters, saltreporting.GetReporter())
		reporters = append(reporters, elkreporting.GetReporter())
		initialized = true
	}
	return reporters
}

func StartTranslator(en *events.EventNode) error {
	log.L.Infof("Starting translator")
	writeChan := make(chan events.Event, queueSize)

	reporters := GetReporters()

	//Set the write channel for all of the reporters
	for _, currentReporter := range reporters {
		log.L.Infof("Starting reporter")
		currentReporter.SetOutChan(writeChan)
	}

	// start publihser, wait for events to come into writeChan
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
		//write the events
		for i := range reporters {
			reporters[i].Write(msg)
		}
	}
}
