package server

import (
	"log"
	"net"

	"github.com/hannesi/udp-rdt-go/internal/reliability"
)

type ReliabilityLayerWithBitErrorDetection struct {
	Socket *net.UDPConn
}

func (r *ReliabilityLayerWithBitErrorDetection) Receive() ([]byte, error) {
    internalBuffer := make([]byte, 1024)
	n, _, err := r.Socket.ReadFromUDP(internalBuffer)
	if err != nil {
		log.Fatal(err)
	}

    packet, err := reliability.Deserialize(internalBuffer[:n])

    if !packet.IsChecksumValid() {
        log.Println("Bit error detected! Checksum does not match.")
    }

	return packet.Payload, nil
}
