//go:build windows
// +build windows

package core

import (
	"fmt"
	"net"
	"strings"
	"syscall"
)

func init() {
	SetSocketOptions = func(network, address string, c syscall.RawConn, interfaceName string) (err error) {
		switch network {
		case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
		default:
			return fmt.Errorf("unsupported network type: %s", network)
		}

		var innerErr error
		err = c.Control(func(fd uintptr) {
			host, _, errSplit := net.SplitHostPort(address)
			if errSplit == nil {
				if ip := net.ParseIP(host); ip != nil && !ip.IsGlobalUnicast() {
					return
				}
			}

			isIPv6 := false
			if strings.HasSuffix(network, "6") {
				isIPv6 = true
			} else if strings.HasSuffix(network, "4") {
				isIPv6 = false
			} else if errSplit == nil {
				ip := net.ParseIP(host)
				if ip != nil && ip.To4() == nil {
					isIPv6 = true
				}
			}

			innerErr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			if innerErr != nil {
				return
			}

			if interfaceName != "" {
				iface, err := net.InterfaceByName(interfaceName)
				if err != nil {
					innerErr = fmt.Errorf("interface %s not found: %v", interfaceName, err)
					return
				}

				addrs, err := iface.Addrs()
				if err != nil {
					innerErr = fmt.Errorf("failed to get interface addresses: %v", err)
					return
				}

				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok {
						var sockaddr syscall.Sockaddr
						if !isIPv6 && ipnet.IP.To4() != nil {
							sa := &syscall.SockaddrInet4{}
							copy(sa.Addr[:], ipnet.IP.To4())
							sockaddr = sa
						} else if isIPv6 && ipnet.IP.To16() != nil && ipnet.IP.To4() == nil {
							sa := &syscall.SockaddrInet6{}
							copy(sa.Addr[:], ipnet.IP.To16())
							sockaddr = sa
						}

						if sockaddr != nil {
							innerErr = syscall.Bind(syscall.Handle(fd), sockaddr)
							if innerErr == nil {
								return
							}
						}
					}
				}

				if innerErr == nil {
					innerErr = fmt.Errorf("no suitable address found for interface %s", interfaceName)
				}
			}
		})

		if innerErr != nil {
			err = innerErr
		}
		return
	}
}
