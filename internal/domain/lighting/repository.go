package lighting

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
)

type Repository interface {
	GetLightingSegmentDefinition() (segment types.LedSegment, ok bool)
	SetLightingSegmentDefinition(types.LedSegment) error
	GetLightingSegmentCount() (count int, ok bool)
	SetLightingSegmentCount(count int) error
}
