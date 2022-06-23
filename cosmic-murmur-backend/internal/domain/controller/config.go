package controller

import "github.com/polis-interactive/2023-CosmicMurmur/internal/types"

type Config interface {
	GetControllerLocalAddress() string
	GetControllerNodeDefinitions() types.NodeDefinitions
}
