package commands

import (
  "fmt"
  "github.com/piratey7007/rediss/rerror"
  "github.com/piratey7007/rediss/resp"
)

func init() {
  Registry.Register("get", get)
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: rerror.ErrWrongNumberOfArguments.FormatAndError("get")}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
    fmt.Println("Key not found")
		return resp.Value{Typ: "null"}
	}
  fmt.Println("Value: ", value)

	return resp.Value{Typ: "bulk", Bulk: value}
}
