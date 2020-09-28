package msg

import (
	"fmt"
	//  "github.com/google/gofuzz"
	//"fmt"
	"testing"
)

var envlope_checksum []byte = []byte{240, 92, 160, 21, 230, 160, 253, 136, 229, 100, 234, 13, 5, 184, 73, 246, 131, 58, 106, 41}

func TestEnvelope(t *testing.T) {

	d := []byte("asd i am data")
	e := NewEnvelope()
	e.Message.SetData(d)
	//	e.Message.CompressionType = GZIP
	_, _, err := e.Generate()
	if err != nil {
		t.Fatal(err)
	}
	ck, err := e.Checksum()

	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%+v", ck) != fmt.Sprintf("%+v", envlope_checksum) {
		t.Fatal(fmt.Sprintf("checksum mismatch, expected %+v, got %+v", ck, envlope_checksum))
	}
	//fmt.Printf("output: %+v\n", *out)
}
