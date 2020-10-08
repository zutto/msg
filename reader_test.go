package msg

import (
	//	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	e := NewEnvelope()
	e.EnvelopeLabels.MessageId = []byte("paskaa")
	e.EnvelopeLabels.MessageFeature = []byte("kilokaupalla")
	e.Message.SetData([]byte("1"))
	ior, ioo := io.Pipe()
	var in io.Writer = ioo
	var out io.Reader = ior
	stream, err := NewStreamWriter(e, in)
	if err != nil {
		t.Fatal(err)
	}
	var output_e *Envelope = NewEnvelope()
	var d []byte = []byte("kissa käveli kuussa")
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context, input io.Reader) {
		reader, err := NewStreamReader(output_e, input)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 5; i++ {
			output_e = NewEnvelope()
			reader.Reset(output_e, input)
			for !reader.IsFull() {
				_, err := reader.ReadWhole()
				if err != nil {
					t.Fatal(err)
				}
			}

			env := reader.GetEnvelope()
			if fmt.Sprintf("%+v", (*env.Message.Data)) != fmt.Sprintf("%+v", d) {
				t.Fatal(fmt.Errorf("incorrect data received in envelope? Received: \n%+v\nexpected: \n%+v", (*env.Message.Data), d))
			}
			//			fmt.Printf("\nenvelope:\n---\n%+v\n---\nid: %+v\nfeat: %+v\nmsg: %+v\n", reader.GetEnvelope(), string(env.EnvelopeLabels.MessageId[:]), string(env.EnvelopeLabels.MessageFeature[:]), string((*env.Message.Data)[:]))
		}
		cancel()
	}(ctx, out)
	for i := 0; i < 5; i++ {
		_, err := stream.Write(d[:])

		if err != nil {
			t.Fatal(err)
		}
		//		fmt.Printf("wrote: %d", n)
	}
	<-ctx.Done()
}

func BenchmarkReader(t *testing.B) {
	e := NewEnvelope()
	e.Message.SetData([]byte("1"))
	ior, ioo := io.Pipe()
	var in io.Writer = ioo
	var out io.Reader = ior
	var output_e *Envelope = NewEnvelope()
	d, _, _ := e.Generate()
	var okcn chan bool = make(chan bool)
	reader, err := NewStreamReader(output_e, out)
	if err != nil {
		t.Fatal(err)
	}
	t.ResetTimer()
	go func() {
		for {
			output_e = NewEnvelope()
			reader.Reset(output_e, out)
			for !reader.IsFull() {
				_, err := reader.ReadWhole()
				if err != nil {
					t.Fatal(err)
				}
			}
			okcn <- true
		}
	}()
	for i := 0; i < t.N; i++ {
		in.Write(*d)
		if err != nil {
			t.Fatal(err)
		}
		<-okcn
	}
}
