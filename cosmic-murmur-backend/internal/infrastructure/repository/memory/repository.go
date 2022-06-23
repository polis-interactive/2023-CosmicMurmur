package memory

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
	"sync"
	"time"
)

var (
	defaultRepository = Repository{
		lightingSegmentDefinition: nil,
		lightingSegmentCount:      -1,
		graphicsReloadOnUpdate:    -1,
		graphicsShaderName:        "",
		graphicsFrequency:         nil,
		controllerLocalAddress:    "",
		controllerNodeDefinitions: nil,
		mu:                        &sync.RWMutex{},
	}
)

type Repository struct {
	lightingSegmentDefinition *types.LedSegment
	lightingSegmentCount      int
	graphicsReloadOnUpdate    int
	graphicsShaderName        string
	graphicsFrequency         *time.Duration
	controllerLocalAddress    string
	controllerNodeDefinitions types.NodeDefinitions
	mu                        *sync.RWMutex
}

func NewMemoryRepository() *Repository {
	r := &Repository{
		mu: &sync.RWMutex{},
	}
	*r = defaultRepository
	return r
}

func (r *Repository) ResetRepository() error {
	*r = defaultRepository
	return nil
}

func (r *Repository) GetLightingSegmentDefinition() (segment types.LedSegment, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.lightingSegmentDefinition != nil {
		return *r.lightingSegmentDefinition, true
	} else {
		return types.LedSegment{}, false
	}
}

func (r *Repository) SetLightingSegmentDefinition(segment types.LedSegment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	*r.lightingSegmentDefinition = segment
	return nil
}

func (r *Repository) GetLightingSegmentCount() (count int, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.lightingSegmentCount != -1 {
		return r.lightingSegmentCount, true
	} else {
		return -1, false
	}
}

func (r *Repository) SetLightingSegmentCount(count int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lightingSegmentCount = count
	return nil
}

func (r *Repository) GetGraphicsReloadOnUpdate() (reloadOnUpdate bool, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.graphicsReloadOnUpdate == 1 {
		return true, true
	} else if r.graphicsReloadOnUpdate == 0 {
		return false, true
	} else {
		return false, false
	}
}

func (r *Repository) SetGraphicsReloadOnUpdate(reloadOnUpdate bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if reloadOnUpdate {
		r.graphicsReloadOnUpdate = 1
	} else {
		r.graphicsReloadOnUpdate = 0
	}
	return nil
}

func (r *Repository) GetGraphicsShader() (shaderName string, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.graphicsShaderName != "" {
		return r.graphicsShaderName, true
	} else {
		return "", false
	}
}

func (r *Repository) SetGraphicsShader(shaderName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.graphicsShaderName = shaderName
	return nil
}

func (r *Repository) GetGraphicsFrequency() (frequency time.Duration, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.graphicsFrequency != nil {
		return *r.graphicsFrequency, true
	} else {
		return 0 * time.Millisecond, false
	}
}

func (r *Repository) SetGraphicsFrequency(frequency time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	*r.graphicsFrequency = frequency
	return nil
}

func (r *Repository) SetControllerLocalAddress(addr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.controllerLocalAddress = addr
	return nil
}

func (r *Repository) GetControllerLocalAddress() (addr string, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.controllerLocalAddress != "" {
		return r.controllerLocalAddress, true
	} else {
		return "", false
	}
}

func (r *Repository) SetControllerNodeDefinitions(definitions types.NodeDefinitions) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy(r.controllerNodeDefinitions, definitions)
	return nil
}

func (r *Repository) GetControllerNodeDefinitions() (definitions types.NodeDefinitions, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.controllerNodeDefinitions != nil {
		copy(definitions, r.controllerNodeDefinitions)
		return definitions, true
	} else {
		return nil, false
	}
}
