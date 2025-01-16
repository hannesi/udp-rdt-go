package client

import (
	"github.com/hannesi/udp-rdt-go/internal/reliability"
	"github.com/hannesi/udp-rdt-go/internal/virtualsocket"
)

type ReliabilityLayerWithBitErrorDetection struct {
	Socket *virtualsocket.VirtualSocket
}

func (r *ReliabilityLayerWithBitErrorDetection) Send(data []byte) error {
	packet := reliability.NewReliableDataTransferPacket(data)

	serializedPacket, err := packet.Serialize()
	if err != nil {
		return err
	}

	return r.Socket.Send(serializedPacket)
}
