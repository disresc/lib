package transmitter

import (
	"context"
	fmt "fmt"

	drModels "github.com/disresc/lib/models"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/transport"
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
	publisher    micro.Publisher
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
	//service.microService.Init()

	service.publisher = micro.NewPublisher("monitoring", service.microService.Client())

}

// Publish allows to transmit an event through the communication service
func (service *Service) Publish(event drModels.Event) error {
	if err := service.publisher.Publish(context.TODO(), event); err != nil {
		return fmt.Errorf("error publishing: %v", err)
	}
	return nil
}

// Close terminates the transmitter service
func (service *Service) Close() {
	service.cancel()
}
