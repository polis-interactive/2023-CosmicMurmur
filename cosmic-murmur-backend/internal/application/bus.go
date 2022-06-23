package application

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain/controller"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain/graphics"
)

type applicationBus interface {
	Startup() error
	Shutdown()
	BindGraphicsService(graphicsClient domain.GraphicsService)
	BindLightingService(lightingService domain.LightingService)
	BindControllerService(controllerClient domain.ControllerService)
	graphics.Bus
	controller.Bus
}
