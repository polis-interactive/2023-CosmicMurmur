package api

type Config interface {
	GetWebServerPort() int
	GetWebServerRootDirectory() string
	GetProgramName() string
	GetWebServerIsProduction() bool
}
