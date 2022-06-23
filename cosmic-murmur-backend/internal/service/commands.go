package service

import (
	"errors"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"log"
	"time"
)

/*
	GraphicsService bus commands
*/

func (b *bus) GetGridDimensions() *types.Grid {
	responseChannel := make(chan *types.Grid)
	defaultResponse := &types.Grid{
		MinX: -1,
		MaxX: 1,
		MinY: -1,
		MaxY: 1,
	}
	err := tryEnqueueEvent(b, RequestGridDimensions, responseChannel)
	if err != nil {
		return defaultResponse
	}
	resp, err := waitForResponse[*types.Grid](b, responseChannel)
	if err != nil {
		return defaultResponse
	} else {
		return resp
	}
}

func (b *bus) EmitGraphicsReady() {
	err := tryEnqueueEvent(b, GraphicsReady, nil)
	if err != nil {
		log.Printf("coulnd't enqueue event")
	}
}

func (b *bus) EmitGraphicsCrashed() {
	err := tryEnqueueEvent(b, GraphicsCrashed, nil)
	if err != nil {
		log.Printf("coulnd't enqueue event")
	}
}

/*
	API bus commands
*/

func (b *bus) FetchGraphicsSettings() (*domain.GraphicsSettings, error) {
	responseChannel := make(chan *domain.GraphicsSettings)
	err := tryEnqueueEvent(b, FetchGraphicsSettings, responseChannel)
	if err != nil {
		return nil, err
	}
	resp, err := waitForResponse[*domain.GraphicsSettings](b, responseChannel)
	if err != nil || resp == nil {
		return nil, err
	}
	return resp, nil
}

func (b *bus) SetGraphicsSettings(shaderName string, refreshInMs int64, reloadOnUpdate bool) error {
	responseChannel := make(chan struct{})
	err := tryEnqueueEvent(b, SetGraphicsSettings, &setGraphicsSettingsPayload{
		DispatchChannel: responseChannel, GraphicsFrequency: time.Duration(refreshInMs) * time.Millisecond,
		ShaderName: shaderName, ReloadOnUpdate: reloadOnUpdate,
	})
	if err != nil {
		return err
	}
	_, err = waitForResponse[struct{}](b, responseChannel)
	return err
}

func (b *bus) FetchLightingSettings() (*domain.LightingSettings, error) {
	responseChannel := make(chan *domain.LightingSettings)
	err := tryEnqueueEvent(b, FetchLightingSettings, responseChannel)
	if err != nil {
		return nil, err
	}
	resp, err := waitForResponse[*domain.LightingSettings](b, responseChannel)
	if err != nil || resp == nil {
		return nil, err
	}
	return resp, nil
}

func (b *bus) SetLightingSettings(segmentDefinition types.LedSegment, segmentCount int) error {
	responseChannel := make(chan struct{})
	err := tryEnqueueEvent(b, SetLightingSettings, &setLightingSettingsPayload{
		DispatchChannel: responseChannel, SegmentCount: segmentCount,
		SegmentDefinition: segmentDefinition,
	})
	if err != nil {
		return err
	}
	_, err = waitForResponse[struct{}](b, responseChannel)
	return err
}

/*
	Common abstractions
*/

func tryEnqueueEvent(b *bus, e eventType, payload interface{}) error {
	eh := b.eventHandler
	eh.eventQueueLock.RLock()
	defer eh.eventQueueLock.RUnlock()
	if eh.eventQueue != nil {
		eh.eventQueue <- &event{
			Name:          e,
			TimeRequested: time.Now(),
			TraceId:       b.GetEventTraceId(),
			Payload:       payload,
		}
		return nil
	} else {
		return errors.New("event queue closed")
	}
}

func waitForResponse[T any](b *bus, responseChan chan T) (t T, err error) {
	eh := b.eventHandler
	select {
	case _, ok := <-eh.shutdowns:
		if !ok {
			return t, errors.New("eventloop shutdown")
		} else {
			return t, errors.New("illegal state")
		}
	case resp, ok := <-responseChan:
		if ok {
			return resp, nil
		} else {
			return t, errors.New("eventloop closed connection")
		}
	case <-time.After(eh.eventBusyTimeout):
		return t, errors.New("eventloop not responding")
	}
}
