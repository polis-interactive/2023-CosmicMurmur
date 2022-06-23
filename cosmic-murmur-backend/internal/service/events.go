package service

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"time"
)

type eventType int64

const (
	RequestGridDimensions eventType = iota
	GraphicsCrashed
	GraphicsReady

	FetchGraphicsSettings
	SetGraphicsSettings
	FetchLightingSettings
	SetLightingSettings
)

func (s eventType) String() string {
	switch s {
	case RequestGridDimensions:
		return "Request Grid Dimensions"
	case GraphicsCrashed:
		return "Graphics Crashed"
	case GraphicsReady:
		return "Graphics Ready"
	case FetchGraphicsSettings:
		return "Fetch Settings, Graphics"
	case SetGraphicsSettings:
		return "Set Settings, Graphics"
	case FetchLightingSettings:
		return "Fetch Settings, Lighting"
	case SetLightingSettings:
		return "Set Settings, lighting"

	}
	return "UNHANDLED_EVENT"
}

type event struct {
	Name          eventType
	CallerId      int64
	TimeRequested time.Time
	TraceId       uint64
	Payload       interface{}
}

type setGraphicsSettingsPayload struct {
	DispatchChannel   chan struct{}
	ShaderName        string
	GraphicsFrequency time.Duration
	ReloadOnUpdate    bool
}

type setLightingSettingsPayload struct {
	DispatchChannel   chan struct{}
	SegmentDefinition types.LedSegment
	SegmentCount      int
}
