package config

type Config struct {
	IPAddrString           string
	ServerPort             int
	VirtualSocketDelayRate float64
	VirtualSocketDelayMs   int
	VirtualSocketDropRate  float64
	VirtualSocketErrorRate float64
}

var DefaultConfig = Config{
	IPAddrString:           "127.0.0.1",
	ServerPort:             42069,
	VirtualSocketDelayRate: 0.25,
	VirtualSocketDelayMs:   1000,
	VirtualSocketDropRate:  0.25,
	VirtualSocketErrorRate: 0.25,
}
