package common

import (
	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-translator-microservice/awsshadowreporting"
	"github.com/byuoitav/event-translator-microservice/elkreporting"
	"github.com/byuoitav/event-translator-microservice/saltreporting"
)

var initialized = false

var reporters = []*Reporter{}

/*
Reporter is used to translate events internal to the system to external event systems, and vice versa.

Write() will be called each time an event is recieved from the internal router.
SetOutChan() will be called on each reporter, and events placed in this channel will be echoed to the local router.

*/
type Reporter interface {
	Write(eventinfrastructure.Event)
	SetOutChan(chan<- eventinfrastructure.Event)
}

/*
Get Reporters returns a list of the reporters
*/
func GetReporters() []Reporter {
	if !initialized {
		reporters = append(reporters, awsshadowreporting.GetReporter())
		reporters = append(reporters, saltreporting.GetReporter())
		reporters = append(reporters, elkreporting.GetReporter())

	}
	return reporters
}
