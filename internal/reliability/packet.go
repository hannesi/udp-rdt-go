package reliability

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

// Represents a packet used for detecting bit errors.
type ReliableDataTransferPacket struct {
    Payload []byte  // The payload of the packet.
    Checksum uint32 // The checksum calculated for the payload.
}

// Create a new rdt packet from provided payload. The checksum is calculated automagically.
func NewReliableDataTransferPacket(payload []byte) *ReliableDataTransferPacket {
    packet := &ReliableDataTransferPacket{
        Payload: payload,
    }
    packet.Checksum = packet.computeChecksum()
    return packet
}

func (p *ReliableDataTransferPacket) computeChecksum() uint32 {
    return crc32.ChecksumIEEE(p.Payload)
}

// Returns true if the rdt packet's checksum matches the packet's payload's calculated checksum.
func (p *ReliableDataTransferPacket) IsChecksumValid() bool {
    return p.Checksum == p.computeChecksum()
}

// Serialize a rdt packet into transferable form.
func (p *ReliableDataTransferPacket) Serialize() ([]byte, error) {
    buffer := new(bytes.Buffer)

    err := binary.Write(buffer, binary.BigEndian, p.Checksum)
    if err != nil {
        return nil, err
    }

    _, err = buffer.Write(p.Payload)
    if err != nil {
        return nil, err
    }

    return buffer.Bytes(), nil
}


// Deserialize a byte array into a rdt packet.
func Deserialize(data []byte) (*ReliableDataTransferPacket, error) {
    buffer := bytes.NewReader(data)

    var checksum uint32

    err := binary.Read(buffer, binary.BigEndian, &checksum)
    if err != nil {
        return nil, err
    }

    payload := make([]byte, buffer.Len())
    _, err = buffer.Read(payload)
    if err != nil {
        return nil, err
    }

    return &ReliableDataTransferPacket{
        Payload: payload,
        Checksum: checksum,
    }, nil
}
