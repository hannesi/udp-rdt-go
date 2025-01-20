package reliability

import (
	"bytes"
	"io"
)

// Serializes ack data to the following format: [sequence, ...msg]
func SerializeAckData(sequence uint8, msg string) ([]byte, error) {
    buffer := new(bytes.Buffer)
    
    err := buffer.WriteByte(sequence)
    if err != nil {
        return nil, err
    }

    _, err = buffer.Write([]byte(msg))
    if err != nil {
        return nil, err
    }

    return buffer.Bytes(), nil
}

// Deserializes ack data from the following format: [sequence, ...msg]
func DeserializeAckData(data []byte) (uint8, string, error) {
    buffer := bytes.NewReader(data)

    sequence, err := buffer.ReadByte()
    if err != nil {
        return 0, "", err
    }

    msg, err := io.ReadAll(buffer)
    if err != nil {
        return 0, "", err
    }

    msgString := string(msg)

    return sequence, msgString, nil
}
