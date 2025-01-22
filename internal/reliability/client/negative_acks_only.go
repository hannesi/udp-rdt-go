package client

import (
	"fmt"
	"log"

	"github.com/hannesi/udp-rdt-go/internal/reliability"
	"github.com/hannesi/udp-rdt-go/internal/virtualsocket"
	"github.com/hannesi/udp-rdt-go/pkg/utils"
)

type ReliabilityLayerWithNegativeAcks struct {
	Socket    *virtualsocket.VirtualSocket
	sequencer utils.Sequencer
}

func NewReliabilityLayerWithNegativeAcks(socket *virtualsocket.VirtualSocket) *ReliabilityLayerWithNegativeAcks {
	return &ReliabilityLayerWithNegativeAcks{
		Socket:    socket,
		sequencer: *utils.NewSequencer(1),
	}
}

func (r *ReliabilityLayerWithNegativeAcks) Send(data []byte) error {
	packet := reliability.NewReliableDataTransferPacket(r.sequencer.Next(), data)

	serializedPacket, err := packet.Serialize()
	if err != nil {
		return err
	}

	r.sendWithRetransmission(serializedPacket)
	return nil
}

func (r *ReliabilityLayerWithNegativeAcks) sendWithRetransmission(data []byte) {
	r.Socket.Send(data)
	sequence, ack, _ := r.receiveAck()
	if ack == "NAK" {
		log.Printf("Received valid ack packet: %s %d.\n", ack, sequence)
        log.Println("Resending packet.")
		r.sendWithRetransmission(data)
	}
}

func (r *ReliabilityLayerWithNegativeAcks) receiveAck() (uint8, string, error) {
	buffer := make([]byte, 4)

    log.Println("Waiting for ack...")
	_, err := r.Socket.Receive(buffer)

	if err != nil {
		err = fmt.Errorf("Error receiving ACK: %s", err.Error())
		return 0, "", err
	}

	sequence, ack, err := reliability.DeserializeAckData(buffer)
	if err != nil {
		err = fmt.Errorf("Error deserializing ACK: %s", err.Error())
		return 0, "", err
	}

    return sequence, ack, nil
}
