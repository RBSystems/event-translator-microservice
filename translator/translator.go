package translator

import (
	"encoding/json"
	"log"

	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
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

func StartTranslator(publisherPort string, sub *eventinfrastructure.Subscriber) error {
	log.Printf("Starting trasnlator")
	writeChan := make(chan eventinfrastructure.Event, queueSize)

	reporters := GetReporters()

	//Set the write channel for all of the reporters
	for _, currentReporter := range reporters {
		currentReporter.SetOutChan(writeChan)
	}

	// start publihser, wait for events to come into writeChan
	pub := eventinfrastructure.NewPublisher(publisherPort)
	go func() {
		for {
			select {
			case event, ok := <-writeChan:
				if ok {
					pub.PublishEvent(event, eventinfrastructure.External)
				} else {
					log.Fatal("[Publisher] Write chan closed.")
				}
			}
		}
	}()

	// start subscriber
	// create connection with router
	var req eventinfrastructure.ConnectionRequest
	req.PublisherAddr = "localhost:" + publisherPort
	req.SubscriberEndpoint = "http://localhost:6998/subscribe"

	go eventinfrastructure.SendConnectionRequest("http://localhost:6999/subscribe", req, true)

	// listen to events and echo them out
	for {
		select {
		case message, ok := <-sub.MessageChan:
			if !ok {
				log.Fatalf("[error] subscriber read channel closed")
			}

			var event eventinfrastructure.Event
			err := json.Unmarshal(message.MessageBody, &event)
			if err != nil {
				log.Printf("[error] there was a problem decoding a message to an event: %s", err.Error())
				continue
			}

			//write the events
			for i := range reporters {
				reporters[i].Write(event)
			}
		}
	}
}
