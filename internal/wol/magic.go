package wol

import (
	"fmt"
	"net"
)

// SendMagicPacket sends a Wake-on-LAN magic packet to the given MAC address.
// The magic packet is 6 bytes of 0xFF followed by the MAC address repeated 16 times.
// An optional SecureOn password (6 bytes) is appended if provided.
func SendMagicPacket(mac net.HardwareAddr, broadcast string, port int, secureOn []byte) error {
	if len(mac) != 6 {
		return fmt.Errorf("MAC address must be 6 bytes, got %d", len(mac))
	}

	// Build magic packet: 6x 0xFF + 16x MAC
	packetLen := 6 + 16*6
	if len(secureOn) > 0 {
		packetLen += len(secureOn)
	}
	packet := make([]byte, 0, packetLen)

	// Header: 6 bytes of 0xFF
	for i := 0; i < 6; i++ {
		packet = append(packet, 0xFF)
	}

	// Body: MAC address repeated 16 times
	for i := 0; i < 16; i++ {
		packet = append(packet, mac...)
	}

	// Optional SecureOn password
	if len(secureOn) > 0 {
		packet = append(packet, secureOn...)
	}

	addr := fmt.Sprintf("%s:%d", broadcast, port)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return fmt.Errorf("connecting to %s: %w", addr, err)
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("sending packet: %w", err)
	}

	return nil
}
