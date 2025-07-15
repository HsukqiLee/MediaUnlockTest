//go:build netbsd
// +build netbsd

package main

import (
	"fmt"
	"net"
	"syscall"
)

func init() {
	setSocketOptions = func(network, address string, c syscall.RawConn, interfaceName string) (err error) {
		switch network {
		case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
			// 继续处理
		default:
			return fmt.Errorf("unsupported network type: %s", network)
		}

		var innerErr error
		err = c.Control(func(fd uintptr) {
			// 绑定到指定的 IP 地址
			if address != "" {
				sockaddr := &syscall.SockaddrInet4{}
				copy(sockaddr.Addr[:], net.ParseIP(address).To4())
				innerErr = syscall.Bind(int(fd), sockaddr)
				if innerErr != nil {
					return
				}
			}

			// 设置 SO_REUSEADDR 选项
			innerErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			if innerErr != nil {
				return
			}

			// 绑定到特定接口 - 使用更安全的方法
			if interfaceName != "" {
				// 查找指定接口的信息
				iface, err := net.InterfaceByName(interfaceName)
				if err != nil {
					innerErr = fmt.Errorf("interface %s not found: %v", interfaceName, err)
					return
				}

				// 获取接口的IP地址
				addrs, err := iface.Addrs()
				if err != nil {
					innerErr = fmt.Errorf("failed to get interface addresses: %v", err)
					return
				}

				// 寻找第一个可用的IP地址进行绑定
				for _, addr := range addrs {
					if ipnet, ok := addr.(*net.IPNet); ok {
						var sockaddr syscall.Sockaddr
						if ip4 := ipnet.IP.To4(); ip4 != nil {
							// IPv4地址
							sa := &syscall.SockaddrInet4{}
							copy(sa.Addr[:], ip4)
							sockaddr = sa
						} else if ip6 := ipnet.IP.To16(); ip6 != nil {
							// IPv6地址
							sa := &syscall.SockaddrInet6{}
							copy(sa.Addr[:], ip6)
							sockaddr = sa
						}

						if sockaddr != nil {
							// 使用标准的bind方法
							innerErr = syscall.Bind(int(fd), sockaddr)
							if innerErr == nil {
								return // 成功绑定
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
