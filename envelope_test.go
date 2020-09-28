package msg

import (
	/*"fmt"
	  "github.com/google/gofuzz"*/
	//"fmt"
	"testing"
)

func TestEnvelope(t *testing.T) {

	d := []byte("asd i am data")
	e := NewEnvelope()
	e.Message.SetData(d)
	//	e.Message.CompressionType = GZIP
	_, _, err := e.Generate()
	if err != nil {
		t.Fatal(err)
	}

	//fmt.Printf("output: %+v\n", *out)
}
