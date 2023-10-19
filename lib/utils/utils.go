package utils

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/piratey7007/rediss/lib/data"
)

var commandsDir, commandsDirErr = fs.Sub(data.EmbeddedFiles, "commands")

// GetCommandJSON loads the JSON file data into the passed interface for the provided command name. Name must be in kebab-case-format.
func GetCommandJSON(commandName string, v interface{}) error {
  file, err := commandsDir.Open(fmt.Sprintf("%s.json", commandName))
  if err != nil {
    return err
  }
  defer file.Close()

  // Decode the command JSON file.
  decoder := json.NewDecoder(file)
  err = decoder.Decode(v)
  if err != nil {
    return err
  }

  return nil
}
