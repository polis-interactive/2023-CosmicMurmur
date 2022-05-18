package domain

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"sync"
)

type GraphicsService interface {
	Startup()
	Reset()
	Shutdown()
	GetPb() (pb *types.PixelBuffer, preLockedMutex *sync.RWMutex)
}

type EventHandler interface {
	SubscribeToEvent(e Event, fn func()) int
	UnsubscribeToEvent(handlerId int)
	HandleEvent(e Event)
}
