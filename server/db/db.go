package db

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/piratey7007/rediss/lib/resp"
)

type Persistence interface {
  Load() error
  Read(fn func(value resp.Value)) error
  Write(value resp.Value) error
  Close(value resp.Value)
}

var DB Persistence

type Aof struct {
  file *os.File
  rd *bufio.Reader
  mu sync.Mutex
}

type Rdb struct {
  file *os.File
  interval int
  mutex sync.Mutex
  buffer bytes.Buffer
}

type AofOptions struct {
  Path string
  Interval int
}

type RdbOptions struct {
  Path string
  Interval int
}

func NewAof(options AofOptions) (*Aof, error) {
  f, err := os.OpenFile(options.Path, os.O_CREATE|os.O_RDWR, 0666)
  if err != nil {
    return nil, err
  }

  aof := &Aof{
    file: f,
    rd: bufio.NewReader(f),
  }

  go func() {
    for {
      aof.mu.Lock()
      aof.file.Sync()
      aof.mu.Unlock()
      time.Sleep(time.Duration(options.Interval) * time.Second)
    }
  }()

  return aof, nil
}

func (aof *Aof) Close(value resp.Value) {
  aof.mu.Lock()
  defer aof.mu.Unlock()

  aof.file.Close()
}

func (aof *Aof) Write(value resp.Value) error {
  aof.mu.Lock()
  defer aof.mu.Unlock()

  _, err := aof.file.Write(value.Marshal())
  if err != nil {
    return err
  }

  return nil
}

func (aof *Aof) Read(fn func(value resp.Value)) error {
  aof.mu.Lock()
  defer aof.mu.Unlock()

  aof.file.Seek(0, io.SeekStart)

  reader := resp.NewResp(aof.file)

  for {
    value, err := reader.Read()
    if err != nil {
      if err == io.EOF {
        break
      }
      return err
    }

    fn(value)
  }

  return nil
}

// Read from file and execute commands
func (aof *Aof) Load() error {


}


func NewRdb(options RdbOptions) (*Rdb, error) {
  f, err := os.OpenFile(options.Path, os.O_CREATE|os.O_RDWR, 0666)
  if err != nil {
    return nil, err
  }

  rdb := &Rdb{
    file: f,
    interval: options.Interval,
  }

  return rdb, nil
}

func (rdb *Rdb) Close() {
  rdb.mutex.Lock()
  defer rdb.mutex.Unlock()

  rdb.file.Close()
}

func (rdb *Rdb) Write() error {
  rdb.mutex.Lock()
  defer rdb.mutex.Unlock()

  rdb.buffer.Reset()
  rdb.writeHeader(6, 0)

  // this is where we write the keys from memory
  

  rdb.writeFooter()
}

func (rdb *Rdb) Load() error {
  rdb.mutex.Lock()
  defer rdb.mutex.Unlock()

  rdb.file.Seek(0, io.SeekStart)

  return nil
}

func (r *Rdb) writeHeader(version int, dbNumber byte) {
	r.buffer.WriteString("REDIS")

	versionStr := fmt.Sprintf("%04d", version)
	r.buffer.WriteString(versionStr)
	r.buffer.WriteByte(0xFE)
	r.buffer.WriteByte(dbNumber)
}

func (w *Rdb) writeFooter() {
  w.buffer.WriteByte(0xFF)

  var checksum uint64 = 0
  binary.Write(&w.buffer, binary.LittleEndian, checksum)
}
