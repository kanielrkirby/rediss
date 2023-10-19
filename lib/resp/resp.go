package resp

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/piratey7007/rediss/lib/rerror"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readLine reads a line from the reader until it reaches \r\n.
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, rerror.ErrWrap(err).Format("Error reading byte")
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// readInteger reads the line and parses the result as an integer.
func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, rerror.ErrWrap(err).Format("Error reading line")
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, rerror.ErrWrap(err).Format("Error parsing integer")
	}
	return int(i64), n, nil
}

// Read reads a RESP value from the reader. This is the primary entry point for
// reading RESP values. This is a blocking call.
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

  if err != nil {
    return Value{}, rerror.ErrWrap(err).Format("Error reading type")
  }

  if _type == 0 {
    return Value{}, errors.New("EOF")
  }

	switch _type {
	case ARRAY:
		return r.readArrayAsValue()
	case BULK:
		return r.readBulkAsValue()
  case STRING:
    return r.readStringAsValue()
  case INTEGER:
    return r.readIntegerAsValue()
  case ERROR:
    return r.readErrorAsValue()
	default:
		fmt.Printf(rerror.ErrUnknownType.FormatAndError(string(_type)))
		return Value{}, nil
	}
}

// readArray reads an array from the reader, and reads each of the values in
// the array. Formatting: *<length>\r\n<values>
func (r *Resp) readArrayAsValue() (Value, error) {
	v := Value{
    Typ: "array",
  }

	len, _, err := r.readInteger()
	if err != nil {
		return v, rerror.ErrWrap(err).Format("Error reading array length")
	}

	v.Array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, rerror.ErrWrap(err).Format("Error reading array value")
		}

		v.Array = append(v.Array, val)
	}

	return v, nil
}

// readBulk reads a bulk string from the reader.
// Formatting: $<length>\r\n<bytes>\r\n
func (r *Resp) readBulkAsValue() (Value, error) {
	v := Value{}

	v.Typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, rerror.ErrWrap(err).Format("Error reading bulk length")
	}

  if len < 0 {
    v.Typ = "null"
    return v, nil
  }

	Bulk := make([]byte, len)

	_, err = r.reader.Read(Bulk)
	if err != nil {
		return v, rerror.ErrWrap(err).Format("Error reading bulk")
	}

	v.Bulk = string(Bulk)

	_, _, err = r.readLine()
	if err != nil {
		return v, rerror.ErrWrap(err).Format("Error reading line")
	}

	return v, nil
}

// readString reads a string from the reader.
// Formatting: +<string>\r\n
func (r *Resp) readStringAsValue() (Value, error) {
  v := Value{
    Typ: "string",
  }

  line, _, err := r.readLine()
  if err != nil {
    return v, rerror.ErrWrap(err).Format("Error reading string")
  }

  v.Str = string(line)

  return v, nil
}

// readError reads an error from the reader.
// Formatting: -<string>\r\n
func (r *Resp) readErrorAsValue() (Value, error) {
  v := Value{
    Typ: "error",
  }

  line, _, err := r.readLine()
  if err != nil {
    return v, rerror.ErrWrap(err).Format("Error reading error")
  }

  v.Str = string(line)

  return v, nil
}

// readInteger reads an integer from the reader.
// Formatting: :<integer>\r\n
func (r *Resp) readIntegerAsValue() (Value, error) {
  v := Value{
    Typ: "int",
  }

  num, _, err := r.readInteger()
  if err != nil {
    return v, rerror.ErrWrap(err).Format("Error reading integer")
  }

  v.Num = num

  return v, nil
}

// Marshal returns the RESP encoding of the value.
func (v Value) Marshal() []byte {
	switch v.Typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		return []byte{}
	}
}

// marshalString returns the RESP encoding of a string.
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshalBulk returns the RESP encoding of a bulk string.
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshalArray returns the RESP encoding of an array.
func (v Value) marshalArray() []byte {
	len := len(v.Array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.Array[i].Marshal()...)
	}

	return bytes
}

// marshallError returns the RESP encoding of an error.
func (v Value) marshallError() []byte {
	if rerror.DEBUG {
		fmt.Print(rerror.ErrWrap(errors.New(v.Str)).FormatAndError(v.Str))
	}

	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshallNull returns the RESP encoding of a null.
func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}

type Writer struct {
	writer io.Writer
}

// NewWriter returns a new RESP writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// Write writes a RESP value to the writer.
func (w *Writer) Write(v Value) error {
  fmt.Println("Writing value start: ", v.Typ, v.Str, v.Num, v.Bulk, v.Array)
	var bytes = v.Marshal()
  fmt.Println("Marshalled: ", bytes, string(bytes))


	_, err := w.writer.Write(bytes)
	if err != nil {
		return rerror.ErrWrap(err).Format("Error writing value")
	}

	return nil
}
