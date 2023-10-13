package main

import (
	"fmt"
	"io"
	"net"
	"os"
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
		buf := make([]byte, 1024)

		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		conn.Write([]byte("+OK\r\n"))
	}

}
