package graphics

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
)

type Bus interface {
	GetGridDimensions() *types.Grid
	EmitGraphicsCrashed()
	EmitGraphicsReady()
}
