package server

import (
	"log"
	"net"

	"github.com/hannesi/udp-rdt-go/internal/reliability"
)

type ReliabilityLayerWithNegativeAcks struct {
	Socket     *net.UDPConn
	lastBadAck uint8
}

func NewReliabilityLayerWithNegativeAcks(socket *net.UDPConn) *ReliabilityLayerWithNegativeAcks {
	return &ReliabilityLayerWithNegativeAcks{
		Socket:     socket,
		lastBadAck: ^uint8(0),
	}
}

func (r *ReliabilityLayerWithNegativeAcks) Receive() ([]byte, error) {
	internalBuffer := make([]byte, 1024)
	n, addr, err := r.Socket.ReadFromUDP(internalBuffer)
	if err != nil {
		log.Fatal(err)
	}

	packet, err := reliability.Deserialize(internalBuffer[:n])

	if packet.IsChecksumValid() {
		log.Println("Checksum OK! Staying silent.")
		return packet.Payload, nil
	}

	r.lastBadAck = packet.Sequence

	log.Printf("Sending NAK %d.", r.lastBadAck)
	ackMsg, _ := reliability.SerializeAckData(r.lastBadAck, "NAK")
	r.Socket.WriteToUDP([]byte(ackMsg), addr)

	return r.Receive()
}
