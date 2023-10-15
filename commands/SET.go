package commands

import (
  "github.com/piratey7007/rediss/resp"
)

func init() {
  Registry.Register("set", set)
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

