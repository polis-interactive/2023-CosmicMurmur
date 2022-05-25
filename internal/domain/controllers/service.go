package controllers

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"log"
	"sync"
)

type service struct {
	repo Repository
	bus  Bus

	mu        *sync.RWMutex
	wg        *sync.WaitGroup
	shutdowns chan struct{}

	controllers  []*controller
	localAddress string
}

var _ domain.ControllerService = (*service)(nil)

func (s *service) doStartup() {
	if s.shutdowns == nil {
		s.shutdowns = make(chan struct{})
		// s.wg.Add(1)
		// go s.runMainLoop()
	}
}

func (s *service) Startup() {
	log.Println("ControllerService Startup: starting")
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
	log.Println("ControllerService Shutdown: shutting down")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.doShutdown()
}

func (s *service) doReset() {
	s.doShutdown()
	s.doStartup()
}

func (s *service) doCreateControllers(definitions types.ControllerDefinitions) {
	s.controllers = make([]*controller, len(definitions))
	for i, d := range definitions {
		s.controllers[i] = newController(d)
	}
}

func (s *service) GetSettings() *domain.ControllerSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &domain.ControllerSettings{
		ControllerDefinitions: mapControllersToDefinitions(s.controllers),
		LocalAddress:          s.localAddress,
	}
}

func (s *service) SetControllerDefinitions(definitions types.ControllerDefinitions) error {
	err := s.repo.SetControllerDefinitions(definitions)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bus.StopRenderer()
	s.doShutdown()
	s.doCreateControllers(definitions)
	s.doStartup()
	s.bus.StartRenderer()
	return nil
}

func (s *service) SetLocalAddress(addr string) error {
	err := s.repo.SetControllerLocalAddress(addr)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.localAddress = addr
	s.doReset()
	return nil
}
