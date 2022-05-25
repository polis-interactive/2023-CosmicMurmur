package domain

type Event int64

const (
	ClientConnected Event = iota
	ClientDisconnected
	LightingUpdated
	RenderUpdated
	StartedDMXConsole
	StoppedDmxConsole
	GraphicsCrashed
	GraphicsReady
)

func (s Event) String() string {
	switch s {
	case ClientConnected:
		return "Client Connected"
	case ClientDisconnected:
		return "Client Disconnected"
	case LightingUpdated:
		return "Lighting Updated"
	case RenderUpdated:
		return "Render Updated"
	case StartedDMXConsole:
		return "Started DMX Console"
	case StoppedDmxConsole:
		return "Stopped DMX Console"
	case GraphicsCrashed:
		return "Graphics Crashed"
	case GraphicsReady:
		return "Graphics Ready"
	}
	return "UNHANDLED_EVENT"
}
