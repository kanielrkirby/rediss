package connections

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/piratey7007/rediss/lib/resp"
)

type ConnectionOptions struct {
	Host     string
	Port     string
  Command   string
}

func ConnectToServer(options ConnectionOptions) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", options.Host, options.Port))
	if err != nil {
		fmt.Println("Failed to connect to Redis", err)
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	RESP := resp.NewResp(conn)

  if options.Command != "" {
    value := resp.Value{
      Typ:   "array",
      Array: []resp.Value{},
    }

    elems := strings.Split(options.Command, " ")

    for _, elem := range elems {
      value.Array = append(value.Array, resp.Value{
        Typ:  "bulk",
        Bulk: elem,
      })
    }

    bytes := value.Marshal()

    if _, err := conn.Write(bytes); err != nil {
      fmt.Println("Failed to send to Redis:", err)
    }

    respResponse, err := RESP.Read()
    if err != nil {
      fmt.Println("Failed to convert response:", err)
    }

    switch respResponse.Typ {
    case "string":
      fmt.Println(respResponse.Str)
    case "error":
      fmt.Println("Error:", respResponse.Str)
    case "bulk":
      fmt.Println(respResponse.Bulk)
    case "int":
      fmt.Println(respResponse.Num)
    case "array":
      for _, respResponse := range respResponse.Array {
        fmt.Println(respResponse.Bulk)
      }
    default:
      fmt.Println("Unknown response type:", respResponse.Typ)
    }
    return
  }

  fmt.Println("Connected to Redis server. You may start typing commands.")
	for {
		fmt.Print("redis-cli> ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
        fmt.Println("Error: Server closed the connection.")
        fmt.Println("not connected>")
			}
			break
		}

		input := scanner.Text()
		if input == "exit" {
			break
		}

		elems := strings.Split(input, " ")

		value := resp.Value{
			Typ:   "array",
			Array: []resp.Value{},
		}

		for _, elem := range elems {
			value.Array = append(value.Array, resp.Value{
				Typ:  "bulk",
				Bulk: elem,
			})
		}

		bytes := value.Marshal()

		if _, err := conn.Write(bytes); err != nil {
			fmt.Println("Failed to send to Redis:", err)
			continue
		}

		respResponse, err := RESP.Read()
		if err != nil {
			fmt.Println("Failed to convert response:", err)
			continue
		}

    switch respResponse.Typ {
    case "string":
      fmt.Println(respResponse.Str)
    case "error":
      fmt.Println("Error:", respResponse.Str)
    case "bulk":
      fmt.Println(respResponse.Bulk)
    case "int":
      fmt.Println(respResponse.Num)
    case "array":
      for _, respResponse := range respResponse.Array {
        fmt.Println(respResponse.Bulk)
      }
    default:
      fmt.Println("Unknown response type:", respResponse.Typ)
    }

    if strings.HasPrefix(input, "subscribe") {
      handleSubscribe(RESP)
    }
	}
}

func handleSubscribe(RESP *resp.Resp) {
  for {
    respResponse, err := RESP.Read()
    if err != nil {
      fmt.Println("Failed to convert response:", err)
      continue
    }

    switch respResponse.Typ {
    case "string":
      fmt.Println(respResponse.Str)
    case "error":
      fmt.Println("Error:", respResponse.Str)
    case "bulk":
      fmt.Println(respResponse.Bulk)
    case "int":
      fmt.Println(respResponse.Num)
    case "array":
      for _, respResponse := range respResponse.Array {
        fmt.Println(respResponse.Bulk)
      }
    default:
      fmt.Println("Unknown response type:", respResponse.Typ)
    }
  }
}
