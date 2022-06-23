package service

import "time"

type Config interface {
	GetServiceBusEventQueueSize() int
	GetServiceBusBusyTimeout() time.Duration
}
