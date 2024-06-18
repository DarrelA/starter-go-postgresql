package utils

import "net"

func GetIPVersion(ip string) string {
	parsedIP := net.ParseIP(ip)

	if parsedIP == nil {
		return "unknown"
	}
	if parsedIP.To4() != nil {
		return "IPv4"
	}

	return "IPv6"
}
