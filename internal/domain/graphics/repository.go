package graphics

import "time"

type Repository interface {
	GetGraphicsReloadOnUpdate() (reloadOnUpdate bool, ok bool)
	SetGraphicsReloadOnUpdate(reloadOnUpdate bool) error
	GetGraphicsShader() (shaderName string, ok bool)
	SetGraphicsShader(shaderName string) error
	GetGraphicsFrequency() (frequency time.Duration, ok bool)
	SetGraphicsFrequency(frequency time.Duration) error
}
