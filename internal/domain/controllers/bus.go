package controllers

import "github.com/polis-interactive/2023-CosmicMurmur/internal/domain"

type Bus interface {
	domain.EventHandler
	StopRenderer()
	StartRenderer()
}
