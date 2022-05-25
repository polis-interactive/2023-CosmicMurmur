package domain

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"sync"
	"time"
)

type GraphicsSettings struct {
	Shaders        []string
	RunningShader  string
	Frequency      time.Duration
	ReloadOnUpdate bool
}

type GraphicsService interface {
	Startup()
	Reset()
	Shutdown()
	GetSettings() (*GraphicsSettings, error)
	SetShader(shaderName string) error
	SetFrequency(updateFrequency time.Duration) error
	SetReloadOnUpdate(bool) error
	GetPb() (pb *types.PixelBuffer, preLockedMutex *sync.RWMutex)
}

type LightingSettings struct {
	SegmentDefinition types.LedSegment
	SegmentCount      int
}

type LightingService interface {
	GetSettings() *LightingSettings
	SetSegmentDefinition(types.LedSegment) error
	SetSegmentCount(int) error
	GetGridDimensions() *types.Grid
}

type ControllerSettings struct {
	ControllerDefinitions types.ControllerDefinitions
	LocalAddress          string
}

type ControllerService interface {
	Startup()
	Shutdown()
	GetSettings() *ControllerSettings
	SetControllerDefinitions(types.ControllerDefinitions) error
	SetLocalAddress(string) error
}

type EventSubscriber interface {
	SubscribeToEvent(e Event, fn func()) int
	UnsubscribeToEvent(handlerId int)
}

type EventHandler interface {
	HandleEvent(e Event)
}

type EventBus interface {
	EventSubscriber
	EventHandler
}
