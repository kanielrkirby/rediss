package resp

import (
  "bufio"
  "errors"
  "fmt"
  "strconv"
  "strings"
)

func ConvertToRESP(commandPieces []string) string {
	
	var resp strings.Builder

	fmt.Fprintf(&resp, "*%d\r\n", len(commandPieces))

	for _, piece := range commandPieces {
		fmt.Fprintf(&resp, "$%d\r\n%s\r\n", len(piece), piece)
	}

	return resp.String()
}

func ConvertFromRESP(reader *bufio.Reader) (string, error) {
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
