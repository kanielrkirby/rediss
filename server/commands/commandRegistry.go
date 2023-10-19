package commands

import (
  _ "embed"
	"encoding/json"
	"github.com/piratey7007/rediss/lib/resp"
	"os"
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

var jsonContent string

func (r *registry) Register(name string, cmd func(args []resp.Value) resp.Value) error {
//	path := filepath.Join(
//		"commands",
//		"json",
//		fmt.Sprintf("%s.json",
//			strings.ReplaceAll(name, " ", "-")))
//   metadata, err := ReadJSON(path)
  var metadata CommandMetadata

  content := strings.Replace(jsonContent, name, strings.ReplaceAll(name, " ", "-"), -1)
  err := json.Unmarshal([]byte(content), &metadata)
  if err != nil {
    return err
  }

	r.Commands[name] = Command{
		Execute:         cmd,
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
