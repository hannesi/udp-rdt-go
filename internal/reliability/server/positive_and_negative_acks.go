package server

import (
	"log"
	"net"

	"github.com/hannesi/udp-rdt-go/internal/reliability"
)

type ReliabilityLayerWithPositiveAndNegativeAcks struct {
	Socket *net.UDPConn
}

func (r *ReliabilityLayerWithPositiveAndNegativeAcks) Receive() ([]byte, error) {
	internalBuffer := make([]byte, 1024)
	n, addr, err := r.Socket.ReadFromUDP(internalBuffer)
	if err != nil {
		log.Fatal(err)
	}

	packet, err := reliability.Deserialize(internalBuffer[:n])

	if !packet.IsChecksumValid() {
		log.Println("Bit error detected! Checksum does not match.")
		log.Println("Sending NACK.")
		r.Socket.WriteToUDP([]byte("NACK"), addr)
		return r.Receive()
	}

	log.Println("Sending ACK.")
	r.Socket.WriteToUDP([]byte("ACK"), addr)

	return packet.Payload, nil
}
