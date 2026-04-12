package wol

import (
	"fmt"
	"net"
	"time"
)

// IsHostUp checks if a host is reachable by attempting a TCP connection
// on common ports (22, 80, 443, 3389).
func IsHostUp(ip string) bool {
	ports := []int{22, 80, 443, 3389}
	for _, port := range ports {
		addr := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
		if err == nil {
			conn.Close()
			return true
		}
	}
	return false
}

// WaitForHost polls the host until it comes online or the timeout (in seconds) expires.
func WaitForHost(ip string, timeoutSec int) bool {
	deadline := time.Now().Add(time.Duration(timeoutSec) * time.Second)
	for time.Now().Before(deadline) {
		if IsHostUp(ip) {
			return true
		}
		time.Sleep(2 * time.Second)
	}
	return false
}
