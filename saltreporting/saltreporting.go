package saltreporting

import (
	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-translator-microservice/common"
)

type saltReporter struct{}

func (s *saltReporter) Write(event eventinfrastructure.Event) {
}

func (s *saltReporter) SetOutChan(chan<- eventinfrastructure.Event) {
}

func GetReporter() common.Reporter {
	return &saltReporter{}
}
