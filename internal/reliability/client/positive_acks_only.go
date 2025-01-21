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
		Socket:    socket,
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
	if ack == "ACK" && sequence == r.sequencer.Current() && err == nil {
		log.Printf("Received valid ACK: %s %d.\n", ack, sequence)
		return
	}

	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Invalid ack or sequence mismatch.")
	}
	log.Println("Resending packet.")
    r.sendWithRetransmission(data)
}

func (r *ReliabilityLayerWithPositiveAcks) receiveAck() (uint8, string, error) {
	buffer := make([]byte, 4)

    log.Println("Waiting for ack...")
	_, err := r.Socket.Receive(buffer)

	if err != nil {
		log.Println("Error receiving ACK:", err.Error())
		return 0, "", err
	}

	sequence, ack, err := reliability.DeserializeAckData(buffer)
	if err != nil {
		log.Println("Error deserializing ACK:", err.Error())
		return 0, "", err
	}

    return sequence, ack, nil
}
