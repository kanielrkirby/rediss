package commands

import (
	"fmt"

	"github.com/piratey7007/rediss/lib/rerror"
	"github.com/piratey7007/rediss/lib/resp"
)

func init() {
  Registry.Register("publish", publish)
}

func publish(args []resp.Value) resp.Value {
  if len(args) != 2 {
    return resp.Value{Typ: "error", Str: rerror.ErrWrongNumberOfArguments.FormatAndError("publish")}
  }

  channelName := args[0].Bulk
  message := args[1].Bulk

  fmt.Println("Publishing to channel", channelName, "message", message)

  SUBsMu.RLock()
  chanContainer, exists := SUBs[channelName]
  fmt.Println(len(chanContainer))
  SUBsMu.RUnlock()

  if !exists {
    return resp.Value{Typ: "error", Str: rerror.ErrNotFound.Error()}
  }

  for _, ch := range chanContainer {
    fmt.Println("Sending message to channel", channelName)
    ch <- resp.Value{Typ: "bulk", Bulk: message}
  }

  return resp.Value{Typ: "string", Str: "OK"}
}
