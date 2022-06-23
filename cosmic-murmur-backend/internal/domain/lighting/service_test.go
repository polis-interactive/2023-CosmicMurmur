package lighting

import (
	"fmt"
	"github.com/polis-interactive/2023-CosmicMurmur/data"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"log"
	"testing"
)

func testLightsEq(a, b [][]*types.Light) bool {
	if len(a) != len(b) {
		log.Println(fmt.Sprintf("len(a) = %d, len(b) = %d", len(a), len(b)))
		return false
	}
	for i, universe := range a {
		for j := range universe {
			if a[i][j].Pixel != b[i][j].Pixel {
				log.Println(fmt.Sprintf(
					"pixel a: %s, pixel b: %s", a[i][j].Print(), b[i][j].Print(),
				))
				return false
			} else if a[i][j].Position.X != b[i][j].Position.X {
				log.Println(fmt.Sprintf(
					"pixel a: %s, pixel b: %s", a[i][j].Print(), b[i][j].Print(),
				))
				return false
			} else if a[i][j].Position.Y != b[i][j].Position.Y {
				log.Println(fmt.Sprintf(
					"pixel a: %s, pixel b: %s", a[i][j].Print(), b[i][j].Print(),
				))
				return false
			}
		}

	}
	return true
}

func testGridEq(a, b *types.Grid) bool {
	return a.MinY == b.MinY && a.MinX == b.MinX &&
		a.MaxY == b.MaxY && a.MaxX == b.MaxX
}

func TestService_doCreateLights(t *testing.T) {
	s1 := &service{
		segmentDefinition: types.LedSegment{
			types.LedUniverse{
				types.LedString{
					LedCount:    3,
					StringCount: 3,
				},
			},
		},
		segmentCount: 1,
	}
	s1.doCreateLights()
	universeLights1 := [][]*types.Light{
		{
			{
				Position: types.Point{
					X: 0,
					Y: -1,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 0,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 1,
				},
				Pixel: 2,
			},
			{
				Position: types.Point{
					X: 1,
					Y: 1,
				},
				Pixel: 3,
			},
			{
				Position: types.Point{
					X: 1,
					Y: 0,
				},
				Pixel: 4,
			},
			{
				Position: types.Point{
					X: 1,
					Y: -1,
				},
				Pixel: 5,
			},
			{
				Position: types.Point{
					X: 2,
					Y: -1,
				},
				Pixel: 6,
			},
			{
				Position: types.Point{
					X: 2,
					Y: 0,
				},
				Pixel: 7,
			},
			{
				Position: types.Point{
					X: 2,
					Y: 1,
				},
				Pixel: 8,
			},
		},
	}
	grid1 := &types.Grid{
		MinX: 0,
		MaxX: 2,
		MinY: -1,
		MaxY: 1,
	}
	if !testLightsEq(s1.universeLights, universeLights1) {
		t.Fatal("Lighting array 1 does not match template")
	} else if !testGridEq(s1.grid, grid1) {
		t.Fatal("Lighting grid 1 does not match template")
	}

	s2 := &service{
		segmentDefinition: types.LedSegment{
			types.LedUniverse{
				types.LedString{
					LedCount:    3,
					StringCount: 1,
				},
				types.LedString{
					LedCount:    5,
					StringCount: 2,
				},
				types.LedString{
					LedCount:    3,
					StringCount: 1,
				},
			},
		},
		segmentCount: 1,
	}
	s2.doCreateLights()
	universeLights2 := [][]*types.Light{
		{
			{
				Position: types.Point{
					X: 0,
					Y: -1,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 0,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 1,
				},
				Pixel: 2,
			},
			{
				Position: types.Point{
					X: 1,
					Y: 2,
				},
				Pixel: 3,
			},
			{
				Position: types.Point{
					X: 1,
					Y: 1,
				},
				Pixel: 4,
			},
			{
				Position: types.Point{
					X: 1,
					Y: 0,
				},
				Pixel: 5,
			},
			{
				Position: types.Point{
					X: 1,
					Y: -1,
				},
				Pixel: 6,
			},
			{
				Position: types.Point{
					X: 1,
					Y: -2,
				},
				Pixel: 7,
			},
			{
				Position: types.Point{
					X: 2,
					Y: -2,
				},
				Pixel: 8,
			},
			{
				Position: types.Point{
					X: 2,
					Y: -1,
				},
				Pixel: 9,
			},
			{
				Position: types.Point{
					X: 2,
					Y: 0,
				},
				Pixel: 10,
			},
			{
				Position: types.Point{
					X: 2,
					Y: 1,
				},
				Pixel: 11,
			},
			{
				Position: types.Point{
					X: 2,
					Y: 2,
				},
				Pixel: 12,
			},
			{
				Position: types.Point{
					X: 3,
					Y: 1,
				},
				Pixel: 13,
			},
			{
				Position: types.Point{
					X: 3,
					Y: 0,
				},
				Pixel: 14,
			},
			{
				Position: types.Point{
					X: 3,
					Y: -1,
				},
				Pixel: 15,
			},
		},
	}

	grid2 := &types.Grid{
		MinX: 0,
		MaxX: 3,
		MinY: -2,
		MaxY: 2,
	}
	if !testLightsEq(s2.universeLights, universeLights2) {
		t.Fatal("Lighting array 2 does not match template")
	} else if !testGridEq(s2.grid, grid2) {
		t.Fatal("Lighting grid 2 does not match template")
	}

	s3 := &service{
		segmentDefinition: types.LedSegment{
			types.LedUniverse{
				types.LedString{
					LedCount:    3,
					StringCount: 1,
				},
			},
			types.LedUniverse{
				types.LedString{
					LedCount:    3,
					StringCount: 2,
				},
			},
			types.LedUniverse{
				types.LedString{
					LedCount:    3,
					StringCount: 1,
				},
			},
		},
		segmentCount: 1,
	}
	s3.doCreateLights()
	universeLights3 := [][]*types.Light{
		{
			{
				Position: types.Point{
					X: 0,
					Y: -1,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 0,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 1,
				},
				Pixel: 2,
			},
		},
		{

			{
				Position: types.Point{
					X: 1,
					Y: -1,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 1,
					Y: 0,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 1,
					Y: 1,
				},
				Pixel: 2,
			},
			{
				Position: types.Point{
					X: 2,
					Y: 1,
				},
				Pixel: 3,
			},
			{
				Position: types.Point{
					X: 2,
					Y: 0,
				},
				Pixel: 4,
			},
			{
				Position: types.Point{
					X: 2,
					Y: -1,
				},
				Pixel: 5,
			},
		},
		{
			{
				Position: types.Point{
					X: 3,
					Y: -1,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 3,
					Y: 0,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 3,
					Y: 1,
				},
				Pixel: 2,
			},
		},
	}
	grid3 := &types.Grid{
		MinX: 0,
		MaxX: 3,
		MinY: -1,
		MaxY: 1,
	}
	if !testLightsEq(s3.universeLights, universeLights3) {
		t.Fatal("Lighting array 3 does not match template")
	} else if !testGridEq(s3.grid, grid3) {
		t.Fatal("Lighting grid 3 does not match template")
	}

	s4 := &service{
		segmentDefinition: types.LedSegment{
			types.LedUniverse{
				types.LedString{
					LedCount:    7,
					StringCount: 1,
				},
				types.LedString{
					LedCount:    3,
					StringCount: 2,
				},
			},
			types.LedUniverse{
				types.LedString{
					LedCount:    5,
					StringCount: 1,
				},
			},
		},
		segmentCount: 3,
	}
	s4.doCreateLights()
	universeLights4 := [][]*types.Light{
		{
			{
				Position: types.Point{
					X: 0,
					Y: -3,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 0,
					Y: -2,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 0,
					Y: -1,
				},
				Pixel: 2,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 0,
				},
				Pixel: 3,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 1,
				},
				Pixel: 4,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 2,
				},
				Pixel: 5,
			},
			{
				Position: types.Point{
					X: 0,
					Y: 3,
				},
				Pixel: 6,
			},
			{
				Position: types.Point{
					X: 1,
					Y: 1,
				},
				Pixel: 7,
			},
			{
				Position: types.Point{
					X: 1,
					Y: 0,
				},
				Pixel: 8,
			},
			{
				Position: types.Point{
					X: 1,
					Y: -1,
				},
				Pixel: 9,
			},
			{
				Position: types.Point{
					X: 2,
					Y: -1,
				},
				Pixel: 10,
			},
			{
				Position: types.Point{
					X: 2,
					Y: 0,
				},
				Pixel: 11,
			},
			{
				Position: types.Point{
					X: 2,
					Y: 1,
				},
				Pixel: 12,
			},
		},
		{
			{
				Position: types.Point{
					X: 3,
					Y: -2,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 3,
					Y: -1,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 3,
					Y: 0,
				},
				Pixel: 2,
			},
			{
				Position: types.Point{
					X: 3,
					Y: 1,
				},
				Pixel: 3,
			},
			{
				Position: types.Point{
					X: 3,
					Y: 2,
				},
				Pixel: 4,
			},
		},
		{
			{
				Position: types.Point{
					X: 4,
					Y: -3,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 4,
					Y: -2,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 4,
					Y: -1,
				},
				Pixel: 2,
			},
			{
				Position: types.Point{
					X: 4,
					Y: 0,
				},
				Pixel: 3,
			},
			{
				Position: types.Point{
					X: 4,
					Y: 1,
				},
				Pixel: 4,
			},
			{
				Position: types.Point{
					X: 4,
					Y: 2,
				},
				Pixel: 5,
			},
			{
				Position: types.Point{
					X: 4,
					Y: 3,
				},
				Pixel: 6,
			},
			{
				Position: types.Point{
					X: 5,
					Y: 1,
				},
				Pixel: 7,
			},
			{
				Position: types.Point{
					X: 5,
					Y: 0,
				},
				Pixel: 8,
			},
			{
				Position: types.Point{
					X: 5,
					Y: -1,
				},
				Pixel: 9,
			},
			{
				Position: types.Point{
					X: 6,
					Y: -1,
				},
				Pixel: 10,
			},
			{
				Position: types.Point{
					X: 6,
					Y: 0,
				},
				Pixel: 11,
			},
			{
				Position: types.Point{
					X: 6,
					Y: 1,
				},
				Pixel: 12,
			},
		},
		{
			{
				Position: types.Point{
					X: 7,
					Y: -2,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 7,
					Y: -1,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 7,
					Y: 0,
				},
				Pixel: 2,
			},
			{
				Position: types.Point{
					X: 7,
					Y: 1,
				},
				Pixel: 3,
			},
			{
				Position: types.Point{
					X: 7,
					Y: 2,
				},
				Pixel: 4,
			},
		},
		{
			{
				Position: types.Point{
					X: 8,
					Y: -3,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 8,
					Y: -2,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 8,
					Y: -1,
				},
				Pixel: 2,
			},
			{
				Position: types.Point{
					X: 8,
					Y: 0,
				},
				Pixel: 3,
			},
			{
				Position: types.Point{
					X: 8,
					Y: 1,
				},
				Pixel: 4,
			},
			{
				Position: types.Point{
					X: 8,
					Y: 2,
				},
				Pixel: 5,
			},
			{
				Position: types.Point{
					X: 8,
					Y: 3,
				},
				Pixel: 6,
			},
			{
				Position: types.Point{
					X: 9,
					Y: 1,
				},
				Pixel: 7,
			},
			{
				Position: types.Point{
					X: 9,
					Y: 0,
				},
				Pixel: 8,
			},
			{
				Position: types.Point{
					X: 9,
					Y: -1,
				},
				Pixel: 9,
			},
			{
				Position: types.Point{
					X: 10,
					Y: -1,
				},
				Pixel: 10,
			},
			{
				Position: types.Point{
					X: 10,
					Y: 0,
				},
				Pixel: 11,
			},
			{
				Position: types.Point{
					X: 10,
					Y: 1,
				},
				Pixel: 12,
			},
		},
		{
			{
				Position: types.Point{
					X: 11,
					Y: -2,
				},
				Pixel: 0,
			},
			{
				Position: types.Point{
					X: 11,
					Y: -1,
				},
				Pixel: 1,
			},
			{
				Position: types.Point{
					X: 11,
					Y: 0,
				},
				Pixel: 2,
			},
			{
				Position: types.Point{
					X: 11,
					Y: 1,
				},
				Pixel: 3,
			},
			{
				Position: types.Point{
					X: 11,
					Y: 2,
				},
				Pixel: 4,
			},
		},
	}
	grid4 := &types.Grid{
		MinX: 0,
		MaxX: 11,
		MinY: -3,
		MaxY: 3,
	}
	if !testLightsEq(s4.universeLights, universeLights4) {
		t.Fatal("Lighting array 4 does not match template")
	} else if !testGridEq(s4.grid, grid4) {
		t.Fatal("Lighting grid 4 does not match template")
	}

	s5 := &service{
		segmentDefinition: data.DefaultLightingSegmentDefinition,
		segmentCount:      1,
	}
	s5.doCreateLights()
	grid5 := &types.Grid{
		MinX: 0,
		MaxX: 144,
		MinY: -5,
		MaxY: 5,
	}
	if !testGridEq(s5.grid, grid5) {
		t.Fatal("Lighting grid 5 does not match template")
	}
}
