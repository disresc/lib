package transmitter

import (
	"github.com/disresc/lib/models"
	"github.com/disresc/lib/util"
	micro "github.com/micro/go-micro"
	"time"
)

func newRequesthandler(microService micro.Service, request *models.Request) *requesthandler {
	rh := requesthandler{
		microService: microService,
		topic:        util.GetTopicFromRequest(request),
		endoflife:    time.Now().Add(time.Duration(request.Timeout) * time.Second),
	}
	rh.createPublisher()
	return &rh
}

type requesthandler struct {
	microService micro.Service
	topic        string
	publisher    micro.Publisher
	endoflife    time.Time
}

func (requesthandler *requesthandler) createPublisher() {
	requesthandler.publisher = micro.NewPublisher(requesthandler.topic, requesthandler.microService.Client())
}

func (requesthandler *requesthandler) GetPublisher() micro.Publisher {
	return requesthandler.publisher
}

func (requesthandler *requesthandler) Update(request *models.Request) {
	newEndoflife := time.Now().Add(time.Duration(request.Timeout) * time.Second)
	if requesthandler.endoflife.Before(newEndoflife) {
		// update deadline
		requesthandler.endoflife = newEndoflife
	}
}
