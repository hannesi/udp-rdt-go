package config


type Config struct {
    IPAddrString string
	ServerPort int
}

var DefaultConfig = Config {
    IPAddrString: "127.0.0.1",
    ServerPort: 42069,
}

