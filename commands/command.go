package commands

import (
	"fmt"
	"github.com/piratey7007/rediss/resp"
	"strings"
)

func init() {
	Registry.Register("command", command)
}

func command(args []resp.Value) resp.Value {
	if len(args) == 0 {
		var builder strings.Builder

		for _, cmd := range Registry.Commands {
			builder.WriteString(cmd.Name + "\n")
		}

		return resp.Value{
			Typ: "string",
			Str: builder.String(),
		}
	}
	subcommand := strings.ToLower(args[0].Bulk)
  switch subcommand {
  case "count":
    return commandCount(args[1:])
  case "list":
    return commandList(args[1:])
  case "info":
    return commandInfo(args[1:])
  case "docs":
    return commandDocs(args[1:])
  case "getkeys":
    return commandGetKeys(args[1:])
  case "getkeysandflags":
    return commandGetKeysAndFlags(args[1:])
  case "help":
    return commandHelp(args[1:])
  default:
    return resp.Value{Typ: "error", Str: fmt.Sprintf("ERR unknown subcommand '%s'. Try COMMAND HELP.", subcommand)}
  }
}

func commandCount(args []resp.Value) resp.Value {
	return resp.Value{
		Typ: "integer",
		Num: len(Registry.Commands),
	}
}

func commandList(args []resp.Value) resp.Value {
	var builder strings.Builder

	for _, cmd := range Registry.Commands {
		builder.WriteString(cmd.Name + "\n")
	}

	return resp.Value{
		Typ: "string",
		Str: builder.String(),
	}
}

func commandInfo(args []resp.Value) resp.Value {
}

func commandDocs(args []resp.Value) resp.Value {
}

func commandGetKeys(args []resp.Value) resp.Value {
}

func commandGetKeysAndFlags(args []resp.Value) resp.Value {
}

func commandHelp(args []resp.Value) resp.Value {
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
