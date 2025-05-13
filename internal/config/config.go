package config

import "time"

type Config struct {
	IPAddrString           string
	ServerPort             int
	VirtualSocketDelayRate float64
	VirtualSocketDelay     time.Duration
	VirtualSocketDropRate  float64
	VirtualSocketErrorRate float64
    ReliabilityLayerAckWaitTime time.Duration
}

var DefaultConfig = Config{
	IPAddrString:           "127.0.0.1",
	ServerPort:             42069,
	VirtualSocketDelayRate: 0.2,
	VirtualSocketDelay:     1500 * time.Millisecond,
	VirtualSocketDropRate:  0.2,
	VirtualSocketErrorRate: 0.2,
    ReliabilityLayerAckWaitTime: 1000 * time.Millisecond,
}
