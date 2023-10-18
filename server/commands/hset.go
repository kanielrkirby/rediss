package commands

import (
  "github.com/piratey7007/rediss/lib/resp"
  "github.com/piratey7007/rediss/lib/rerror"
)

func init() {
  Registry.Register("hset", hset)
}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: rerror.ErrWrongNumberOfArguments.FormatAndError("hset")}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

