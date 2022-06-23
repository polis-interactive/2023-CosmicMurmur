package domain

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"sync"
	"time"
)

type GraphicsSettableSettings struct {
	ShaderName     string
	Frequency      time.Duration
	ReloadOnUpdate bool
}

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
	SetSettings(settings *GraphicsSettableSettings) error
	GetPb() (pb *types.PixelBuffer, preLockedMutex *sync.RWMutex)
}

type LightingSettings struct {
	SegmentDefinition types.LedSegment
	SegmentCount      int
}

type LightingService interface {
	SetupLightingService()
	GetSettings() *LightingSettings
	SetSettings(settings *LightingSettings) error
	GetGridDimensions() *types.Grid
	GetLightUniverses() [][]*types.Light
}

type ControllerSettings struct {
	NodeDefinitions types.NodeDefinitions
	LocalAddress    string
}

type ControllerService interface {
	Startup()
	Shutdown()
	SetupControllerService()
	BlackoutNodes()
	GetUniverseBuffer(universe int) (universeBuffer *[512]byte, ok bool)
	SendUniverseUpdate(universe int)
	GetSettings() *ControllerSettings
	SetSettings(settings *ControllerSettings) error
}
