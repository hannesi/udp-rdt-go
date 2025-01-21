package client

import (
	"log"

	"github.com/hannesi/udp-rdt-go/internal/reliability"
	"github.com/hannesi/udp-rdt-go/internal/virtualsocket"
)

type ReliabilityLayerWithPositiveAndNegativeAcks struct {
	Socket *virtualsocket.VirtualSocket
}

func (r *ReliabilityLayerWithPositiveAndNegativeAcks) Send(data []byte) error {
	packet := reliability.NewReliableDataTransferPacket(0, data)

	serializedPacket, err := packet.Serialize()
	if err != nil {
        return err
	}

    r.sendWithRetransmission(serializedPacket)
    return nil
}

func (r *ReliabilityLayerWithPositiveAndNegativeAcks) sendWithRetransmission(data []byte) {
    r.Socket.Send(data)
    response, err := r.receiveAck()
    if response == "NACK" || err != nil {
        log.Println("Resending packet.")
        r.sendWithRetransmission(data)
    }
}

func (r *ReliabilityLayerWithPositiveAndNegativeAcks) receiveAck() (string, error) {
	responseBuffer := make([]byte, 4)
	n, err := r.Socket.Receive(responseBuffer)
    responseString := string(responseBuffer[:n])

    log.Printf("Received %s response.\n", responseString)

	if err != nil {
		return responseString, err
	}
    
    return responseString, nil
}
