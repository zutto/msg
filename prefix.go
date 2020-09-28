package msg

import (
	"encoding/binary"
	"fmt"
)

type Prefix struct {
	Init         uint8  // [0:1] always 1?
	Version      uint16 // [1:3] 1..to be changed?
	Length       uint16 // [3:5] length of whole packet
	HeaderLength uint16 // [5:7] length of headers
	Type         uint8  // [7:8] packet type
}

func (p *Prefix) Generate() (*[]byte, error) {
	var data []byte = make([]byte, 8)
	data[0] = byte(INIT_BYTE)

	binary.LittleEndian.PutUint16(data[1:3], uint16(p.Version))
	binary.LittleEndian.PutUint16(data[3:5], uint16(p.Length))
	binary.LittleEndian.PutUint16(data[5:7], uint16(p.HeaderLength))

	data[7] = byte(p.Type)

	//fmt.Printf("prefix: %+v\n", data)
	return &data, nil
}

func (p *Prefix) Parse(data *[]byte, pad int) error {

	p.Init = uint8((*data)[pad+0])
	p.Version = binary.LittleEndian.Uint16((*data)[pad+1 : pad+3])
	p.Length = binary.LittleEndian.Uint16((*data)[pad+3 : pad+5])
	p.HeaderLength = binary.LittleEndian.Uint16((*data)[pad+5 : pad+7])
	p.Type = uint8((*data)[pad+7])

	if p.Init != uint8(INIT_BYTE) {
		return fmt.Errorf("Invalid init byte. Received %d, expected %d.", p.Init, INIT_BYTE)
	}
	return nil
}
