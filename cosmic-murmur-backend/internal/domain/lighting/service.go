package lighting

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/domain"
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"log"
)

type service struct {
	repo Repository
	cfg  Config

	segmentDefinition types.LedSegment
	segmentCount      int

	universeLights [][]*types.Light
	grid           *types.Grid
}

var _ domain.LightingService = (*service)(nil)

func NewService(cfg Config, repo Repository) *service {
	log.Println("Lighting, NewService: creating")
	s := &service{
		repo:              repo,
		cfg:               cfg,
		segmentDefinition: nil,
		segmentCount:      0,
		universeLights:    nil,
		grid:              nil,
	}
	s.SetupLightingService()
	return s
}

func (s *service) SetupLightingService() {
	s.initializeVariables()
	s.doCreateLights()
}

func (s *service) initializeVariables() {
	var ok bool
	var segmentDefinition types.LedSegment
	segmentDefinition, ok = s.repo.GetLightingSegmentDefinition()
	if !ok {
		log.Println("Lighting, initializeVariables: no segment definition found, using default")
		segmentDefinition = s.cfg.GetLightingSegmentDefinition()
	}
	s.segmentDefinition = segmentDefinition
	var segmentCount int
	segmentCount, ok = s.repo.GetLightingSegmentCount()
	if !ok {
		log.Println("Lighting, initializeVariables: no segment count found, using default")
		segmentCount = s.cfg.GetLightingSegmentCount()
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
	maxAboveBelow := 0
	stringsPerSegment := 0
	universesPerSegment := len(s.segmentDefinition)
	universeLights := make([][]*types.Light, s.segmentCount*len(s.segmentDefinition))
	for segment := 0; segment < s.segmentCount; segment++ {
		// keeps track of strings in the current section
		seenStrings := 0
		for universeNumber, universe := range s.segmentDefinition {
			lights := make([]*types.Light, 0, 171)
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
							Color:    types.Color{},
						}
						nextPixel += 1
						lights = append(lights, newLight)
					}
					evenString = !evenString
				}
				seenStrings += ledString.StringCount
			}
			universeLights[universeNumber+segment*universesPerSegment] = lights
		}
		/* this will technically only do something on the first iteration;
		subsequent segments will be copies of the original
		*/
		stringsPerSegment = seenStrings
	}
	s.universeLights = universeLights
	s.grid = &types.Grid{
		MinX: 0,
		MaxX: stringsPerSegment*s.segmentCount - 1,
		MinY: -maxAboveBelow,
		MaxY: maxAboveBelow,
	}
}

func (s *service) GetLightUniverses() [][]*types.Light {
	return s.universeLights
}

func (s *service) GetSettings() *domain.LightingSettings {
	return &domain.LightingSettings{
		SegmentDefinition: s.segmentDefinition,
		SegmentCount:      s.segmentCount,
	}
}

func (s *service) SetSettings(settings *domain.LightingSettings) error {
	err := s.repo.SetLightingSegmentDefinition(settings.SegmentDefinition)
	if err != nil {
		return err
	}
	err = s.repo.SetLightingSegmentCount(settings.SegmentCount)
	if err != nil {
		return err
	}
	s.segmentDefinition = settings.SegmentDefinition
	s.segmentCount = settings.SegmentCount
	s.doCreateLights()
	return nil
}

func (s *service) GetGridDimensions() *types.Grid {
	return s.grid
}
