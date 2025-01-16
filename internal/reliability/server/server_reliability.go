package server

type ServerReliabilityLayer interface {
    Receive(packet []byte) 
}
