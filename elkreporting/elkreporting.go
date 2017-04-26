package elkreporting

import (
	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/byuoitav/event-translator-microservice/common"
)

type elkReporter struct{}

func (e *elkReporter) Write(event eventinfrastructure.Event) {
}

func GetReporter() *common.Reporter {
	//Do whatever initialization is necessary

	return elkReporter{}
}
