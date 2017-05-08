package translator

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
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

const QueueSize = 1000

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
	sub, err := subscriber.NewSubscriber(queueSize)
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR: Could not create subscriber: %s", err.Error()))
		return err
	}
	for retryCount > -1 {
		err = sub.Subscribe(routerAddr, []string{eventinfrastructure.Translator})
		if err != nil {
			if retryCount > 0 { //retry
				retryCount--
				log.Printf("Susbcription to router failed with error %s,  will try again %v times", err.Error(), retryCount)

				timer := time.NewTimer(2 * time.Second)
				<-timer.C

				log.Printf("Retrying subscription to router at %s", routerAddr)
				continue
			} else {
				log.Fatal(fmt.Sprintf("ERROR: Could not subscribe to router %s, exceeded retry attempts error was: %s", routerAddr, err.Error()))
				return err
			}
		} else {
			break
		}
	}

	// listen to events and echo them out
	for {

		msg := sub.Read()

		var event eventinfrastructure.Event

		err := json.Unmarshal(msg.MessageBody, &event)
		if err != nil {
			log.Printf("ERROR: Ther was a problem decoding a message to an event : %s", err.Error())
			continue
		}

		//write the events
		for i := range reporters {
			reporters[i].Write(event)
		}
	}

	return nil
}
