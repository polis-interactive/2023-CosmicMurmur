package controller

import "github.com/polis-interactive/2023-CosmicMurmur/internal/types"

type Repository interface {
	SetControllerLocalAddress(addr string) error
	GetControllerLocalAddress() (addr string, ok bool)
	SetControllerNodeDefinitions(definitions types.NodeDefinitions) error
	GetControllerNodeDefinitions() (definitions types.NodeDefinitions, ok bool)
}
