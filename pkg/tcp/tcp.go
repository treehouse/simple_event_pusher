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

// this is dead code now, keeping just in case

// tcp.OpenPort(":8080", func(tcpConn net.Listener) {
// 	listen for new connections
// 	http.Serve(
// 		tcpConn,
// 		nil, /* Handler (DefaultServeMux if nil) */
// 	)
// })
