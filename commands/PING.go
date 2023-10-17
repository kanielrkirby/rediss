package commands

import (
  "github.com/piratey7007/rediss/resp"
  "github.com/piratey7007/rediss/errorHandler"
)

func init() {
  Registry.Register("ping", ping)
}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Typ: "string", Str: "PONG"}
	}

  return resp.Value{Typ: "error", Str: errorHandler.ArgumentCount("ping")}
}
