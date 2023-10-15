package commands

import (
	"github.com/piratey7007/rediss/resp"
  "strings"
)

func init() {
	Registry.Register("command", command)
}

// Command: COMMAND [INFO] [COMMAND_NAME]
func command(args []resp.Value) resp.Value {
	// this either returns a list of all commands or the info for a specific command
  // returns a multiline string
	if len(args) == 1 {
    // get command info from registry
    metadata := Registry.Commands[args[0].Bulk].CommandMetadata
    str := ""
    str += "name: " + metadata.Name + "\n"
    str += "summary: " + metadata.Summary + "\n"
    str += "complexity: " + metadata.Complexity + "\n"

		return resp.Value{
      Typ: "bulk",
      Bulk: str,
		}
	}

  var builder strings.Builder

  for _, cmd := range Registry.Commands {
    builder.WriteString(cmd.Name + "\n")
  }

  return resp.Value{
    Typ: "string",
    Str: builder.String(),
  }
}
