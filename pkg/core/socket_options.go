package core

import "syscall"

var SetSocketOptions = func(network, address string, c syscall.RawConn, interfaceName string) (err error) {
	return
}
