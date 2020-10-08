package msg

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math"
)

//Headers contain the headers of the message, which includes the index of the packet (if multipart packet), total number of packets, and total length of the combined data contained of all the packets.
type Headers struct {
	PacketIndex    uint32   //[0:4] number of packets in the message
	TotalPackets   uint32   //[4:8] total number of packets in this message
	TotalLength    uint64   //[8:16] total length of the message
	PacketChecksum [20]byte //[16:36] Checksum (potentially zeroes.)
}

//Generates generates the headers for envelope.
func (h *Headers) Generate() (*[]byte, error) {
	var data []byte = make([]byte, HEADER_LEN)
	binary.LittleEndian.PutUint32(data[0:4], uint32(h.PacketIndex))
	binary.LittleEndian.PutUint32(data[4:8], uint32(h.TotalPackets))
	binary.LittleEndian.PutUint64(data[8:16], uint64(h.TotalLength))
	copy(data[16:], h.PacketChecksum[:])
	return &data, nil
}

func (h *Headers) IncrementIndex() {
	h.PacketIndex++ //???
}

func (h *Headers) CalculatePacketCount(size int) {
	h.TotalPackets = uint32(math.Ceil(float64(h.TotalLength) / float64(size)))
}

//SetChecksum is a convenience function to set checksum
func (h *Headers) SetChecksum(checksum []byte) error {
	copy(h.PacketChecksum[:], checksum[:20])
	return nil
}

func (h *Headers) ValidateChecksum(checksum [20]byte) bool {
	if fmt.Sprintf("%+v", h.PacketChecksum) != fmt.Sprintf("%+v", checksum) {
		return false
	}
	return true
}

//Parse parses input data, data can be read from padded starting position.
func (h *Headers) Parse(data *[]byte, pad int) error {
	h.PacketIndex = binary.LittleEndian.Uint32((*data)[pad+0 : pad+4])
	h.TotalPackets = binary.LittleEndian.Uint32((*data)[pad+4 : pad+8])
	h.TotalLength = binary.LittleEndian.Uint64((*data)[pad+8 : pad+16])
	copy(h.PacketChecksum[:], (*data)[pad+16:pad+36])
	return nil
}

//Checksum returns SHA-1 checksum of the headers.
func (h *Headers) Checksum() ([]byte, error) {
	data, err := h.Generate()
	if err != nil {
		return nil, err
	}

	output := sha1.Sum((*data)[0 : len(*data)-20])
	return output[:], nil

}
