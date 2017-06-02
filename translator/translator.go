package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-router-microservice/subscription"
	"github.com/byuoitav/event-translator-microservice/awsshadowreporting"
	"github.com/byuoitav/event-translator-microservice/common"
	"github.com/byuoitav/event-translator-microservice/elkreporting"
	"github.com/byuoitav/event-translator-microservice/saltreporting"
	message "github.com/xuther/go-message-router/common"
	"github.com/xuther/go-message-router/publisher"
	"github.com/xuther/go-message-router/subscriber"
)

var initialized = false

var reporters = []common.Reporter{}

const queueSize = 1000

var retryCount = 60

var Sub subscriber.Subscriber

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

func StartTranslator(routerAddr string, publisherPort string) error {
	writeChan := make(chan eventinfrastructure.Event, queueSize)

	reporters := GetReporters()

	//Set the write channel for all of the reporters
	for _, currentReporter := range reporters {
		currentReporter.SetOutChan(writeChan)
	}

	//Start our publisher
	go func() {
		pub, err := publisher.NewPublisher(publisherPort, queueSize, 10)
		if err != nil {
			log.Fatal("Error creating publisher: " + err.Error())
		}
		pub.Listen()
		header := [24]byte{}
		copy(header[:], []byte(eventinfrastructure.External))
		for {
			select {
			case event, ok := <-writeChan:
				if ok {
					b, err := json.Marshal(event)
					if err != nil {
						log.Printf("ERROR marshalling event into event struct: %s", err.Error())
						continue
					}
					pub.Write(message.Message{MessageHeader: header, MessageBody: b})
				} else {
					log.Fatal("Write chan closed.")
				}
			}
		}
	}()

	//start our subscriber
	Sub, err := subscriber.NewSubscriber(queueSize)
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR: Could not create subscriber: %s", err.Error()))
		return err
	}

	// tell router to subscribe to me, and me to it
	var s subscription.SubscribeRequest
	s.Address = "localhost:7002"
	s.PubAddress = "localhost:6998/subscribe"
	body, err := json.Marshal(s)
	if err != nil {
		log.Printf("[error] %s", err)
	}

	_, err = http.Post("http://localhost:6999/subscribe", "application/json", bytes.NewBuffer(body))
	for err != nil {
		log.Printf("[error] failed to connect to the router. Trying again...")
		time.Sleep(3 * time.Second)
		_, err = http.Post("http://localhost:6999/subscribe", "application/json", bytes.NewBuffer(body))
	}
	log.Printf("The event router is subscribed to me.")

	// listen to events and echo them out
	for {

		msg := Sub.Read()

		var event eventinfrastructure.Event

		err := json.Unmarshal(msg.MessageBody, &event)
		if err != nil {
			log.Printf("ERROR: Ther was a problem decoding a message to an event : %s", err.Error())
			continue
		}

		log.Printf("Event recieved: %s", msg.MessageBody)

		//write the events
		for i := range reporters {
			reporters[i].Write(event)
		}
	}
}
