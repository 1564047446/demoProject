package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func sender(conn net.Conn) {
	words := "hello world!"
	conn.Write([]byte(words))
	fmt.Println("send over")
	time.Sleep(1 * time.Second)

}

func main() {
	server := "192.168.132.211:1236"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	sender(conn)

}
