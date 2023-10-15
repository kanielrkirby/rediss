package commands

import (
  "os"
  "encoding/json"
  "github.com/piratey7007/rediss/resp"
  "log"
  "fmt"
  "path/filepath"
  "strings"
)

type CommandMetadata struct {
	Name       string
	Summary    string
	Complexity string
}

type Command struct {
  Execute func(args []resp.Value) resp.Value
  CommandMetadata
}


type registry struct {
  Commands map[string]Command
}

var Registry = &registry{
  Commands: make(map[string]Command),
}

func (r *registry) Register(name string, cmd func (args []resp.Value) resp.Value) error {
  kebabName := strings.ReplaceAll(name, " ", "-")
  path := filepath.Join("commands", "json", fmt.Sprintf("%s.json", kebabName))
  metadata, err := ReadJSON(path)
  if err != nil {
    log.Fatal(err)
  }

  r.Commands[name] = Command{
    Execute: cmd,
    CommandMetadata: metadata,
  }

  return nil
}

func ReadJSON(path string) (CommandMetadata, error) {
  file, err := os.Open(path)
  if err != nil {
    return CommandMetadata{}, err
  }
  defer file.Close()

  var metadata CommandMetadata
  err = json.NewDecoder(file).Decode(&metadata)
  if err != nil {
    return CommandMetadata{}, err
  }

  return metadata, nil
}

