package graphics

import "time"

type Config interface {
	GetProgramName() string
	GetGraphicsDefaultShader() string
	GetGraphicsPixelSize() int
	GetGraphicsFrequency() time.Duration
	GetGraphicsReloadOnUpdate() bool
}
