package msg

import (
	"fmt"
	"testing"
)

//[0:4] index of packet
//[4:8] packet count
//[8:16] total length of message (not the packet..)
var static_header []byte = []byte{1, 0, 0, 0, 2, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0}

var static_header2 []byte = []byte{2, 0, 0, 0, 4, 0, 0, 0, 200, 0, 0, 0, 0, 0, 0, 0}

func TestGenerate(t *testing.T) {
	h := Headers{
		PacketIndex: 0,
		TotalLength: 100,
	}
	h.IncrementIndex()
	h.CalculatePacketCount(50)
	d, err := h.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%+v", *d) != fmt.Sprintf("%+v", static_header) {
		t.Fatal(fmt.Errorf("generated does not match static."))
	}

	//2
	h.TotalLength = 200
	h.IncrementIndex()
	h.CalculatePacketCount(50)
	d, err = h.Generate()
	if err != nil {
		t.Fatal(err)
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
