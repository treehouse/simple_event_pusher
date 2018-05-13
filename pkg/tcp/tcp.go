package tcp

import (
	"net"
)

func OpenPort(incomingPort string, open func(net.Listener)) {
	tcp, err := net.Listen("tcp", incomingPort)
	if err != nil {
		return
	}
	defer tcp.Close()

	open(tcp)
}
