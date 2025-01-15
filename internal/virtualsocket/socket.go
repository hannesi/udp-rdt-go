package virtualsocket

import (
	"log"
	"math/rand/v2"
	"net"
	"time"

	"github.com/hannesi/udp-rdt-go/internal/config"
)

// VirtualSocket wraps an UDP connection and simulates an unreliable network, introducing delay and bit errors before passing it to the actual UDP socket, or dropping the packet.
type VirtualSocket struct {
	socket    *net.UDPConn
	delay     time.Duration
	delayRate float64
	dropRate  float64
}

// Creates a new virtual socket.
func NewVirtualSocket() (*VirtualSocket, error) {
	destAddr := net.UDPAddr{
		IP:   net.ParseIP(config.DefaultConfig.IPAddrString),
		Port: config.DefaultConfig.ServerPort,
	}

	socketAddr := net.UDPAddr{
		IP:   net.ParseIP(config.DefaultConfig.IPAddrString),
		Port: 0,
	}

	socket, err := net.DialUDP("udp", &socketAddr, &destAddr)

	if err != nil {
		return nil, err
	}

	log.Println("Virtual socket initialized.")

	return &VirtualSocket{
		socket:    socket,
		delay:     config.DefaultConfig.VirtualSocketDelay,
		delayRate: config.DefaultConfig.VirtualSocketDelayRate,
		dropRate:  config.DefaultConfig.VirtualSocketDropRate,
	}, nil
}

// Send data using the virtual socket.
func (vs *VirtualSocket) Send(data []byte) error {
	if vs.shouldDrop() {
		return nil
	}

	vs.handlePacketDelay()

	_, err := vs.socket.Write(data)
	return err
}

// Close the socket wrapped inside the virtual socket.
func (vs *VirtualSocket) Close() {
	if vs.socket != nil {
		vs.socket.Close()
	}
}

func (vs *VirtualSocket) shouldDrop() bool {
	packetDropped := rand.Float64() < vs.dropRate
	if packetDropped {
		log.Println("Packet dropped.")
	}
	return packetDropped
}

func (vs *VirtualSocket) handlePacketDelay() {
	if rand.Float64() < vs.delayRate {
		log.Printf("Packet delayed by %d milliseconds.\n", vs.delay.Milliseconds())
		time.Sleep(vs.delay)
	}
}
