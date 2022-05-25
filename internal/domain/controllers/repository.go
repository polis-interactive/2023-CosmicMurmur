package controllers

import "github.com/polis-interactive/2023-CosmicMurmur/internal/types"

type Repository interface {
	SetControllerLocalAddress(addr string) error
	GetControllerLocalAddress() (addr string, ok bool)
	SetControllerDefinitions(definitions types.ControllerDefinitions) error
	GetControllerDefinitions() (definitions types.ControllerDefinitions, ok bool)
}
