package commands

import (
  "github.com/piratey7007/rediss/server/resp"
  "github.com/piratey7007/rediss/server/rerror"
)

func init() {
  Registry.Register("getset", getset)
}

func getset(args []resp.Value) resp.Value {
  if len(args) != 2 {
    return resp.Value{Typ: "error", Str: rerror.ErrWrongNumberOfArguments.FormatAndError("getset")}
  }

  key := args[0].Bulk
  value := args[1].Bulk

  SETsMu.Lock()
  prev := SETs[key]
  SETs[key] = value
  SETsMu.Unlock()

  return resp.Value{Typ: "bulk", Bulk: prev}
}
