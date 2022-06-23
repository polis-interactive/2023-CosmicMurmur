package controller

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"log"
)

type service struct {
	repo Repository
	bus  Bus
	cfg  Config

	localAddress    string
	nodeDefinitions types.NodeDefinitions

	controller *controller
}

var _ domain.ControllerService = (*service)(nil)

func NewService(cfg Config, repo Repository, bus Bus) *service {
	log.Println("Lighting, NewService: creating")
	s := &service{
		repo:       repo,
		bus:        bus,
		cfg:        cfg,
		controller: nil,
	}
	s.SetupControllerService()
	return s
}

func (s *service) SetupControllerService() {
	s.initializeVariables()
	s.doCreateController()
}

func (s *service) initializeVariables() {
	var ok bool
	var localAddress string
	localAddress, ok = s.repo.GetControllerLocalAddress()
	if !ok {
		log.Println("Controller, initializeVariables: no local address found, using default")
		localAddress = s.cfg.GetControllerLocalAddress()
	}
	s.localAddress = localAddress
	var nodeDefinitions types.NodeDefinitions
	nodeDefinitions, ok = s.repo.GetControllerNodeDefinitions()
	if !ok {
		log.Println("Controller, initializeVariables: no node definitions found, using default")
		nodeDefinitions = s.cfg.GetControllerNodeDefinitions()
	}
	s.nodeDefinitions = nodeDefinitions
}

func (s *service) doCreateController() {
	s.controller = newController(s.localAddress, s.nodeDefinitions)
}

func (s *service) Startup() {
	for _, n := range s.controller.nodes {
		n.startup()
	}
}

func (s *service) Shutdown() {
	for _, n := range s.controller.nodes {
		n.shutdown()
	}
}

func (s *service) BlackoutNodes() {
	for _, n := range s.controller.nodes {
		func() {
			n.mu.RLock()
			defer n.mu.RUnlock()
			if n.sendChan == nil {
				return
			}
			for _, u := range n.universeNumbers {
				data := n.universePackets[u].Data
				// loop is optimized in assembly by go
				for i := range data {
					data[i] = 0
				}
				n.sendChan <- u
			}
		}()
	}
}

func (s *service) GetUniverseBuffer(universe int) (*[512]byte, bool) {
	if universeBuffer, ok := s.controller.universeBufferMap[universe]; ok {
		return universeBuffer, true
	}
	return nil, false
}

func (s *service) SendUniverseUpdate(universe int) {
	n := s.controller.universeNodeMap[universe]
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.sendChan != nil {
		n.sendChan <- universe
	}
}

func (s *service) GetSettings() *domain.ControllerSettings {
	return &domain.ControllerSettings{
		NodeDefinitions: s.nodeDefinitions,
		LocalAddress:    s.localAddress,
	}
}

func (s *service) SetSettings(settings *domain.ControllerSettings) error {
	err := s.repo.SetControllerNodeDefinitions(settings.NodeDefinitions)
	if err != nil {
		return err
	}
	err = s.repo.SetControllerLocalAddress(settings.LocalAddress)
	if err != nil {
		return err
	}
	s.localAddress = settings.LocalAddress
	s.nodeDefinitions = settings.NodeDefinitions
	s.Shutdown()
	s.doCreateController()
	s.Startup()
	return nil
}
