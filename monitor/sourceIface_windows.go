//go:build windows
// +build windows

package main

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

const SIO_RCVALL = 0x98000001 // 控制代码，用于接收所有数据包

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
				innerErr = syscall.Bind(syscall.Handle(fd), sockaddr)
				if innerErr != nil {
					return
				}
			}

			// 设置 SO_REUSEADDR 选项
			innerErr = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			if innerErr != nil {
				return
			}

			// 绑定到特定接口
			if interfaceName != "" {
				var ifreq [256]byte
				copy(ifreq[:], interfaceName)
				var bytesReturned uint32
				innerErr = syscall.WSAIoctl(syscall.Handle(fd), SIO_RCVALL, (*byte)(unsafe.Pointer(&ifreq[0])), uint32(len(interfaceName)), nil, 0, &bytesReturned, nil, 0)
				if innerErr != nil {
					return
				}
			}
		})

		if innerErr != nil {
			err = innerErr
		}
		return
	}
}
