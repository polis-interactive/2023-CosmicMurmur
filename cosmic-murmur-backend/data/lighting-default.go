package data

import "github.com/polis-interactive/2023-CosmicMurmur/internal/types"

var DefaultNodeDefinitions = types.NodeDefinitions{
	types.NodeDefinition{
		Address:   "2.0.0.2",
		Universes: []int{0, 1, 2, 3, 4, 5},
	},
}

var DefaultLightingSegmentDefinition = types.LedSegment{
	// U1 (most of seg 1)
	types.LedUniverse{
		types.LedString{
			LedCount:    3,
			StringCount: 5,
		},
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    11,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
	},
	// U2 (remaining seg1, most of seg2)
	types.LedUniverse{
		types.LedString{
			LedCount:    3,
			StringCount: 10,
		},
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    11,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
	},
	// U3 (remaining seg2, most of seg3)
	types.LedUniverse{
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
		types.LedString{
			LedCount:    3,
			StringCount: 10,
		},
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    11,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
	},
	// U4 (remaining seg3, half of seg4)
	types.LedUniverse{
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
		types.LedString{
			LedCount:    3,
			StringCount: 10,
		},
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    11,
			StringCount: 3,
		},
	},
	// U5 (half of seg4, half of seg5)
	types.LedUniverse{
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
		types.LedString{
			LedCount:    3,
			StringCount: 10,
		},
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
	},
	// U6 (half of seg5)
	types.LedUniverse{
		types.LedString{
			LedCount:    11,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    9,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    7,
			StringCount: 3,
		},
		types.LedString{
			LedCount:    5,
			StringCount: 2,
		},
		types.LedString{
			LedCount:    3,
			StringCount: 5,
		},
	},
}
