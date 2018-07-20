package awsshadowreporting

import (
	"github.com/byuoitav/common/events"
	"github.com/byuoitav/event-translator-microservice/common"
)

type awsShadowReporter struct {
}

//Write fulfils the Reporter interface requirment.
//This function will be called each time an event arrives from the router. All operations taken as a response to this function MUST be threadsafe.
func (aws *awsShadowReporter) Write(event events.Event) {

}

//SetOutChan() will be called on each reporter, and events placed in the channel will be echoed to the local router.
func (s *awsShadowReporter) SetOutChan(chan<- events.Event) {
}

func GetReporter() common.Reporter {
	//Do whatever initialization is necessary here

	return &awsShadowReporter{}
}
