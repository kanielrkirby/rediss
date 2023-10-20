package db

import (
	"github.com/piratey7007/rediss/lib/resp"
)

type Persistence interface {
  Load() error
  Read(fn func(value resp.Value)) error
  Write(value resp.Value) error
  Close(value resp.Value)
}

var DB Persistence

