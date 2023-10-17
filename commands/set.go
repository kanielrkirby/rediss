package commands

import (
  "github.com/piratey7007/rediss/resp"
  "github.com/piratey7007/rediss/rerror"
)

func init() {
  Registry.Register("set", set)
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: rerror.ArgumentCount("set")}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

