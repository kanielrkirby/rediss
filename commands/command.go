package commands

import (
	"fmt"
	"github.com/piratey7007/rediss/resp"
	"strings"
	"github.com/piratey7007/rediss/rerror"
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
  return resp.Value{
    Typ: "error",
    Str: rerror.ErrUnimplemented.Error(),
  }
}

func docs(args []resp.Value) resp.Value {
  return resp.Value{
    Typ: "error",
    Str: rerror.ErrUnimplemented.Error(),
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
