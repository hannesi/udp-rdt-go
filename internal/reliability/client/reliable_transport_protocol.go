package client

import (
	"fmt"
	"log"

	"github.com/hannesi/udp-rdt-go/internal/reliability"
	"github.com/hannesi/udp-rdt-go/internal/virtualsocket"
	"github.com/hannesi/udp-rdt-go/pkg/utils"
)

type ReliableTransportProtocol struct {
	socket    *virtualsocket.VirtualSocket
	sequencer *utils.Sequencer
}

func NewReliableTransportProtocol(socket *virtualsocket.VirtualSocket) *ReliableTransportProtocol {
	return &ReliableTransportProtocol{
		socket:    socket,
		sequencer: utils.NewSequencer(255),
	}
}

func (r *ReliableTransportProtocol) Send(data []byte) error {
	packet := reliability.NewReliableDataTransferPacket(r.sequencer.Next(), data)

	serializedPacket, err := packet.Serialize()
	if err != nil {
		return err
	}

	log.Printf("Sending: \033[32m%s\033[0m", packet.Payload)
	r.sendWithRetransmission(serializedPacket)
	return nil
}

func (r *ReliableTransportProtocol) sendWithRetransmission(data []byte) {
	r.socket.Send(data)
	sequence, _, err := r.receiveAck()
	if sequence != r.sequencer.Current() || err != nil {
		log.Println("Resending packet.")
		r.sendWithRetransmission(data)
	}
}

func (r *ReliableTransportProtocol) receiveAck() (uint8, string, error) {
	buffer := make([]byte, 4)

    log.Println("Waiting for ack...")
	_, err := r.socket.Receive(buffer)

	if err != nil {
		err = fmt.Errorf("Error receiving ACK: %s", err.Error())
		return 0, "", err
	}

	seq, ack, err := reliability.DeserializeAckData(buffer)
	if err != nil {
		err = fmt.Errorf("Error deserializing ACK: %s", err.Error())
		return 0, "", err
	}

	log.Printf("Received ack: %s %d", ack, seq)

    return seq, ack, nil
}
