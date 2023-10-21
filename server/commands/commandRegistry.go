package commands

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"
  "sync"

	"github.com/piratey7007/rediss/lib/resp"
  "github.com/piratey7007/rediss/lib/utils"
)

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}
var SUBs = map[string][]chan resp.Value{}
var SUBsMu = sync.RWMutex{}

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
	Commands: map[string]Command{},
}

var jsonContent string

func (r *registry) Register(name string, cmd func(args []resp.Value) resp.Value) error {
  var metadata CommandMetadata

  err := utils.GetCommandJSON(strings.ReplaceAll(name, " ", "-"), &metadata)
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
