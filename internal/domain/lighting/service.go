package lighting

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"log"
	"sync"
)

type service struct {
	repo Repository
	bus  Bus

	mu                *sync.RWMutex
	segmentDefinition types.LedSegment
	segmentCount      int

	lights []*types.Light
	grid   *types.Grid
}

var _ domain.LightingService = (*service)(nil)

func NewService(cfg Config, repo Repository, bus Bus) *service {
	log.Println("Lighting, NewService: creating")
	s := &service{
		repo:              repo,
		bus:               bus,
		mu:                &sync.RWMutex{},
		segmentDefinition: nil,
		segmentCount:      0,
		lights:            nil,
		grid:              nil,
	}
	s.initializeVariables(cfg)
	s.doCreateLights()
	return s
}

func (s *service) initializeVariables(cfg Config) {
	var ok bool
	var segmentDefinition types.LedSegment
	segmentDefinition, ok = s.repo.GetLightingSegmentDefinition()
	if !ok {
		log.Println("Lighting, initializeVariables: no segment definition found, using default")
		segmentDefinition = cfg.GetLightingSegmentDefinition()
	}
	s.segmentDefinition = segmentDefinition
	var segmentCount int
	segmentCount, ok = s.repo.GetLightingSegmentCount()
	if !ok {
		log.Println("Lighting, initializeVariables: no segment count found, using default")
		segmentCount = cfg.GetLightingSegmentCount()
	}
	s.segmentCount = segmentCount
}

func MaxInt(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func (s *service) doCreateLights() {
	lights := make([]*types.Light, 0)
	maxAboveBelow := 0
	stringsPerSegment := 0
	universesPerSegment := len(s.segmentDefinition)
	for segment := 0; segment < s.segmentCount; segment++ {
		// keeps track of strings in the current section
		seenStrings := 0
		for universeNumber, universe := range s.segmentDefinition {
			// first string in universe always starts on the bottom
			evenString := true
			nextPixel := 0
			for _, ledString := range universe {
				// strings are always odd numbered
				aboveBelow := (ledString.LedCount - 1) / 2
				maxAboveBelow = MaxInt(maxAboveBelow, aboveBelow)
				for stringNumber := 0; stringNumber < ledString.StringCount; stringNumber++ {
					/* addition of stringsPerSegment will always be zero
					on the first segment; we then set it after the first iteration
					where it actually does something
					*/
					stringXPosition := stringsPerSegment*segment + seenStrings + stringNumber
					for ledNumber := -aboveBelow; ledNumber <= aboveBelow; ledNumber++ {
						ledYPosition := ledNumber
						// on odd strings, flip the position order so they "snake"
						if !evenString {
							ledYPosition = -ledYPosition
						}
						newLight := &types.Light{
							Position: types.CreatePoint(stringXPosition, ledYPosition),
							Pixel:    nextPixel,
							Universe: universeNumber + segment*universesPerSegment,
							Color:    types.Color{},
						}
						nextPixel += 1
						lights = append(lights, newLight)
					}
					evenString = !evenString
				}
				seenStrings += ledString.StringCount
			}
		}
		/* this will technically only do something on the first iteration;
		subsequent segments will be copies of the original
		*/
		stringsPerSegment = seenStrings
	}
	s.lights = lights
	s.grid = &types.Grid{
		MinX: 0,
		MaxX: stringsPerSegment*s.segmentCount - 1,
		MinY: -maxAboveBelow,
		MaxY: maxAboveBelow,
	}
}

func (s *service) GetSettings() *domain.LightingSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &domain.LightingSettings{
		SegmentDefinition: s.segmentDefinition,
		SegmentCount:      s.segmentCount,
	}
}

func (s *service) SetSegmentDefinition(segment types.LedSegment) error {
	err := s.repo.SetLightingSegmentDefinition(segment)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.segmentDefinition = segment
	s.doCreateLights()
	s.mu.Unlock()
	s.bus.HandleEvent(domain.LightingUpdated)
	return nil
}

func (s *service) SetSegmentCount(u int) error {
	err := s.repo.SetLightingSegmentCount(u)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.segmentCount = u
	s.doCreateLights()
	s.mu.Unlock()
	s.bus.HandleEvent(domain.LightingUpdated)
	return nil
}

func (s *service) GetGridDimensions() *types.Grid {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.grid
}
