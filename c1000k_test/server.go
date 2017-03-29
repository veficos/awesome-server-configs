package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	var array []byte = make([]byte, 512)

	for i := 0; i < 512; i += 1 {
		array[i] = 'A'
	}

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", ":8888")
	listenFd, _ := net.ListenTCP("tcp", tcpAddr)

	for {
		conn, err := listenFd.Accept()
		if err != nil {
			fmt.Printf("accept err: %s\n", err)
			continue
		}
		go func(conn net.Conn) {
			defer func() { conn.Close() }()
			for {
				conn.SetReadDeadline(time.Now().Add(180 * time.Second))
				_, err := conn.Write(array)
				if err == nil {
					time.Sleep(1800 * time.Second)
					return
				}
				fmt.Printf("write err: %s\n", err)
				time.Sleep(60 * time.Second)
			}
		}(conn)
	}
}
