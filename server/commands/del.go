package commands

import (
  "github.com/piratey7007/rediss/lib/resp"
  "github.com/piratey7007/rediss/lib/rerror"
)

func init() {
  Registry.Register("del", del)
}

func del(args []resp.Value) resp.Value {
  if len(args) != 1 {
    return resp.Value{Typ: "error", Str: rerror.ErrWrongNumberOfArguments.FormatAndError("del")}
  }

  key := args[0].Bulk

  SETsMu.Lock()
  delete(SETs, key)
  SETsMu.Unlock()

  return resp.Value{Typ: "string", Str: "OK"}
}


