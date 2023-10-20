package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/piratey7007/rediss/lib/rerror"
	"github.com/piratey7007/rediss/lib/resp"
	"github.com/piratey7007/rediss/server/commands"
	"github.com/piratey7007/rediss/server/messages"
  "github.com/piratey7007/rediss/server/db"
)

var options struct {
  AppendOnly bool
}

func main() {
	fmt.Println(messages.RpStartup(messages.Options{
		Version:    "7.2.1",
		Bits:       "64",
		Commit:     "00000000",
		Modified:   "0",
		Pid:        "93873",
		Port:       "6379",
		ConfigFile: false,
	}))

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(fn, line, err)
		return
	}

  // Read into memory
  if options.AppendOnly == true {
    db.DB, err = db.NewAof(db.AofOptions{
      Path: "database.aof",
    })
    db.DB.Load()
  } else {
    db.DB, err = db.NewRdb(db.RdbOptions{
      Path: "database.rdb",
      Interval: 10,
    })
    db.DB.Load()
  }

  // Execute AOF file
	aof.Read(func(value Value) {
		commandName := value.Array[0].Bulk
		args := value.Array[1:]

		cmd, exists := commands.Registry.Commands[commandName]
		if !exists {
			fmt.Printf("Invalid command AOF read: '%s'\n", commandName)
			return
		}

		cmd.Execute(args)
	})

  // Go routine to accept connections
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

  // Handle interrupt
  interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

  // Main thread
  for {
    select {
      case conn := <-connChan:
        go handleConnection(conn, aof)
      case <-interrupt:
        fmt.Println("Shutting down...")
        return
    }
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

		args := []resp.Value{}

		if i < len(value.Array) {
			args = value.Array[i:]
		}

		writer := resp.NewWriter(conn)

		cmd, exists := commands.Registry.Commands[command]
		if !exists {
			continue
		}

    if aof, ok := db.DB.(*db.Aof); ok {
      if command == "set" || command == "hset" {
        aof.Write(value)
      }
    }

		result := cmd.Execute(args)

		writer.Write(result)
	}
}
