package graphics

import (
	"errors"
	"fmt"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"github.com/polis-interactive/go-lighting-utils/pkg/graphicsShader"
	"log"
	"sync"
)

type service struct {
	repo Repository
	bus  Bus
	cfg  Config

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
		repo:          repo,
		bus:           bus,
		cfg:           cfg,
		mu:            &sync.RWMutex{},
		wg:            &sync.WaitGroup{},
		subscriptions: nil,
	}
	g, err := newGraphics(s, cfg)
	if err != nil {
		log.Println("Graphics, NewService: error creating graphics")
		return nil, err
	}
	s.g = g
	return s, nil
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
	log.Println("GraphicsService Shutdown: shutting down")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.doShutdown()
}

func (s *service) Reset() {
	log.Println("GraphicsService Startup: resetting")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.doShutdown()
	s.doStartup()
}

func (s *service) GetPb() (pb *types.PixelBuffer, preLockedMutex *sync.RWMutex) {
	s.g.mu.RLock()
	return s.g.pb, s.g.mu
}

func (s *service) GetSettings() (*domain.GraphicsSettings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.shutdowns == nil {
		return nil, errors.New("GraphicsService is down")
	}
	s.g.mu.RLock()
	defer s.g.mu.RUnlock()
	if s.g.shaderList == nil {
		return nil, errors.New("GraphicsLoop is down")
	}
	var shaders []string
	for _, v := range s.g.shaderList {
		shaders = append(shaders, v)
	}
	return &domain.GraphicsSettings{
		Shaders:        shaders,
		RunningShader:  s.g.runningShader,
		Frequency:      s.g.runningFrequency,
		ReloadOnUpdate: s.g.runningReloadOnUpdate,
	}, nil
}

func (s *service) SetSettings(settings *domain.GraphicsSettableSettings) error {
	shaderKey := graphicsShader.ShaderKey(settings.ShaderName)
	err := func() error {
		s.g.mu.RLock()
		defer s.g.mu.RUnlock()
		if s.g.shaderList == nil {
			return errors.New("GraphicsLoop is Down")
		}
		if _, ok := s.g.shaderList[shaderKey]; !ok {
			return errors.New(fmt.Sprintf("Shader %s not found", settings.ShaderName))
		}
		return nil
	}()
	if err != nil {
		return err
	}
	err = s.repo.SetGraphicsShader(settings.ShaderName)
	if err != nil {
		return err
	}
	err = s.repo.SetGraphicsFrequency(settings.Frequency)
	if err != nil {
		return err
	}
	err = s.repo.SetGraphicsReloadOnUpdate(settings.ReloadOnUpdate)
	if err != nil {
		return err
	}
	s.g.mu.Lock()
	defer s.g.mu.Unlock()
	s.g.runningFrequency = settings.Frequency
	s.g.runningReloadOnUpdate = settings.ReloadOnUpdate
	return s.g.setShader(shaderKey)
}
