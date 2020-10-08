package msg

/*
const (
	_= iota // ignore first value by assigning to blank identifier
     KB float64 = 1 << (10 * iota)
     MB
     GB)

*/
import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"testing"
)

const (
	KB int = 1000
	MB     = 1000 * KB
)

func TestBox(t *testing.T) {
	var message []byte = make([]byte, 500*MB)
	rand.Read(message)
	e := NewEnvelope()
	//	e.Messagie.SetData(message)
	e.IncludeChecksum = false
	b := NewBox()

	ior, ioo := io.Pipe()
	var in io.Writer = ioo

	var out io.Reader = ior
	stream, err := NewStreamWriter(e, in)
	if err != nil {
		t.Fatal(err)
	}

	//
	var output_e *Envelope = NewEnvelope()
	ctx, cancel := context.WithCancel(context.Background())
	//	var checksum [20]byte
	go func(ctx context.Context, input io.Reader) {
		reader, err := NewStreamReader(output_e, input)
		if err != nil {
			t.Fatal(err)
		}
		for !b.IsFull() {
			output_e = NewEnvelope()
			reader.Reset(output_e, input)
			for !reader.IsFull() {
				_, err := reader.ReadWhole()
				if err != nil {
					t.Fatal(err)
				}
			}

			env := reader.GetEnvelope()
			/*ck, _ := env.Checksum(0, 0)
			copy(checksum[0:20], ck[:20])*/
			/*			if !env.Headers.ValidateChecksum(checksum) {
						fmt.Printf("%+v\nwrong checksum:\n%+v\n%+v\n", env, env.Headers.PacketChecksum, checksum)
						t.Fatal("wrong checksum!")
					}*/
			b.InsertEnvelope(env)
			//		fmt.Printf("read %d/%d\n", env.Headers.PacketIndex, env.Headers.TotalPackets)
		}
		cancel()
	}(ctx, out)

	fmt.Printf("writing message")
	_, err = stream.Write(message)
	if err != nil {
		t.Fatal(fmt.Errorf("stream writer errored: %+v", err))
	}
	<-ctx.Done()

	/*
		ne, err := b.RecornstructMessage()
		//	ne.Message.Data = &[]byte{}
		//	e.Message.Data = &[]byte{}
		e.Headers.PacketIndex = 0
		e.Headers.TotalPackets = 0
		orig, err := e.Checksum(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		nex, err := e.Checksum(0, 0)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("checksum of complete: \n%+v\nold:\n%+v\n\n%+v\n\n%+v", orig, nex, ne, e)*/

}
