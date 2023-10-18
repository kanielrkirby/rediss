package main

import (
	"fmt"
	"net"
	"runtime"
	"strings"

	"github.com/piratey7007/rediss/server/commands"
	"github.com/piratey7007/rediss/lib/rerror"
	"github.com/piratey7007/rediss/lib/resp"
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
		fmt.Println(rerror.ErrWrap(err).FormatAndError("Error creating aof file"))
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

	connChan := make(chan net.Conn)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println(rerror.ErrWrap(err).FormatAndError("Error accepting connection"))
				continue
			}
			connChan <- conn
		}
	}()

	for conn := range connChan {
		go handleConnection(conn, aof)
	}

}

func handleConnection(conn net.Conn, aof *Aof) {
	defer conn.Close()
	for {
		RESP := resp.NewResp(conn)
		value, err := RESP.Read()
		if err != nil {
      if strings.HasSuffix(err.Error(), "EOF") {
        fmt.Println("Connection closed")
				return
			}
			errString := rerror.ErrInternal.FormatAndError("Error reading value")
			if rerror.DEBUG {
				custom := fmt.Sprintf("value.Typ: %s, value.Str: %s, value.Num: %d, value.Bulk: %s, value.Array: %v", value.Typ, value.Str, value.Num, value.Bulk, value.Array)
				errString = rerror.ErrInternal.FormatAndError("Error reading value: %s", custom)
			}
			fmt.Println(rerror.ErrWrap(err).FormatAndError(errString))
			return
		}

		if value.Typ != "array" {
			fmt.Println(rerror.ErrInternal.FormatAndError("Invalid request, expected array"))
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println(rerror.ErrWrap(err).FormatAndError("Invalid request, expected array length > 0"))
			continue
		}

		str := ""
		tempVal := ""
		i := 0
		for _, v := range value.Array {
			tempVal = str + " " + v.Bulk
			if _, exists := commands.Registry.Commands[strings.ToLower(strings.TrimSpace(tempVal))]; !exists {
				break
			}
			str = tempVal
			i++
		}

		command := strings.ToLower(strings.TrimSpace(str))
		if rerror.DEBUG {
			fmt.Println("command: ", command)
		}

		args := []resp.Value{}

		if i < len(value.Array) {
			args = value.Array[i:]
		}

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
