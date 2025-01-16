package config

import "time"

type Config struct {
	IPAddrString           string
	ServerPort             int
	VirtualSocketDelayRate float64
	VirtualSocketDelay     time.Duration
	VirtualSocketDropRate  float64
	VirtualSocketErrorRate float64
}

var DefaultConfig = Config{
	IPAddrString:           "127.0.0.1",
	ServerPort:             42069,
	VirtualSocketDelayRate: 0.0,
	VirtualSocketDelay:     500 * time.Millisecond,
	VirtualSocketDropRate:  0.0,
	VirtualSocketErrorRate: 0.25,
}
