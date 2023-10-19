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
	Port     int
	Password string
	Db       int
}

func ConnectToServer(options ConnectionOptions) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", options.Host, options.Port))
	if err != nil {
		fmt.Println("Failed to connect to Redis", err)
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	RESP := resp.NewResp(conn)
	fmt.Println("Connected to Redis server. You may start typing commands.")

	for {
		fmt.Print("redis-cli> ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading from input: %s\n", err)
				os.Exit(1)
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

    fmt.Println("'1'")
		respResponse, err := RESP.Read()
		if err != nil {
			fmt.Println("Failed to convert response:", err)
			continue
		}

    fmt.Println("'2'")
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
    fmt.Println("'3'")
	}
}
