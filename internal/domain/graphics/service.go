package graphics

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"log"
	"sync"
)

type service struct {
	cfg  Config
	repo Repository
	bus  Bus

	mu        *sync.RWMutex
	wg        *sync.WaitGroup
	shutdowns chan struct{}

	subscriptions []int

	g *Graphics
}

var _ domain.GraphicsService = (*service)(nil)

func NewService(cfg Config, repo Repository, bus Bus) (*service, error) {
	log.Println("Graphics, NewService: creating")

	s := &service{
		cfg:           cfg,
		repo:          repo,
		bus:           bus,
		mu:            &sync.RWMutex{},
		wg:            &sync.WaitGroup{},
		subscriptions: nil,
	}
	g, err := newGraphics(s)
	if err != nil {
		log.Println("Graphics, NewService: error creating graphics")
		return nil, err
	}
	s.g = g
	s.subscribeToEvents()
	return s, nil
}

func (s *service) subscribeToEvents() {
	s.subscriptions = make([]int, 4)
	s.subscriptions[0] = s.bus.SubscribeToEvent(domain.LightingUpdated, s.Reset)
	s.subscriptions[1] = s.bus.SubscribeToEvent(domain.StartedDMXConsole, s.Shutdown)
	s.subscriptions[2] = s.bus.SubscribeToEvent(domain.StoppedDmxConsole, s.Startup)
	s.subscriptions[3] = s.bus.SubscribeToEvent(domain.ClientDisconnected, s.Startup)
}

func (s *service) doStartup() {
	if s.shutdowns == nil {
		s.shutdowns = make(chan struct{})
		s.wg.Add(1)
		go s.g.runMainLoop()
	}
}

func (s *service) Startup() {
	log.Println("GraphicsService Startup: starting")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.doStartup()
}

func (s *service) doShutdown() {
	if s.shutdowns != nil {
		close(s.shutdowns)
		s.wg.Wait()
		s.shutdowns = nil
	}
}

func (s *service) Shutdown() {
	log.Println("RenderService Shutdown: shutting down")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.doShutdown()
}

func (s *service) Reset() {
	log.Println("RenderService Startup: resetting")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.doShutdown()
	s.doStartup()
}

func (s *service) GetPb() (pb *types.PixelBuffer, preLockedMutex *sync.RWMutex) {
	//TODO implement me
	panic("implement me")
}
