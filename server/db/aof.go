package db

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/piratey7007/rediss/lib/rerror"
	"github.com/piratey7007/rediss/lib/resp"
	"github.com/piratey7007/rediss/server/commands"
)

type Aof struct {
  File *os.File
  Rd *bufio.Reader
  Mu sync.Mutex
  Path string
  Interval int
}

type AofOptions struct {
  Path string
  Interval int
}

func NewAof(options AofOptions) (*Aof, error) {
  f, err := os.OpenFile(options.Path, os.O_CREATE|os.O_RDWR, 0666)
  if err != nil {
    return nil, rerror.ErrWrap(err)
  }

  aof := &Aof{
    File: f,
    Rd: bufio.NewReader(f),
  }

  go func() {
    for {
      aof.Mu.Lock()
      aof.File.Sync()
      aof.Mu.Unlock()
      time.Sleep(time.Duration(options.Interval) * time.Second)
    }
  }()

  return aof, nil
}

func (a *Aof) Close(value resp.Value) {
  a.Mu.Lock()
  defer a.Mu.Unlock()

  a.File.Close()
}

func (a *Aof) Write(value resp.Value) error {
  a.Mu.Lock()
  defer a.Mu.Unlock()

  _, err := a.File.Write(value.Marshal())
  if err != nil {
    return rerror.ErrWrap(err)
  }

  return nil
}

func (a *Aof) Read() error {
  a.Mu.Lock()
  defer a.Mu.Unlock()

  a.File.Seek(0, io.SeekStart)

  reader := resp.NewResp(a.File)

  for {
    value, err := reader.Read()
    if err != nil {
      if err == io.EOF {
        break
      }
      return rerror.ErrWrap(err)
    }

    commandName := value.Array[0].Bulk
    args := value.Array[1:]

    cmd, exists := commands.Registry.Commands[commandName]
    if !exists {
      return rerror.ErrWrap(errors.New("Invalid command AOF read %s")).Format(commandName)
    }

    cmd.Execute(args)
  }

  return nil
}
