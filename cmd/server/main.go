package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hannesi/udp-rdt-go/internal/config"
	"github.com/hannesi/udp-rdt-go/internal/reliability/server"
)

func main() {
    fmt.Println("SERVER")

    addr := net.UDPAddr{
        IP: net.ParseIP(config.DefaultConfig.IPAddrString),
        Port: config.DefaultConfig.ServerPort,
    }

    socket, err := net.ListenUDP("udp", &addr) 

    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    defer socket.Close()

    reliabilityLayer := server.ReliabilityLayerWithPositiveAndNegativeAcks{
        Socket: socket,
    }

    fmt.Printf("UDP server is listening on %s\n", addr.String())

    for {
        buffer, err := reliabilityLayer.Receive()
        if err != nil {
            log.Fatal(err)
            break
        }

        fmt.Printf("%s: \"%s\"\n", addr.String(), string(buffer))
    }
}
