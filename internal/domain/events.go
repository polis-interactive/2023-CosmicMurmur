package domain

type Event int64

const (
	ClientConnected Event = iota
	ClientDisconnected
	LightingUpdated
	RenderUpdated
	StartedDMXConsole
	StoppedDmxConsole
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
	}
	return "UNHANDLED_EVENT"
}
