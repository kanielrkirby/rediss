package commands

func del(args []Value) Value {
  if len(args) != 1 {
    return Value{typ: "error", str: "ERR wrong number of arguments for 'del' command"}
  }

  key := args[0].bulk

  SETsMu.Lock()
  delete(SETs, key)
  SETsMu.Unlock()

  return Value{typ: "string", str: "OK"}
}
