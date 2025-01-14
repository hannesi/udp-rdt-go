package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/hannesi/udp-rdt-go/internal/config"
)

func main() {
	fmt.Println("CLIENT")

	destAddr := net.UDPAddr{
		IP:   net.ParseIP(config.DefaultConfig.IPAddrString),
		Port: config.DefaultConfig.ServerPort,
	}

	socketAddr := net.UDPAddr{
		IP:   net.ParseIP(config.DefaultConfig.IPAddrString),
		Port: 0,
	}

    socket, err := net.DialUDP("udp", &socketAddr, &destAddr)

    if err != nil {
        log.Fatalf("Failed to create UDP socket: %v", err)
    }

    defer socket.Close()

    fmt.Printf("Ready to send messages to %s:%d\n", destAddr.IP, destAddr.Port)
    fmt.Println("Usage: Type a message and hit enter :)")

    reader := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("Enter message: ")
        msg, err := reader.ReadString('\n')
        if err != nil {
            log.Println("Failed to read input:", err)
            continue
        }

        trimmedMsg := strings.TrimSpace(msg)

        _, err = socket.Write([]byte(trimmedMsg))
        		if err != nil {
			log.Printf("Failed to send the message: %v\n", err)
		} else {
			fmt.Printf("Sending \"%s\" to %s:%d\n", trimmedMsg, destAddr.IP, destAddr.Port)
		}
    }

}
