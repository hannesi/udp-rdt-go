package server

import (
	"log"
	"net"

	"github.com/hannesi/udp-rdt-go/internal/reliability"
)

type ReliabilityLayerWithPositiveAcks struct {
	Socket    *net.UDPConn
	lastOkAck uint8
}

func NewReliabilityLayerWithPositiveAcks(socket *net.UDPConn) *ReliabilityLayerWithPositiveAcks {
    return &ReliabilityLayerWithPositiveAcks{
        Socket: socket,
        lastOkAck: ^uint8(0),
    }
}

func (r *ReliabilityLayerWithPositiveAcks) Receive() ([]byte, error) {
	internalBuffer := make([]byte, 1024)
	n, addr, err := r.Socket.ReadFromUDP(internalBuffer)
	if err != nil {
		log.Fatal(err)
	}

	packet, err := reliability.Deserialize(internalBuffer[:n])

	if !packet.IsChecksumValid() {
		log.Println("Bit error detected! Checksum does not match.")
		log.Printf("Sending ACK %d.\n", r.lastOkAck)
        ackMsg, _ := reliability.SerializeAckData(r.lastOkAck, "ACK")
		r.Socket.WriteToUDP([]byte(ackMsg), addr)
		return r.Receive()
	}

    r.lastOkAck = packet.Sequence

	log.Printf("Sending ACK %d.", r.lastOkAck)
    ackMsg, _ := reliability.SerializeAckData(r.lastOkAck, "ACK")
    r.Socket.WriteToUDP([]byte(ackMsg), addr)

	return packet.Payload, nil
}
