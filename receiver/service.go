package transmitter

import (
	"context"
	"fmt"

	"github.com/disresc/lib/models"
	"github.com/disresc/lib/util"
	micro "github.com/micro/go-micro"
	transport "github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/util/log"
)

// NewService instantiates a new communication service
func NewService(name string) *Service {
	service := Service{
		name: name,
	}
	return &service
}

// Service implements the kvmtop communication service to transmit events to the
// message bus
type Service struct {
	name             string
	microService     micro.Service
	cancel           context.CancelFunc
	eventChannel     chan *models.Event
	requestPublisher micro.Publisher
}

// Start establishs the communication service
func (service *Service) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	service.cancel = cancel

	service.microService = micro.NewService(
		micro.Name(service.name),
		micro.Context(ctx),
		// for debugging: disable TLS
		micro.Transport(transport.NewTransport(transport.Secure(false))),
		//micro.RegisterInterval(time.Duration(5)*time.Second),
	)

	service.eventChannel = make(chan *models.Event, 10) // buffer max. 10 events
	service.requestPublisher = micro.NewPublisher("requests", service.microService.Client())

	go func() {
		err := service.microService.Run()
		if err != nil {
			log.Errorf("Error while running micro service: %s", err)
		}
	}()
}

// Request sends request and subscribes to matching events
func (service *Service) Request(request *models.Request) {
	topic := util.GetTopicFromRequest(request)
	log.Infof("DisResc receiver request to topic %s", topic)

	// send request
	if err := service.requestPublisher.Publish(context.TODO(), request); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
	micro.RegisterSubscriber(topic, service.microService.Server(), service.handle)
}

// handle takes events and forwards via eventChannel
func (service *Service) handle(ctx context.Context, event *models.Event) error {
	service.eventChannel <- event
	return nil
}

// EventChannel returns a channel providing incoming events
func (service *Service) EventChannel() <-chan *models.Event {
	return service.eventChannel
}

// Close terminates the transmitter service
func (service *Service) Close() {
	service.cancel()
}
