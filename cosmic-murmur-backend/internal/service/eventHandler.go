package service

import (
	"errors"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type eventHandler struct {
	b         *bus
	mu        *sync.Mutex
	shutdowns chan struct{}
	wg        *sync.WaitGroup

	eventQueueSize   int
	eventBusyTimeout time.Duration

	eventQueue     chan *event
	eventQueueLock *sync.RWMutex
}

func newEventHandler(b *bus, conf Config) *eventHandler {

	log.Info().
		Str("package", "service").Str("method", "newEventHandler").
		Msg("creating")

	e := &eventHandler{
		b:         b,
		mu:        &sync.Mutex{},
		shutdowns: nil,
		wg:        &sync.WaitGroup{},

		eventQueueSize:   conf.GetServiceBusEventQueueSize(),
		eventBusyTimeout: conf.GetServiceBusBusyTimeout(),

		eventQueueLock: &sync.RWMutex{},
	}
	e.eventQueue = make(chan *event, e.eventQueueSize)

	log.Info().
		Str("package", "service").Str("method", "newEventHandler").
		Msg("created")

	return e
}

func (e *eventHandler) startup() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	log.Info().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "startup").Msg("starting")

	if e.shutdowns != nil {
		log.Error().
			Str("package", "service").Str("struct", "eventHandler").
			Str("method", "startup").Msg("event loop already running")

		return errors.New("already running")
	}

	e.shutdowns = make(chan struct{})
	e.wg.Add(1)
	go e.runEventLoop()

	log.Info().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "startup").Msg("started")

	return nil
}

func (e *eventHandler) shutdown() {
	e.mu.Lock()
	defer e.mu.Unlock()

	log.Info().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "shutdown").Msg("shutting down")

	// event handler is already stopped
	if e.shutdowns == nil {
		log.Warn().
			Str("package", "service").Str("struct", "eventHandler").
			Str("method", "shutdown").Msg("event loop is already down")
		return
	}

	close(e.shutdowns)
	e.wg.Wait()
	e.shutdowns = nil

	log.Info().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "shutdown").Msg("stopped")

}

func (e *eventHandler) runEventLoop() {

	defer func() {
		log.Info().
			Str("package", "service").Str("struct", "eventHandler").
			Str("method", "runEventLoop").Msg("stopping")
		e.wg.Done()
	}()

	log.Info().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "runEventLoop").Msg("running")

	for {
		select {
		case _, ok := <-e.shutdowns:
			if !ok {
				return
			}
		case eventInstance := <-e.eventQueue:
			e.handleEvent(eventInstance)
		}
	}
}

func (e *eventHandler) handleEvent(eventInstance *event) {

	log.Debug().
		Str("package", "service").Str("struct", "eventHandler").
		Str("method", "handleEvent").Uint64("trace", eventInstance.TraceId).
		Msg("running event")

	/*
		since we are using commands to create each of the events, we don't need
		to type check on the event payload; it does mean we have to manually cast
		the payload to its specific type though
	*/
	switch eventInstance.Name {
	// graphics bus
	case RequestGridDimensions:
		e.RetrieveGridDimensions(eventInstance, eventInstance.Payload.(chan *types.Grid))
	case GraphicsReady:
		e.UpdateRenderFromGraphics(eventInstance)
	case GraphicsCrashed:
		e.ClearGraphics(eventInstance)
	// api calls
	case FetchGraphicsSettings:
		e.FetchGraphicsSettings(eventInstance, eventInstance.Payload.(chan *domain.GraphicsSettings))
	case SetGraphicsSettings:
		e.SetGraphicsSettings(eventInstance, eventInstance.Payload.(*setGraphicsSettingsPayload))
	case FetchLightingSettings:
		e.FetchLightingSettings(eventInstance, eventInstance.Payload.(chan *domain.LightingSettings))
	case SetLightingSettings:
		e.SetLightingSettings(eventInstance, eventInstance.Payload.(*setLightingSettingsPayload))
	}

	if l := log.Debug(); l.Enabled() {
		l.Str("package", "service").Str("struct", "eventHandler").
			Str("method", "handleEvent").Uint64("trace", eventInstance.TraceId).
			Msgf("Handled event in %v", time.Since(eventInstance.TimeRequested))
	}
}

func (e *eventHandler) pumpQueue() error {
	for {
		select {
		case eventInstance, ok := <-e.eventQueue:
			if ok {
				e.closeoutEvent(eventInstance)
			} else {
				return errors.New("queue unexpectedly closed")
			}
		default:
			return nil
		}
	}
}

func (e *eventHandler) closeoutEvent(eventInstance *event) {

	switch eventInstance.Name {
	// graphics bus
	case RequestGridDimensions:
		close(eventInstance.Payload.(chan *types.Grid))
	case GraphicsReady:
		return
	case GraphicsCrashed:
		return
	// api calls
	case FetchGraphicsSettings:
		close(eventInstance.Payload.(chan *domain.GraphicsSettings))
	case SetGraphicsSettings:
		close(eventInstance.Payload.(*setGraphicsSettingsPayload).DispatchChannel)
	case FetchLightingSettings:
		close(eventInstance.Payload.(chan *domain.LightingSettings))
	case SetLightingSettings:
		close(eventInstance.Payload.(*setLightingSettingsPayload).DispatchChannel)
	}
}
