package db

// these are where the data is stored in memory
// var SETs = map[string]string{}
// var SETsMu = sync.RWMutex{}
// var HSETs = map[string]map[string]string{}
// var HSETsMu = sync.RWMutex{}

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/piratey7007/rediss/server/commands"
)

// Semantic Bytes for RDB
const (
	rdb6bitLen             = 0
	rdb14bitLen            = 1
	rdb32bitLen            = 2
	rdbEncVal              = 3
	rdbOpcodeExpiry        = 0xFD
	rdbOpcodeSelectDB      = 0xFE
	rdbOpcodeEOF           = 0xFF
	rdbStringEnc           = 0
	rdbListEnc             = 1
	rdbSetEnc              = 2
	rdbSortedSetEnc        = 3
	rdbHashEnc             = 4
	rdbZipmapEnc           = 9
	rdbZiplistEnc          = 10
	rdbIntsetEnc           = 11
	rdbSortedSetZiplistEnc = 12
	rdbHashmapZiplistEnc   = 13
  rdb6bitStr             = 0
  rdb14bitStr            = 1
  rdb32bitStr            = 2
)

type Rdb struct {
	file     *os.File
	mutex    sync.Mutex
	interval int
	buffer   bytes.Buffer
}

type RdbOptions struct {
	Path     string
	Interval int
}

func NewRdb(options RdbOptions) (*Rdb, error) {
	f, err := os.OpenFile(options.Path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	rdb := &Rdb{
		file:     f,
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

	for key, value := range commands.SETs {
    typeByte := getEncodingByte(value)
    if err := rdb.buffer.WriteByte(typeByte); err != nil {
      return err
    }

    if err := rdb.WriteString(key); err != nil {
      return err
    }

    if err := rdb.WriteEncoded(value); err != nil {
    }
  }

	rdb.writeFooter()

	return nil
}

func getEncodingByte(value interface{}) byte {
  // get type of value and return the appropriate byte like this:
  switch value.(type) {
  case string:
    return rdbStringEnc
  case []string:
    return rdbListEnc
  case int:
    return rdbIntsetEnc
  default:
    return rdbStringEnc
  }
}

func (rdb *Rdb) WriteString(str string) error {
	// Convert the string to a byte slice because we're dealing with binary data.
	data := []byte(str)
	length := len(data)

	// Determine the encoding based on the string's length.
	// Then, write the appropriate header and the string itself.
	switch {
	case length <= 63:
		// For strings with length <= 63
		header := byte(rdb6bitStr<<6) | byte(length) // Encoding the length directly into the header
		if err := rdb.buffer.WriteByte(header); err != nil {
			return err
		}

	case length <= 16383:
		// For strings with length <= 16383
		header := byte(rdb14bitStr<<6) | byte(length>>8)&0x3F | byte(length)&0xFF // Splitting length into two bytes
		if _, err := rdb.buffer.Write([]byte{header}); err != nil {
			return err
		}

	default:
		// For longer strings, we might use a different strategy (like 32-bit length)
		// Here we write a fixed header and then the length as a 32-bit integer.
		header := byte(rdb32bitStr << 6)
		if err := rdb.buffer.WriteByte(header); err != nil {
			return err
		}
		lengthBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthBytes, uint32(length)) // Converting length to bytes
		if _, err := rdb.buffer.Write(lengthBytes); err != nil {
			return err
		}
	}

	// After writing the appropriate header, we write the actual string data.
	if _, err := rdb.buffer.Write(data); err != nil {
		return err
	}

	return nil
}

func (rdb *Rdb) Load() error {
	rdb.mutex.Lock()
	defer rdb.mutex.Unlock()

	rdb.file.Seek(0, io.SeekStart)

	var magic [5]byte
	_, err := rdb.file.Read(magic[:])
	if err != nil {
		return err
	}

	if string(magic[:]) != "REDIS" {
		return fmt.Errorf("Invalid magic string: %s", string(magic[:]))
	}

	var version uint16
	err = binary.Read(rdb.file, binary.LittleEndian, &version)
	if err != nil {
		return err
	}

	if version != 6 {
		return fmt.Errorf("Invalid version: %d", version)
	}

	var dbNumber byte
	err = binary.Read(rdb.file, binary.LittleEndian, &dbNumber)
	if err != nil {
		return err
	}

	var opcode byte
	err = binary.Read(rdb.file, binary.LittleEndian, &opcode)
	if err != nil {
		return err
	}

	for opcode != rdbOpcodeEOF {
		switch opcode {
		case rdbOpcodeExpiry:
			var expiry uint64
			err = binary.Read(rdb.file, binary.LittleEndian, &expiry)
			if err != nil {
				return err
			}
		case rdbOpcodeSelectDB:
			var dbNumber byte
			err = binary.Read(rdb.file, binary.LittleEndian, &dbNumber)
			if err != nil {
				return err
			}
		default:
			var key []byte
			err = binary.Read(rdb.file, binary.LittleEndian, &key)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// writeHeader is a helper function to write the header of the RDB file.
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
