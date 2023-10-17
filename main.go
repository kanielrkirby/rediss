package main

import (
	"fmt"
	"net"
	"strings"
  "runtime"

	"github.com/piratey7007/rediss/commands"
	"github.com/piratey7007/rediss/resp"
)

func main() {
	fmt.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
    _, fn, line, _ := runtime.Caller(1)
		fmt.Println(fn, line, err)
		return
	}

	aof, err := NewAof("database.aof")
	if err != nil {
    _, fn, line, _ := runtime.Caller(1)
		fmt.Println(fn, line, err)
		return
	}
	defer aof.Close()

	aof.Read(func(value Value) {
		commandName := value.Array[0].Bulk
		args := value.Array[1:]

    cmd, exists := commands.Registry.Commands[commandName]
    if !exists {
      fmt.Println("Invalid command: ", commandName)
      return
    }

    cmd.Execute(args)
	})

	for {
    conn, err := l.Accept()
    if err != nil {
      _, fn, line, _ := runtime.Caller(1)
      fmt.Println(fn, line, "Error accepting connection: " + err.Error())
			continue
    }

    go handleConnection(conn, aof)
	}
}

func handleConnection(conn net.Conn, aof *Aof) {
  defer conn.Close()
	for {
		RESP := resp.NewResp(conn)
		value, err := RESP.Read()
		if err != nil {
      if err.Error() == "EOF" {
        fmt.Println("Connection closed")
        return
      }
      fmt.Println(err.Error())
			return
		}

		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

    command := strings.ToLower(value.Array[0].Bulk)

		args := value.Array[1:]

		writer := resp.NewWriter(conn)

    cmd, exists := commands.Registry.Commands[command]

    if !exists {
      fmt.Println("Invalid command: ", command)
      continue
    }

		if command == "set" || command == "hset" {
			aof.Write(value)
		}

		result := cmd.Execute(args)
		writer.Write(result)
	}
}
