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
)

var rootCmd = &cobra.Command{
	Use:   "rediss-cli",
	Short: "Redis CLI interacts with a Redis server",
	Long: `A custom, simplified CLI to interact with Rediss server that takes user commands,
converts them to the Redis Serialization Protocol (RESP), and forwards them to the Rediss server.`,
	Run: func(cmd *cobra.Command, args []string) {
		interactWithRedisServer()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func interactWithRedisServer() {
	// Create a scanner to read input line by line, which is more reliable across different systems.

	conn, err := net.Dial("tcp", ":6379")
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

		// Scanner will read the next token (line in this case) when 'Scan' is called
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading from input: %s\n", err)
				os.Exit(1) // or handle the error as appropriate in your application.
			}
			// If it reaches here, it might be that the scanner reached the end of the input (EOF),
			// so you might want to handle that case, for example, exit the loop or the program.
			break
		}

		// Obtain the text that has been scanned
		input := scanner.Text()
    if input == "exit" {
      break
    }

		// Convert the command to RESP format
		respCommand := convertToRESP(strings.Fields(input))

		// Send the command to the server
		if _, err := conn.Write([]byte(respCommand)); err != nil {
			fmt.Println("Failed to send to Redis:", err)
			continue // Or handle the error as you prefer
		}

		respResponse, err := convertFromRESP(responseReader)
		if err != nil {
			fmt.Println("Failed to convert response:", err)
			continue
		}

		// Print the response
		fmt.Println(respResponse)
	}
}

func convertToRESP(commandPieces []string) string {
	// Use a string builder for efficient string concatenation
	var resp strings.Builder

	// Add the array prefix and size
	fmt.Fprintf(&resp, "*%d\r\n", len(commandPieces))

	for _, piece := range commandPieces {
		// Add the bulk string header and the command piece itself
		fmt.Fprintf(&resp, "$%d\r\n%s\r\n", len(piece), piece)
	}

	// For debugging: print the conversion. You might want to remove this in production code.
	fmt.Printf("RESP Conversion: From: %s, To: %s\n", strings.Join(commandPieces, " "), strings.Join(strings.Split(resp.String(), "\r\n"), "\\r\\n"))

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
	case '*': // Arrays
		return "", errors.New("array parsing not implemented")
	default:
		return "", errors.New("unknown data type")
	}
}
