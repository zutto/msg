package msg

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"github.com/pierrec/lz4"
	"sync"
)

//Compression types
const (
	NIL = iota
	LZ4
	GZIP
)

//Message contains the (optional) user defined data for the envelope.
type Message struct {
	Data *[]byte //Data of any size

	Compressed      bool  //true false if current data is compressed
	CompressionType uint8 //compression type used/to be used.

	lock sync.Mutex
}

//Parse parses input data, data can be read from padded starting position.
func (m *Message) Parse(data *[]byte, pad int) error {
	if uint8((*data)[pad+0]) == uint8(1) {
		m.Compressed = true
	}

	m.CompressionType = (*data)[pad+1]
	var d []byte = (*data)[pad+MESSAGE_PREFIX_LEN:]
	m.Data = &d
	return nil
}

func (m *Message) AppendData(data []byte) error {
	nd := append(*m.Data, data...)
	m.Data = &nd
	return nil
}

//SetData is convenience function to set the data.
func (m *Message) SetData(d []byte) error {
	m.Data = &d
	return nil
}

func (m *Message) ChecksumFromBytes(from, to int) ([]byte, error) {
	data, _, err := m.GetData(from, to)
	if err != nil {
		return nil, err
	}

	output := sha1.Sum((*data)[:])
	return output[:], nil

}

//Checksum generates SHA-1 checksum from the message.
func (m *Message) Checksum() ([]byte, error) {
	var cmp uint8 = 0
	if m.Compressed {
		cmp = 1
	}

	var data []byte = []byte{byte(cmp), byte(m.CompressionType)}
	output := sha1.Sum(append(data, (*m.Data)[:]...))

	return output[:], nil
}

//GetData retusn data from X position to Y position. Returns data, length, and error.
func (m *Message) GetData(from, to int) (*[]byte, int, error) {
	if m.Data == nil {
		return nil, 0, errors.New("no data in message.")
	}
	var cmp uint8 = 0
	if m.Compressed {
		cmp = 1
	}
	var result []byte = []byte{byte(cmp), byte(m.CompressionType)}

	if len(*m.Data) < to {
		to = len(*m.Data)
	}

	if from < 0 {
		return nil, 0, errors.New("cannot read from index less than zero")
	}

	if from > to {
		return nil, 0, errors.New("from is higher than to?")
	}

	result = append(result, (*m.Data)[from:to]...)
	return &result, len(result) - MESSAGE_PREFIX_LEN, nil
}

func (m *Message) Len() int {
	if m.Data != nil {
		return len(*m.Data)
	}
	return 0
}

//Compress compresses the current data and marks it as compressed.
func (m *Message) Compress(level int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.Compressed {
		return nil
	}
	switch m.CompressionType {
	case LZ4:
		return m.lz4Compress(level)
	case GZIP:
		return m.gzipCompress(level)
	default:
		return nil
	}
}

//DeCompress decompresses the current data and marks it as uncompressed.
func (m *Message) DeCompress() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if !m.Compressed {
		return nil
	}

	switch m.CompressionType {
	case LZ4:
		return m.lz4DeCompress()
	case GZIP:
		return m.gzipDeCompress()
	default:
		return nil
	}
}

func (m *Message) lz4Compress(level int) error {
	var buffer bytes.Buffer
	lz := lz4.NewWriter(&buffer)
	defer lz.Close() //close here, lz4 flush does not work correctly and duplicates data if close is called..
	var written int = 0
	var length int = len((*m.Data))

	for written < length {
		w, err := lz.Write((*m.Data)[written:length])
		if err != nil {
			return err
		}
		written += w
	}
	if err := lz.Flush(); err != nil {
		return err
	}

	newData := buffer.Bytes()
	m.Data = &newData
	m.Compressed = true

	return nil
}

func (m *Message) lz4DeCompress() error {
	var buffer *bytes.Buffer = bytes.NewBuffer((*m.Data))

	lz := lz4.NewReader(buffer)

	var result bytes.Buffer
	_, err := result.ReadFrom(lz)

	if err != nil {
		return err
	}
	newData := result.Bytes()
	m.Data = &newData
	m.Compressed = false
	return nil
}

func (m *Message) gzipCompress(level int) error {
	var buffer bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buffer, level)
	if err != nil {
		return err
	}

	var written int = 0
	var length int = len((*m.Data))
	for written < length {
		w, err := gz.Write((*m.Data)[written:length])
		if err != nil {
			return err
		}
		written += w
	}

	if err := gz.Flush(); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	newData := buffer.Bytes()
	m.Data = &newData
	m.Compressed = true
	return nil
}

func (m *Message) gzipDeCompress() error {
	var buffer *bytes.Buffer = bytes.NewBuffer((*m.Data))

	gz, err := gzip.NewReader(buffer)
	if err != nil {
		return err
	}

	var result bytes.Buffer
	if _, err := result.ReadFrom(gz); err != nil {
		return err
	}

	newData := result.Bytes()
	m.Data = &newData
	m.Compressed = false
	return nil
}
