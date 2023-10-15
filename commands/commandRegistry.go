package commands

import (
	"fmt"
)

type CommandMetadata struct {
	Name       string
	Summary    string
	Complexity string
}

type Command interface {
	Execute(args []string) (string, error)
}

type Registry struct {
  commands map[string]struct {
    Execute func(args []string) (string, error)
    CommandMetadata
  }
}

func (r *Registry) Register(name string, cmd Command, metadata CommandMetadata) error {
  if _, exists := r.commands[name]; exists {
    return fmt.Errorf("Command %s already registered", name)
  }
  r.commands[name] = struct {
    Execute func(args []string) (string, error)
    CommandMetadata
  }{
    Execute: cmd.Execute,
    CommandMetadata: metadata,
  }
  return nil
}
