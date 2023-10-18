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
  Host string
  Port int
  Password string
  Db int
}

func ConnectToServer(options ConnectionOptions) {
  conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", options.Host, options.Port))
  if err != nil {
    fmt.Println("Failed to connect to Redis", err)
    return
  }
  defer conn.Close()

  scanner := bufio.NewScanner(os.Stdin)
  responseReader := bufio.NewReader(conn)
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

    respCommand := resp.ConvertToRESP(strings.Fields(input))
    
    if _, err := conn.Write([]byte(respCommand)); err != nil {
      fmt.Println("Failed to send to Redis:", err)
      continue 
    }

    respResponse, err := resp.ConvertFromRESP(responseReader)
    if err != nil {
      fmt.Println("Failed to convert response:", err)
      continue
    }

    fmt.Println(respResponse)
  }
}

