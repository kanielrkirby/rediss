package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
  "github.com/spf13/viper"
)

type ConnectionOptions struct {
  host string
  port int
  password string
  db int
}

var rootCmd = &cobra.Command{
	Use:   "rediss-cli",
	Short: "Redis CLI interacts with a Redis server",
	Long: `A custom, simplified CLI to interact with Rediss server that takes user commands,
converts them to the Redis Serialization Protocol (RESP), and forwards them to the Rediss server.`,
	Run: func(cmd *cobra.Command, args []string) {
    options := ConnectionOptions{
      host: viper.GetString("bind"),
      port: viper.GetInt("port"),
      password: viper.GetString("password"),
      db: viper.GetInt("db"),
    }
    connectToServer(options)
	},
  
}

func init() {
  helpFlag := false
  rootCmd.PersistentFlags().BoolVarP(&helpFlag, "help", "", false, "Help default flag")
  cobra.OnInitialize(initConfig)
  rootCmd.PersistentFlags().StringP("host", "h", "localhost", "The host to bind to")
  rootCmd.PersistentFlags().IntP("port", "p", 6379, "The port to bind to")
  rootCmd.PersistentFlags().StringP("password", "a", "", "The password to use when connecting to the server")
  rootCmd.PersistentFlags().IntP("db", "n", 0, "The database to use when connecting to the server")
  viper.BindPFlag("bind", rootCmd.PersistentFlags().Lookup("host"))
  viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
  viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
  viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
}

func initConfig() {
  viper.SetConfigName("config")
  viper.AddConfigPath(".")
  viper.SetConfigType("conf")

  viper.SetDefault("bind", "localhost")
  viper.SetDefault("port", 6379)
  viper.SetDefault("password", "")
  viper.SetDefault("db", 0)

  if err := viper.ReadInConfig(); err != nil {
    fmt.Println("Error reading config file:", err)
  }
}

func connectToServer(options ConnectionOptions) {
  conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", options.host, options.port))
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

    respCommand := convertToRESP(strings.Fields(input))
    
    if _, err := conn.Write([]byte(respCommand)); err != nil {
      fmt.Println("Failed to send to Redis:", err)
      continue 
    }

    respResponse, err := convertFromRESP(responseReader)
    if err != nil {
      fmt.Println("Failed to convert response:", err)
      continue
    }

    fmt.Println(respResponse)
  }
}

func Run() {
	cobra.CheckErr(rootCmd.Execute())
}

func convertToRESP(commandPieces []string) string {
	
	var resp strings.Builder

	fmt.Fprintf(&resp, "*%d\r\n", len(commandPieces))

	for _, piece := range commandPieces {
		fmt.Fprintf(&resp, "$%d\r\n%s\r\n", len(piece), piece)
	}

	return resp.String()
}

func convertFromRESP(reader *bufio.Reader) (string, error) {
	_type, err := reader.ReadByte()
	if err != nil {
		return "", err
	}

	switch _type {
	case '+', '-', ':':
		line, _, err := reader.ReadLine()
		if err != nil {
			return "", err
		}
		return string(line), nil
	case '$':
		lengthStr, _, err := reader.ReadLine()
		if err != nil {
			return "", err
		}

		length, err := strconv.Atoi(string(lengthStr))
		if err != nil {
			return "", err
		}

		if length == -1 {
			return "(nil)", nil
		}

		data := make([]byte, length)
		_, err = reader.Read(data)
		if err != nil {
			fmt.Println("Read error: ", err)
			return "", err
		}
		fmt.Println("Read success")

		_, err = reader.Discard(2)
		if err != nil {
			fmt.Println("Discard error: ", err)
			return "", err
		}
		fmt.Println("Discard success")

		return string(data), nil
	case '*': 
		return "", errors.New("array parsing not implemented")
	default:
		return "", errors.New("unknown data type")
	}
}
