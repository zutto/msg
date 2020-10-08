package msg

import (
	"fmt"
	"testing"
)

//[0:4] index of packet
//[4:8] packet count
//[8:16] total length of message (not the packet..)
//var static_header []byte = []byte{1, 0, 0, 0, 2, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0}
var static_header []byte = []byte{1, 0, 0, 0, 2, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 47, 87, 211, 152, 84, 19, 179, 49, 227, 82, 190, 190, 42, 71, 74, 182, 82, 164, 249, 163}

//111, 178, 187, 51, 168, 51, 170, 150, 45, 90, 209, 209, 67, 198, 126, 222, 17, 67, 4, 229}i

//var header_checksum1 []byte = []byte{47, 87, 211, 152, 84, 19, 179, 49, 227, 82, 190, 190, 42, 71, 74, 182, 82, 164, 249, 163}
var header_checksum1 []byte = []byte{47, 87, 211, 152, 84, 19, 179, 49, 227, 82, 190, 190, 42, 71, 74, 182, 82, 164, 249, 163}

//var static_header2 []byte = []byte{2, 0, 0, 0, 4, 0, 0, 0, 200, 0, 0, 0, 0, 0, 0, 0}
var static_header2 []byte = []byte{2, 0, 0, 0, 4, 0, 0, 0, 200, 0, 0, 0, 0, 0, 0, 0, 47, 87, 211, 152, 84, 19, 179, 49, 227, 82, 190, 190, 42, 71, 74, 182, 82, 164, 249, 163}

//var header_checksum2 []byte = []byte{23, 58, 179, 138, 114, 246, 245, 157, 181, 144, 119, 119, 233, 124, 22, 216, 183, 166, 86, 218}
var header_checksum2 []byte = []byte{23, 58, 179, 138, 114, 246, 245, 157, 181, 144, 119, 119, 233, 124, 22, 216, 183, 166, 86, 218}

func TestGenerate(t *testing.T) {
	h := Headers{
		PacketIndex: 0,
		TotalLength: 100,
	}
	h.IncrementIndex()
	h.CalculatePacketCount(50)
	checksum, err := h.Checksum()
	if err != nil {
		t.Fatal(err)
	}

	if err := h.SetChecksum(checksum); err != nil {
		t.Fatal(err)
	}

	d, err := h.Generate()
	if err != nil {
		t.Fatal(err)
	}
	//	ckk, _ := h.Checksum()
	//	fmt.Printf("headers: %+v\n%+v\n", ckk, *d)
	if fmt.Sprintf("%+v", *d) != fmt.Sprintf("%+v", static_header) {
		t.Fatal(fmt.Errorf("generated does not match static."))
	}
	ck, err := h.Checksum()
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%+v", ck) != fmt.Sprintf("%+v", header_checksum1) {
		t.Fatal(fmt.Sprintf("checksum mismatch, expected \n%+v\n got \n%+v", header_checksum1, ck))
	}
	//2
	h.TotalLength = 200
	h.IncrementIndex()
	h.CalculatePacketCount(50)
	d, err = h.Generate()
	if err != nil {
		t.Fatal(err)
	}
	ck, err = h.Checksum()
	if err != nil {
		t.Fatal(err)
	}

	//	fmt.Printf("\nxx %+v\n%+v", ck, *d)
	if fmt.Sprintf("%+v", ck) != fmt.Sprintf("%+v", header_checksum2) {
		t.Fatal(fmt.Sprintf("checksum mismatch, expected %+v, got %+v", ck, header_checksum2))
	}

	if fmt.Sprintf("%+v", *d) != fmt.Sprintf("%+v", static_header2) {
		t.Fatal(fmt.Errorf("generated does not match static."))
	}
}

func TestParse(t *testing.T) {
	h := Headers{}
	err := h.Parse(&static_header, 0)
	if err != nil {
		t.Fatal(err)
	}

	if h.PacketIndex != uint32(1) || h.TotalPackets != uint32(2) || h.TotalLength != uint64(100) {
		t.Fatal(fmt.Errorf("failed to parse data.."))
	}

}
