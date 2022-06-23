package service

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"github.com/rs/zerolog/log"
	"sync"
)

func (e *eventHandler) RetrieveGridDimensions(eventInstance *event, dispatchChannel chan *types.Grid) {

	log.Trace().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "RetrieveGridDimensions").Uint64("trace", eventInstance.TraceId).
		Msg("getting grid dimensions")

	grid := e.b.lightingService.GetGridDimensions()
	dispatchChannel <- grid
	// dispatch channel should be garbage collected after command returns grid to graphics service
}

func (e *eventHandler) UpdateRenderFromGraphics(eventInstance *event) {

	log.Trace().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "UpdateRenderFromGraphics").Uint64("trace", eventInstance.TraceId).
		Msg("updating renderer")

	pb, gMuPreRLocked := e.b.graphicsService.GetPb()
	lightUniverses := e.b.lightingService.GetLightUniverses()

	wg := &sync.WaitGroup{}
	for universe, lights := range lightUniverses {
		universeBuffer, ok := e.b.controllerService.GetUniverseBuffer(universe)
		if !ok {
			continue
		}
		wg.Add(1)
		go func(universe int, lights []*types.Light) {
			defer func() {
				wg.Done()
			}()
			for _, l := range lights {
				c := pb.GetPixel(&l.Position)
				universeBuffer[l.Pixel*3] = c.R
				universeBuffer[l.Pixel*3+1] = c.G
				universeBuffer[l.Pixel*3+2] = c.B
			}
			e.b.controllerService.SendUniverseUpdate(universe)
		}(universe, lights)
	}
	wg.Wait()
	gMuPreRLocked.RUnlock()
}

func (e *eventHandler) ClearGraphics(eventInstance *event) {
	log.Trace().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "ClearGraphics").Uint64("trace", eventInstance.TraceId).
		Msg("clearing graphics")

	e.b.controllerService.BlackoutNodes()
}

func (e *eventHandler) FetchGraphicsSettings(eventInstance *event, dispatchChannel chan *domain.GraphicsSettings) {
	log.Trace().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "FetchGraphicsSettings").Uint64("trace", eventInstance.TraceId).
		Msg("getting graphics settings")

	settings, err := e.b.graphicsService.GetSettings()
	if err != nil {
		log.Warn().
			Str("package", "service").Str("struct", "eventHandler").
			Str("method", "FetchGraphicsSettings").Uint64("trace", eventInstance.TraceId).
			Err(err).Msg("error retrieving graphics settings")
		close(dispatchChannel)
		return
	}

	dispatchChannel <- settings
	// dispatch channel should be garbage collected after command returns settings to api
}

func (e *eventHandler) SetGraphicsSettings(eventInstance *event, payload *setGraphicsSettingsPayload) {
	log.Trace().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "SetGraphicsSettings").Uint64("trace", eventInstance.TraceId).
		Msg("setting graphics settings")

	err := e.b.graphicsService.SetSettings(&domain.GraphicsSettableSettings{
		ShaderName:     payload.ShaderName,
		Frequency:      payload.GraphicsFrequency,
		ReloadOnUpdate: payload.ReloadOnUpdate,
	})
	if err != nil {
		log.Warn().
			Str("package", "service").Str("struct", "eventHandler").
			Str("method", "SetGraphicsSettings").Uint64("trace", eventInstance.TraceId).
			Err(err).Msg("error setting graphics settings")
		close(payload.DispatchChannel)
		return
	}

	payload.DispatchChannel <- struct{}{}
	// dispatch channel should be garbage collected after command returns success to api
}

func (e *eventHandler) FetchLightingSettings(eventInstance *event, dispatchChannel chan *domain.LightingSettings) {
	log.Trace().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "FetchLightingSettings").Uint64("trace", eventInstance.TraceId).
		Msg("fetching lighting settings")
	settings := e.b.lightingService.GetSettings()
	dispatchChannel <- settings
	// dispatch channel should be garbage collected after command returns settings to api
}

func (e *eventHandler) SetLightingSettings(eventInstance *event, payload *setLightingSettingsPayload) {
	log.Trace().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "SetLightingSettings").Uint64("trace", eventInstance.TraceId).
		Msg("setting lighting settings")

	err := e.b.lightingService.SetSettings(&domain.LightingSettings{
		SegmentDefinition: payload.SegmentDefinition,
		SegmentCount:      payload.SegmentCount,
	})
	if err != nil {
		log.Warn().
			Str("package", "service").Str("struct", "eventHandler").
			Str("method", "SetLightingSettings").Uint64("trace", eventInstance.TraceId).
			Err(err).Msg("error setting lighting settings")
		close(payload.DispatchChannel)
		return
	}

	payload.DispatchChannel <- struct{}{}
	// dispatch channel should be garbage collected after command returns success to api
}

func (e *eventHandler) ResetApplication(eventInstance *event, dispatchChan chan struct{}) {
	log.Trace().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "ResetApplication").Uint64("trace", eventInstance.TraceId).
		Msg("resetting services; this may take up to a second")

	err := e.b.repo.ResetRepository()
	if err != nil {
		log.Warn().
			Str("package", "service").Str("struct", "eventHandler").
			Str("method", "ResetApplication").Uint64("trace", eventInstance.TraceId).
			Err(err).Msg("couldn't reset repository")
		close(dispatchChan)
		return
	}
	/*
		we could technically pump the queue here, but more things could just get queued in the meantime; instead,
		we are going to rely on busyTimeout to force listeners of events to timeout before shutting them down and
		resetting things
	*/
	// make sure services can't still write to event queue
	e.eventQueueLock.Lock()
	e.eventQueue = nil
	e.eventQueueLock.Unlock()
	// shutdown services; rely on the busyTimeout for them to return on response channels
	e.b.graphicsService.Shutdown()
	e.b.controllerService.Shutdown()
	// garbage collect old events
	e.eventQueue = make(chan *event, e.eventQueueSize)
	// reconfig
	e.b.lightingService.SetupLightingService()
	e.b.controllerService.SetupControllerService()
	// start program loops back up
	e.b.controllerService.Startup()
	e.b.graphicsService.Startup()
}
