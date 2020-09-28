package msg

import (
	//	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"
)

func BenchmarkWriterWrite(b *testing.B) {
	e := NewEnvelope()
	//	e.Message.CompressionType = LZ4
	ior, ioo := io.Pipe()
	var in io.Writer = ioo
	var out io.Reader = ior
	stream, err := NewStreamWriter(e, in)
	if err != nil {
		b.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context, input io.Reader) {
		var read []byte = make([]byte, 1000)
		var totalread = 0
		for {

			select {
			case <-ctx.Done():
				//fmt.Printf("read bytes: %d\n", totalread)
				return
			default:
			}

			n, _ := input.Read(read)
			totalread += n
		}
	}(ctx, out)
	written := 0
	var data [1]byte
	for i := 0; i < b.N; i++ {
		n, _ := stream.Write(data[:])
		written += n
	}

	//fmt.Printf("wrote %d amount of data\n", written)
	cancel()
	ioo.Close()

}

func TestWriter(t *testing.T) {
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

	var d []byte = []byte("kissa käveli kuussa")
	ctx, _ := context.WithTimeout(context.Background(), 200*time.Millisecond)
	go func(ctx context.Context, input io.Reader) {
		var pass int = 0
		for {
			select {
			case <-ctx.Done():
				//fmt.Printf("timeout\n")
				return
			default:
			}

			var prefix []byte = make([]byte, 8)
			var prefixRead int = 0
			for prefixRead < PREFIX_LEN {
				n, err := input.Read(prefix)
				if err != nil {
					t.Fatal(err)
				}

				prefixRead += n
			}
			p := Prefix{}
			err := p.Parse(&prefix, 0)
			if err != nil {
				t.Fatal(err)
			}
			var readToEnd []byte = make([]byte, p.Length-8)
			var er int = 0
			for er < int(p.Length-8) {
				n, err := input.Read(readToEnd)
				if err != nil {
					t.Fatal(err)
				}
				er += n
			}

			//fmt.Printf("prefix: %+v\n", p)

			h := Headers{}
			h.Parse(&readToEnd, 0)
			//fmt.Printf("head: %+v\n", h)

			eh := EnvelopeHeaders{}
			eh.Parse(&readToEnd, HEADER_LEN)
			//fmt.Printf("envhed: %+v\n", eh)

			el := EnvelopeLabels{}
			el.Parse(eh, &readToEnd, HEADER_LEN+ENVELOPE_HEADER_LEN)
			//fmt.Printf("labels: %s\n%s\n", string(el.MessageId[:]), string(el.MessageFeature[:]))
			if fmt.Sprintf("%v", el.MessageId) != fmt.Sprintf("%+v", e.EnvelopeLabels.MessageId) {
				t.Fatal(fmt.Sprintf("label message id did not match!"))
			}
			if fmt.Sprintf("%v", el.MessageFeature) != fmt.Sprintf("%+v", e.EnvelopeLabels.MessageFeature) {
				t.Fatal(fmt.Sprintf("label feature did not match!"))
			}
			m := Message{}
			m.Parse(&readToEnd, HEADER_LEN+ENVELOPE_HEADER_LEN+el.Len())
			m.DeCompress()
			//	fmt.Printf("message: %+v\n%+v\n%+v\n%+v/%+v -- %+v\n", m, string((*m.Data)), readToEnd[HEADER_LEN+ENVELOPE_HEADER_LEN+el.Len():], el.Len(), HEADER_LEN+ENVELOPE_HEADER_LEN+el.Len(), len(*m.Data))
			pass++
			if fmt.Sprintf("%+v", *m.Data) != fmt.Sprintf("%+v", d) {
				t.Fatal(fmt.Sprintf("content did not match!\n%+v\n%+v\n", *m.Data, d))
			}
		}
	}(ctx, out)
	_, err = stream.Write(d[:])

	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("wrote: %d", n)
	<-ctx.Done()
}
