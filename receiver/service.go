package transmitter

import (
	"context"
	"fmt"
	"time"

	"github.com/disresc/lib/models"
	"github.com/disresc/lib/util"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/transport"
	"github.com/micro/go-micro/util/log"
)

// NewService instantiates a new communication service
func NewService(name string) *Service {
	service := Service{
		name:    name,
		running: false,
	}

	ctx, cancel := context.WithCancel(context.Background())
	service.cancel = cancel
	service.microService = micro.NewService(
		micro.Name(service.name),
		micro.Context(ctx),
		// for debugging: disable TLS
		micro.Transport(transport.NewTransport(transport.Secure(false))),
		//micro.RegisterInterval(time.Duration(5)*time.Second),
		micro.AfterStart(service.afterStart),
	)

	service.eventChannel = make(chan *models.Event, 10) // buffer max. 10 events
	service.requestPublisher = micro.NewPublisher("requests", service.microService.Client())

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
	running          bool
	requests         []*models.Request
}

// Start establishs the communication service
func (service *Service) Start() {
	go func() {
		err := service.microService.Run()
		if err != nil {
			log.Errorf("Error while running micro service: %s", err)
		}
	}()

	service.running = true
}

// RegisterData specifies the type of data to register for. Has to be called
// BEFORE Start()
func (service *Service) RegisterData(source string, transmitter string, interval int) error {
	if service.running {
		return fmt.Errorf("failed to register data: service already running")
	}

	request := models.Request{
		Timeout:     15,
		Source:      source,
		Transmitter: transmitter,
		Interval:    uint32(interval),
	}
	service.registerRequest(&request)
	service.requests = append(service.requests, &request)

	return nil
}

// registerRequest adds a subscriber to the request's topic
func (service *Service) registerRequest(request *models.Request) {
	topic := util.GetTopicFromRequest(request)
	micro.RegisterSubscriber(topic, service.microService.Server(), service.handle)
}

// sendPeriodicRequest starts a periodic ticker to sendRequest
func (service *Service) sendPeriodicRequest(request *models.Request) {
	ticker := time.NewTicker(time.Duration(request.Timeout) * time.Second)
	service.sendRequest(request)
	go func() {
		for {
			select {
			//case <-done:
			//	return
			case <-ticker.C:
				service.sendRequest(request)
			}
		}
	}()
}

// sendRequest publishes the request to the request topic
func (service *Service) sendRequest(request *models.Request) {
	if err := service.requestPublisher.Publish(context.TODO(), request); err != nil {
		fmt.Printf("error publishing: %v", err)
	}
}

// handle takes events and forwards via eventChannel
func (service *Service) handle(ctx context.Context, event *models.Event) error {
	service.eventChannel <- event
	return nil
}

func (service *Service) afterStart() error {
	// after service started, start periodic request sending
	log.Infof("Start sending %d requests periodically", len(service.requests))
	for _, request := range service.requests {
		service.sendPeriodicRequest(request)
	}
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
