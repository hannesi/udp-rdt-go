package server

import (
	"log"
	"net"

	"github.com/hannesi/udp-rdt-go/internal/reliability"
)

type ReliableTransportProtocol struct {
	Socket    *net.UDPConn
	lastOkAck uint8
}

func NewReliableTransportProtocol(socket *net.UDPConn) *ReliableTransportProtocol {
	return &ReliableTransportProtocol{
		Socket:    socket,
		lastOkAck: ^uint8(0),
	}
}

func (r *ReliableTransportProtocol) Receive() ([]byte, error) {
	internalBuffer := make([]byte, 1024)
	n, addr, err := r.Socket.ReadFromUDP(internalBuffer)
	if err != nil {
		log.Fatal(err)
	}

	packet, err := reliability.Deserialize(internalBuffer[:n])

	// Skips duplicates and e.g. data of seq 6 if expecting 5
	if r.lastOkAck + 1 != packet.Sequence {
		log.Printf("Unexpected sequence number! Skipping data: %s", packet.Payload)
		return r.Receive()
	}

	if !packet.IsChecksumValid() {
		log.Println("Bit error detected! Checksum does not match.")
		ackMsg, _ := reliability.SerializeAckData(packet.Sequence, "NAK")
		log.Printf("Sending NAK %d.", packet.Sequence)
		r.Socket.WriteToUDP([]byte(ackMsg), addr)
		return r.Receive()
	}

	r.lastOkAck = packet.Sequence

	log.Printf("Sending ACK %d.", r.lastOkAck)
	ackMsg, _ := reliability.SerializeAckData(r.lastOkAck, "ACK")
	r.Socket.WriteToUDP([]byte(ackMsg), addr)

	// if packet has same seq as last ok, don't return the duplicate message.
	return packet.Payload, nil
}
