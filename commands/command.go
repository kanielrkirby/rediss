package commands

import (
	"fmt"
	"github.com/piratey7007/rediss/rerror"
	"github.com/piratey7007/rediss/resp"
	"strings"
)

func init() {
	Registry.Register("command", command)
	Registry.Register("command count", count)
	Registry.Register("command list", list)
	Registry.Register("command info", info)
	Registry.Register("command docs", docs)
	Registry.Register("command getkeys", getKeys)
	Registry.Register("command getkeysandflags", getKeysAndFlags)
	Registry.Register("command help", help)
}

func command(args []resp.Value) resp.Value {
	if len(args) != 0 {
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

func count(args []resp.Value) resp.Value {
	return resp.Value{
		Typ: "number",
		Num: len(Registry.Commands),
	}
}

func list(args []resp.Value) resp.Value {
	var builder strings.Builder

	for _, cmd := range Registry.Commands {
		builder.WriteString(cmd.Name + "\n")
	}

	return resp.Value{
		Typ: "string",
		Str: builder.String(),
	}
}

func info(args []resp.Value) resp.Value {
	arr := make([]resp.Value, len(Registry.Commands))
	for _, cmd := range Registry.Commands {
		arr = append(arr, resp.Value{
			Typ: "string",
			Str: cmd.Name,
		})
	}

	return resp.Value{
		Typ:   "array",
		Array: arr,
	}
}

func docs__getDoc(cmd Command) []string {
	//return fmt.Sprintf("Name: %sComplexity: %sSummary: %s", cmd.CommandMetadata.Name, cmd.CommandMetadata.Complexity, cmd.CommandMetadata.Summary)
	return []string{
		"Name",
		fmt.Sprintf("%s", cmd.CommandMetadata.Name),
		"Complexity",
		fmt.Sprintf("%s", cmd.CommandMetadata.Complexity),
		"Summary",
		fmt.Sprintf("%s", cmd.CommandMetadata.Summary),
	}
}

// docs returns the documentation for the given command.
func docs(args []resp.Value) resp.Value {
	if len(args) == 0 {
		allDocs := make([]resp.Value, 0, len(Registry.Commands))

		for _, cmd := range Registry.Commands {
			cmdDocs := docs__getDoc(cmd)

			respArr := make([]resp.Value, len(cmdDocs))
			for i, docStr := range cmdDocs {
				respArr[i] = resp.Value{
					Typ:  "string",
					Str: docStr,
				}
			}

			allDocs = append(allDocs, resp.Value{
				Typ:   "array",
				Array: respArr,
			})
		}

		return resp.Value{
			Typ:   "array",
			Array: allDocs,
		}

	}
	cmd, exists := Registry.Commands[args[0].Str]
	if !exists {
		return resp.Value{
			Typ: "error",
			Str: rerror.ErrInvalidArgument.Error(),
		}
	}
  doc := docs__getDoc(cmd)
  respArr := make([]resp.Value, len(doc))
  for i, docStr := range doc {
    respArr[i] = resp.Value{
      Typ: "string",
      Str: docStr,
    }
  }
  return resp.Value{
    Typ: "array",
    Array: respArr,
  }
}

func getKeys(args []resp.Value) resp.Value {
	return resp.Value{
		Typ: "error",
		Str: rerror.ErrUnimplemented.Error(),
	}
}

func getKeysAndFlags(args []resp.Value) resp.Value {
	return resp.Value{
		Typ: "error",
		Str: rerror.ErrUnimplemented.Error(),
	}
}

func help(args []resp.Value) resp.Value {
	var help = []string{
		"COMMAND <subcommand> [<arg> [value] [opt] ...]. Subcommands are:",
		"(no subcommand)",
		"    Return details about all Redis commands.",
		"COUNT",
		"    Return the total number of commands in this Redis server.",
		"LIST",
		"    Return a list of all commands in this Redis server.",
		"INFO [<command-name> ...]",
		"    Return details about multiple Redis commands.",
		"    If no command names are given, documentation details for all",
		"    commands are returned.",
		"DOCS [<command-name> ...]",
		"    Return documentation details about multiple Redis commands.",
		"    If no command names are given, documentation details for all",
		"    commands are returned.",
		"GETKEYS <full-command>",
		"    Return the keys from a full Redis command.",
		"GETKEYSANDFLAGS <full-command>",
		"    Return the keys and the access flags from a full Redis command.",
		"HELP",
		"    Print this help.",
	}
	for i, line := range help {
		line = fmt.Sprintf("%d) %s", i, line)
	}
	return resp.Value{
		Typ: "string",
		Str: strings.Join(help, "\n"),
	}
}
