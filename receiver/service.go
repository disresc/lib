package transmitter

import (
	"context"
	"fmt"

	drModels "github.com/disresc/lib/models"
	micro "github.com/micro/go-micro"
	transport "github.com/micro/go-micro/transport"
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
	name         string
	microService micro.Service
	cancel       context.CancelFunc
	eventChannel chan *drModels.Event
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
	)

	service.microService.Init()

	// register subscriber
	service.eventChannel = make(chan *drModels.Event, 10)
	micro.RegisterSubscriber("monitoring", service.microService.Server(), service.handle)

	service.microService.Run()
}

func (service *Service) startService() {
	if err := service.microService.Run(); err != nil {
		fmt.Printf("error while starting service: %v", err)
	}
}

func (service *Service) handle(ctx context.Context, event *drModels.Event) error {
	fmt.Printf("handle event %v from context %v\n", event, ctx)
	service.eventChannel <- event
	return nil
}

// EventChannel returns a channel providing incoming events
func (service *Service) EventChannel() <-chan *drModels.Event {
	return service.eventChannel
}

// Close terminates the transmitter service
func (service *Service) Close() {
	service.cancel()
}
