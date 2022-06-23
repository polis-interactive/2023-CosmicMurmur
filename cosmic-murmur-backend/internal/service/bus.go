package service

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"sync/atomic"
)

type bus struct {
	graphicsService   domain.GraphicsService
	lightingService   domain.LightingService
	controllerService domain.ControllerService
	repo              Repository
	eventHandler      *eventHandler
	nextEventTraceId  uint64
}

func NewBus(conf Config, repo Repository) *bus {
	b := &bus{
		nextEventTraceId: 0,
		repo:             repo,
	}
	b.eventHandler = newEventHandler(b, conf)
	return b
}

func (b *bus) BindGraphicsService(graphicsClient domain.GraphicsService) {
	b.graphicsService = graphicsClient
}

func (b *bus) BindLightingService(lightingService domain.LightingService) {
	b.lightingService = lightingService
}

func (b *bus) BindControllerService(controllerClient domain.ControllerService) {
	b.controllerService = controllerClient
}

func (b *bus) Startup() error {
	err := b.eventHandler.startup()
	if err != nil {
		return err
	}
	b.controllerService.Startup()
	b.graphicsService.Startup()
	return nil
}

func (b *bus) Shutdown() {
	b.eventHandler.shutdown()
	b.graphicsService.Shutdown()
	b.controllerService.Shutdown()
}

func (b *bus) GetEventTraceId() uint64 {
	return atomic.AddUint64(&b.nextEventTraceId, 1)
}
