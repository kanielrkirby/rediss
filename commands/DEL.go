package commands

import (
  "github.com/piratey7007/rediss/resp"
)

func init() {
  Registry.Register("del", del)
}

func del(args []resp.Value) resp.Value {
  if len(args) != 1 {
    return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'del' command"}
  }

  key := args[0].Bulk

  SETsMu.Lock()
  delete(SETs, key)
  SETsMu.Unlock()

  return resp.Value{Typ: "string", Str: "OK"}
}


