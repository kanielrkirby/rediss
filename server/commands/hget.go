package commands

import (
  "github.com/piratey7007/rediss/lib/resp"
  "github.com/piratey7007/rediss/lib/rerror"
)

func init() {
  Registry.Register("hget", hget)
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: rerror.ErrWrongNumberOfArguments.FormatAndError("hget")}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{Typ: "bulk", Bulk: value}
}

