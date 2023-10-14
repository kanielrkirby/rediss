package commands

import (
  "sync"
  "github.com/piratey7007/rediss/resp"
)

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}
