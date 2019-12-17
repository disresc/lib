package transmitter

import (
	"context"
	fmt "fmt"
	"sync"

	models "github.com/disresc/lib/models"
	"github.com/disresc/lib/util"
	micro "github.com/micro/go-micro"
	transport "github.com/micro/go-micro/transport"
	"github.com/prometheus/common/log"
)

// NewService instantiates a new communication service
func NewService(name string) *Service {
	service := Service{
		name:                  name,
		requesthandlers:       make(map[string]*requesthandler),
		requesthandlersAccess: sync.Mutex{},
	}
	return &service
}

// Service implements the kvmtop communication service to transmit events to the
// message bus
type Service struct {
	name                  string
	microService          micro.Service
	cancel                context.CancelFunc
	requesthandlers       map[string]*requesthandler
	requesthandlersAccess sync.Mutex
	requestCallback       func(request *models.Request) bool
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

	// wait for requests to open publishers
	micro.RegisterSubscriber("requests", service.microService.Server(), service.handleRequests)

	go func() {
		err := service.microService.Run()
		if err != nil {
			log.Errorf("Error while running micro service: %s", err)
		}
	}()
}

// RegisterRequestCallback registers the callback used for handling requests
func (service *Service) RegisterRequestCallback(callback func(request *models.Request) bool) {
	service.requestCallback = callback
}

func (service *Service) handleRequests(ctx context.Context, request *models.Request) error {
	log.Infof("DisResc transmitter receives request %s", request.String())
	// requesthandler exists already?
	requesthandler, found := service.findRequesthandler(request)
	if !found {
		// create new request handler
		if service.requestCallback == nil {
			return nil
		}
		acceptRequest := service.requestCallback(request)
		if acceptRequest {
			log.Infof("DisResc transmitter accepts request %s", request.String())
			service.createRequesthandler(request)
		}
	} else {
		// update existing request handler
		requesthandler.Update(request)
	}
	return nil
}

func (service *Service) findRequesthandler(request *models.Request) (*requesthandler, bool) {
	topic := util.GetTopicFromRequest(request)
	service.requesthandlersAccess.Lock()
	defer service.requesthandlersAccess.Unlock()
	rh, found := service.requesthandlers[topic]
	return rh, found
}

func (service *Service) createRequesthandler(request *models.Request) {
	topic := util.GetTopicFromRequest(request)
	r := newRequesthandler(service.microService, request)
	service.requesthandlersAccess.Lock()
	defer service.requesthandlersAccess.Unlock()
	service.requesthandlers[topic] = r
}

// Publish allows to transmit an event through the communication service
func (service *Service) Publish(event *models.Event, interval int) error {
	topics := util.GetTopicsFromEvent(event, interval)
	service.requesthandlersAccess.Lock()
	defer service.requesthandlersAccess.Unlock()
	for _, topic := range topics {
		rh, exists := service.requesthandlers[topic]
		if !exists {
			// nobody requested event on this topic
			continue
		}
		log.Infof("DisResc transmitter publishes event to topic %s", topic)
		if err := rh.GetPublisher().Publish(context.TODO(), event); err != nil {
			return fmt.Errorf("error publishing: %v", err)
		}
	}
	return nil
}

// Close terminates the transmitter service
func (service *Service) Close() {
	service.cancel()
}
