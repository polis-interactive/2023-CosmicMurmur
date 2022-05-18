package application

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain/graphics"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/infrastructure/bus"
	"log"
	"sync"
)

type Application struct {
	serviceBus   applicationBus
	shutdown     bool
	shutdownLock *sync.Mutex
}

func NewApplication(conf *Config) (*Application, error) {
	log.Println("Application, NewApplication: creating")

	/* create application instance */
	app := &Application{
		shutdown: true,
	}

	/* create bus */
	app.serviceBus = bus.NewBus()

	/* create services */
	graphicsService, err := graphics.NewService()

	return app, nil
}
