package controller

import (
	"github.com/polis-interactive/2023-CosmicMurmur/internal/types"
)

type controller struct {
	localAddress string

	nodes             []*node
	universeBufferMap map[int]*[512]byte
	universeNodeMap   map[int]*node
}

func newController(localAddress string, definitions types.NodeDefinitions) *controller {
	c := &controller{
		localAddress:      localAddress,
		nodes:             make([]*node, 0, len(definitions)),
		universeBufferMap: make(map[int]*[512]byte),
		universeNodeMap:   make(map[int]*node),
	}
	for _, nodeDefinition := range definitions {
		n := newNode(c, nodeDefinition)
		c.nodes = append(c.nodes, n)
	}
	return c
}
