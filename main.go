package main

import (
	"fmt"
//	"io"
	"net"
//	"os"
)

func main() {
	// Listen
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

  // Accept
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

  // Read
	for {
    resp := NewResp(conn)

    value, err := resp.Read()
		if err != nil {
      fmt.Println(err)
			return
		}

    fmt.Println("value:", value)

		conn.Write([]byte("+OK\r\n"))
	}
}
