package client

import (
	"log"

	"github.com/hannesi/udp-rdt-go/internal/reliability"
	"github.com/hannesi/udp-rdt-go/internal/virtualsocket"
	"github.com/hannesi/udp-rdt-go/pkg/utils"
)

type ReliabilityLayerWithPositiveAcks struct {
	Socket    *virtualsocket.VirtualSocket
	sequencer utils.Sequencer
}

func NewReliabilityLayerWithPositiveAcks(socket *virtualsocket.VirtualSocket) *ReliabilityLayerWithPositiveAcks {
    return &ReliabilityLayerWithPositiveAcks{
        Socket: socket,
        sequencer: *utils.NewSequencer(1),
    }
}

func (r *ReliabilityLayerWithPositiveAcks) Send(data []byte) error {
	packet := reliability.NewReliableDataTransferPacket(r.sequencer.Next(), data)

	serializedPacket, err := packet.Serialize()
	if err != nil {
		return err
	}

	r.sendWithRetransmission(serializedPacket)
	return nil
}

func (r *ReliabilityLayerWithPositiveAcks) sendWithRetransmission(data []byte) {
	r.Socket.Send(data)
	sequence, ack, err := r.receiveAck()
	if ack != "ACK" || sequence != r.sequencer.Current() || err != nil {
		log.Println("Resending packet.")
		r.sendWithRetransmission(data)
	}
}

func (r *ReliabilityLayerWithPositiveAcks) receiveAck() (uint8, string, error) {
	responseBuffer := make([]byte, 4)
	_, err := r.Socket.Receive(responseBuffer)

    sequence, ack, err := reliability.DeserializeAckData([]byte(responseBuffer))

	if err != nil {
		return 0, "", err
	}

	log.Printf("Received response %s %d.\n", ack, sequence)

	return sequence, ack, nil
}
