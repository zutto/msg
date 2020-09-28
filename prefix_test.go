package msg

import (
	"fmt"
	"testing"
)

//[0:1] init packet
//[1:3] version
//[3:5] length
//[5:7] header length
//[7:8] type
var static_data []byte = []byte{1, 1, 0, 0, 0, 34, 0, 8}

func TestPrefixGenerate(t *testing.T) {
	p := Prefix{
		Init:         1,
		Version:      1,
		Length:       0,
		HeaderLength: 34,
		Type:         8,
	}

	d, err := p.Generate()
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%+v", *d) != fmt.Sprintf("%+v", static_data) {
		t.Fatal(fmt.Errorf("prefix generated incorrect data?\nreceived: %+v\nexpected: %+v", *d, static_data))
	}
}

func TestPrefixParse(t *testing.T) {
	p := Prefix{}
	if err := p.Parse(&static_data, 0); err != nil {
		t.Fatal(err)
	}
	if p.Init != uint8(1) || p.Version != uint16(1) || p.Length != uint16(0) || p.HeaderLength != uint16(34) || p.Type != uint8(8) {
		t.Fatal(fmt.Errorf("incorrectly parsed data."))
	}

}
