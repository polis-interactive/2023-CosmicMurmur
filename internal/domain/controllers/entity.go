package controllers

import "github.com/polis-interactive/2023-CosmicMurmur/internal/types"

type controller struct {
	address       string
	rawUniverses  []int
	dataUniverses map[int][512]byte
	needsPoll     bool
}

func newController(definition types.ControllerDefinition) *controller {
	c := &controller{
		address:       definition.Address,
		rawUniverses:  definition.Universes,
		dataUniverses: make(map[int][512]byte),
		needsPoll:     false,
	}
	for _, u := range definition.Universes {
		c.dataUniverses[u] = [512]byte{}
	}
	return c
}

func (c *controller) getDefinition() types.ControllerDefinition {
	return types.ControllerDefinition{
		Address:   c.address,
		Universes: c.rawUniverses,
	}
}

func mapControllersToDefinitions(controllers []*controller) types.ControllerDefinitions {
	definitions := make(types.ControllerDefinitions, len(controllers))
	for i, c := range controllers {
		definitions[i] = c.getDefinition()
	}
	return definitions
}
