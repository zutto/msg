package msg

import (
	"crypto/sha1"
	"encoding/binary"
	"math"
)

type Headers struct {
	PacketIndex  uint32 //[0:4] number of packets in the message
	TotalPackets uint32 //[4:8] total number of packets in this message
	TotalLength  uint64 //[8:16] total length of the message
}

func (h *Headers) Generate() (*[]byte, error) {
	var data []byte = make([]byte, 16)
	binary.LittleEndian.PutUint32(data[0:4], uint32(h.PacketIndex))
	binary.LittleEndian.PutUint32(data[4:8], uint32(h.TotalPackets))
	binary.LittleEndian.PutUint64(data[8:16], uint64(h.TotalLength))
	return &data, nil
}
func (h *Headers) IncrementIndex() {
	h.PacketIndex++ //???
}

func (h *Headers) CalculatePacketCount(size int) {
	h.TotalPackets = uint32(math.Ceil(float64(h.TotalLength) / float64(size)))
}

func (h *Headers) Parse(data *[]byte, pad int) error {
	h.PacketIndex = binary.LittleEndian.Uint32((*data)[pad+0 : pad+4])
	h.TotalPackets = binary.LittleEndian.Uint32((*data)[pad+4 : pad+8])
	h.TotalLength = binary.LittleEndian.Uint64((*data)[pad+8 : pad+16])
	return nil
}

func (h *Headers) Checksum() ([]byte, error) {
	data, err := h.Generate()
	if err != nil {
		return nil, err
	}

	output := sha1.Sum(*data)
	return output[:], nil

}
