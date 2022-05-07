package gognet

import (
	"net"
)

const (
	Protocol_Version byte = 0x01
)

var timeout_time int64 = 900

func ConvertDomain(domain string) net.IP {
	ips, _ := net.LookupIP(domain)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4
		}
	}
	return nil
}
