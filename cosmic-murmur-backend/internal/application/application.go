package application

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain/controller"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain/graphics"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain/lighting"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/infrastructure/repository/memory"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/service"
	"log"
	"sync"
)

type Application struct {
	memoryRepository *memory.Repository
	serviceBus       applicationBus
	shutdown         bool
	shutdownLock     *sync.Mutex
}

func NewApplication(conf *Config) (*Application, error) {
	log.Println("Application, NewApplication: creating")

	/* create application instance */
	app := &Application{
		shutdown:     true,
		shutdownLock: &sync.Mutex{},
	}

	/* create repositories */
	memoryRepository := memory.NewMemoryRepository()
	app.memoryRepository = memoryRepository

	/* create bus */
	app.serviceBus = service.NewBus(conf, app.memoryRepository)

	/* create services */
	lightingService := lighting.NewService(conf, app.memoryRepository)
	app.serviceBus.BindLightingService(lightingService)

	graphicsService, err := graphics.NewService(conf, app.memoryRepository, app.serviceBus)
	if err != nil {
		log.Fatalln("Application, NewApplication: failed to initialize graphics service")
	}
	app.serviceBus.BindGraphicsService(graphicsService)

	controllerService := controller.NewService(conf, app.memoryRepository, app.serviceBus)
	app.serviceBus.BindControllerService(controllerService)

	return app, nil
}

func (app *Application) Startup() error {

	log.Println("Application, Startup: starting")

	app.shutdownLock.Lock()
	defer app.shutdownLock.Unlock()
	if app.shutdown == false {
		return nil
	}

	app.shutdown = false

	err := app.serviceBus.Startup()
	if err != nil {
		return err
	}

	log.Println("Application, Startup: started")

	return nil
}

func (app *Application) Shutdown() error {

	log.Println("Application, Shutdown: shutting down")

	app.shutdownLock.Lock()
	defer app.shutdownLock.Unlock()
	if app.shutdown {
		return nil
	}
	app.shutdown = true

	app.serviceBus.Shutdown()

	log.Println("Application, Shutdown: finished")

	return nil
}
