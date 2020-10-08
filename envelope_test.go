package msg

import (
	"fmt"
	//  "github.com/google/gofuzz"
	//"fmt"
	"testing"
)

//var envlope_checksum []byte = []byte{240, 92, 160, 21, 230, 160, 253, 136, 229, 100, 234, 13, 5, 184, 73, 246, 131, 58, 106, 41}
//var envlope_checksum []byte = []byte{221, 140, 108, 189, 57, 213, 44, 170, 200, 157, 190, 49, 141, 181, 252, 118, 80, 109, 70, 74}
var envlope_checksum []byte = []byte{131, 79, 205, 172, 144, 237, 156, 209, 53, 122, 150, 34, 169, 178, 226, 124, 154, 146, 119, 253}

func TestEnvelope(t *testing.T) {

	d := []byte("asd i am data")
	e := NewEnvelope()
	e.Message.SetData(d)
	//	e.Message.CompressionType = GZIP
	_, _, err := e.Generate()
	if err != nil {
		t.Fatal(err)
	}

	limit := 0
	if e.Message.Len() > e.MessageSizeLimit {
		limit = e.MessageSizeLimit
	} else {
		limit = e.Message.Len()
	}

	limit = limit - e.EnvelopeLabels.Len()
	if limit < 1 {
		t.Fatal(fmt.Sprintf("Envelope labels are longer than the message limit."))
	}
	ck, err := e.Checksum(0, 0+limit)

	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%+v", ck) != fmt.Sprintf("%+v", envlope_checksum) {
		t.Fatal(fmt.Sprintf("checksum mismatch, expected \n%+v, got \n%+v", envlope_checksum, ck))
	}
	//fmt.Printf("output: %+v\n", *out)
}
