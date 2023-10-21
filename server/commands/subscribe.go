package commands

import (
	"github.com/piratey7007/rediss/lib/rerror"
	"github.com/piratey7007/rediss/lib/resp"
)

func init() {
  Registry.Register("subscribe", subscribe)
}

func subscribe(args []resp.Value) resp.Value {
  if len(args) != 1 {
    return resp.Value{Typ: "error", Str: rerror.ErrWrongNumberOfArguments.FormatAndError("subscribe")}
  }

  channelName := args[0].Bulk

  SUBsMu.Lock()

  if SUBs[channelName] != nil {
    return resp.Value{Typ: "error", Str: rerror.ErrAlreadyExists.Error()}
  }

  _, exists := SUBs[channelName]
  if !exists {
    SUBs[channelName] = []chan resp.Value{}
  }
  SUBs[channelName] = append(SUBs[channelName], make(chan resp.Value))

  go func() {
    for {
      msg := <- SUBs[channelName][0]
      Registry.Writer.Write(msg)
    }
  }()

  SUBsMu.Unlock()

  return resp.Value{Typ: "string", Str: "OK"}
}
