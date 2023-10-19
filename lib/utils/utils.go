package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GetCommandJSON loads the JSON file data into the passed interface for the provided command name. Name must be in kebab-case-format.
func GetCommandJSON(commandName string, v interface{}) error {
  path := filepath.Join("data", "commands", fmt.Sprintf("%s.json", commandName))

  // Check if the command JSON file exists.
  if _, err := os.Stat(path); os.IsNotExist(err) {
    return fmt.Errorf("Command data for '%s' does not exist", commandName)
  }

  // Open the command JSON file.
  file, err := os.Open(path)
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
