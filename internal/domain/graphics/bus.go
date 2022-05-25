package graphics

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
)

type Bus interface {
	domain.EventBus
	GetGridDimensions() *types.Grid
}
