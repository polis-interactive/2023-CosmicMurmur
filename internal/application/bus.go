package application

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain/graphics"
)

type applicationBus interface {
	Startup()
	Shutdown()
	BindGraphicsService(graphicsClient domain.GraphicsService)
	graphics.Bus
}
