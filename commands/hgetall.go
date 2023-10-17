package commands

import (
  "github.com/piratey7007/rediss/resp"
  "github.com/piratey7007/rediss/rerror"
)

func init() {
  Registry.Register("hgetall", hgetall)
}

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: rerror.ErrWrongNumberOfArguments.FormatAndError("hgetall")}
	}

	hash := args[0].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	values := []resp.Value{}
	for k, v := range value {
		values = append(values, resp.Value{Typ: "bulk", Bulk: k})
		values = append(values, resp.Value{Typ: "bulk", Bulk: v})
	}

	return resp.Value{Typ: "array", Array: values}
}
