package lighting

import "github.com/polis-interactive/2023-CosmicMurmur/internal/types"

type Config interface {
	GetLightingSegmentDefinition() types.LedSegment
	GetLightingSegmentCount() int
}
