package bus

import "github.com/polis-interactive/2023-CosmicMurmur/internal/domain"

type bus struct {
	*eventHandler
	graphicsService domain.GraphicsService
}

func NewBus() *bus {
	h := newEventHandler()
	return &bus{
		eventHandler: h,
	}
}

func (b *bus) BindGraphicsService(graphicsClient domain.GraphicsService) {
	b.graphicsService = graphicsClient
}

func (b *bus) Startup() {
	b.startupEventLoop()
	b.graphicsService.Startup()
}

func (b *bus) Shutdown() {
	b.shutdownEventLoop()
	b.graphicsService.Shutdown()
}
