package connections

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
	"github.com/piratey7007/rediss/server/db"
	"github.com/piratey7007/rediss/server/messages"
	"github.com/piratey7007/rediss/server/types"
)

var options struct {
  AppendOnly bool
}

func StartServer(options types.ConnectionOptions) {
	fmt.Println(messages.RpStartup(options))

  l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", options.Host, options.Port))
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		fmt.Println(fn, line, err)
		return
	}

  // Read aof into memory
  db.DB, err = db.NewAof(db.AofOptions{
    Path: "database.aof",
    Interval: 1,
  })

	db.DB.Read()

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
        go handleConnection(conn)
      case <-interrupt:
        fmt.Println("Shutting down...")
        return
    }
  }
}

// handleConnection handles a single connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

  RESP := resp.NewResp(conn)
  writer := resp.NewWriter(conn)
  commands.Registry.Writer = *writer

	for {
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
