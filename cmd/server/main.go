package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hannesi/udp-rdt-go/internal/config"
)

func main() {
    fmt.Println("SERVER")

    addr := net.UDPAddr{
        IP: net.ParseIP(config.DefaultConfig.IPAddrString),
        Port: config.DefaultConfig.ServerPort,
    }

    conn, err := net.ListenUDP("udp", &addr) 

    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Printf("UDP server is listening on %s\n", addr.String())

    buffer := make([]byte, 1024)

    for {
        n, addr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            log.Fatal(err)
            break
        }

        fmt.Printf("%s: \"%s\"\n", addr, string(buffer[:n]))

    }
}
