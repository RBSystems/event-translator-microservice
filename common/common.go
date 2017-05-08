package common

import (
	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
)

//This package is so small because of otherwise cyclical dependencies

/*
Reporter is used to translate events internal to the system to external event systems, and vice versa.

Write() will be called each time an event is recieved from the internal router.
SetOutChan() will be called on each reporter, and events placed in this channel will be echoed to the local router.

*/
type Reporter interface {
	Write(eventinfrastructure.Event)
	SetOutChan(chan<- eventinfrastructure.Event)
}
