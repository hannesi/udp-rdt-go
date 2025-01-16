package client

type ClientReliabilityLayer interface {
    Send(packet []byte) error
}
