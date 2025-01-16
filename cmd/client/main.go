package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hannesi/udp-rdt-go/internal/config"
	"github.com/hannesi/udp-rdt-go/internal/reliability/client"
	"github.com/hannesi/udp-rdt-go/internal/virtualsocket"
)

func main() {
	fmt.Println("CLIENT")

    socket, err := virtualsocket.NewVirtualSocket()

	if err != nil {
		log.Fatalf("Failed to create UDP socket: %v", err)
	}

	defer socket.Close()

    reliabilityLayer := client.ReliabilityLayerWithBitErrorDetection {
        Socket: socket, 
    }

	fmt.Printf("Ready to send messages to %s:%d\n", config.DefaultConfig.IPAddrString, config.DefaultConfig.ServerPort)
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

		err = reliabilityLayer.Send([]byte(trimmedMsg))
		if err != nil {
			log.Printf("Failed to send the message: %v\n", err)
		}
	}

}
