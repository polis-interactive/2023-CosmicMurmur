package application

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"time"
)

type LightingConfig struct {
	SegmentDefinition types.LedSegment
	SegmentCount      int
}

func (c *LightingConfig) GetLightingSegmentDefinition() types.LedSegment {
	return c.SegmentDefinition
}

func (c *LightingConfig) GetLightingSegmentCount() int {
	return c.SegmentCount
}

type GraphicsConfig struct {
	DefaultShader  string
	PixelSize      int
	Frequency      time.Duration
	ReloadOnUpdate bool
}

func (c *GraphicsConfig) GetGraphicsDefaultShader() string {
	return c.DefaultShader
}

func (c *GraphicsConfig) GetGraphicsPixelSize() int {
	return c.PixelSize
}

func (c *GraphicsConfig) GetGraphicsFrequency() time.Duration {
	return c.Frequency
}

func (c *GraphicsConfig) GetGraphicsReloadOnUpdate() bool {
	return c.ReloadOnUpdate
}

type ControllerConfig struct {
	LocalAddress    string
	NodeDefinitions types.NodeDefinitions
}

func (c *ControllerConfig) GetControllerLocalAddress() string {
	return c.LocalAddress
}

func (c *ControllerConfig) GetControllerNodeDefinitions() types.NodeDefinitions {
	return c.NodeDefinitions
}

type ServiceBusConfig struct {
	EventQueueSize int
	BusyTimeout    time.Duration
}

func (c *ServiceBusConfig) GetServiceBusEventQueueSize() int {
	return c.EventQueueSize
}

func (c *ServiceBusConfig) GetServiceBusBusyTimeout() time.Duration {
	return c.BusyTimeout
}

type Config struct {
	*LightingConfig
	*GraphicsConfig
	*ControllerConfig
	*ServiceBusConfig
	ProgramName string
}

func (c *Config) GetProgramName() string {
	return c.ProgramName
}
