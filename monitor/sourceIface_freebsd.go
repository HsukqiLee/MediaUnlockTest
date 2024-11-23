//go:build freebsd
// +build freebsd

package main

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"
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

			// 绑定到特定接口
			if interfaceName != "" {
				var ifreq [32]byte
				copy(ifreq[:], interfaceName)
				_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.SIOCGIFADDR, uintptr(unsafe.Pointer(&ifreq[0])))
				if errno != 0 {
					innerErr = errno
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
