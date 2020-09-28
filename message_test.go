package msg

import (
	"fmt"
	"github.com/google/gofuzz"
	"testing"
)

func TestMessage(t *testing.T) {
	var data []byte
	f := fuzz.New()
	f.Fuzz(&data)
	m := Message{
		Data: &data,
	}

	//no compression enabled
	m.Compress(0)

	if fmt.Sprintf("%+v", data) != fmt.Sprintf("%+v", *m.Data) {
		t.Errorf("data is corrupted\n%+v\n%+v\n", data, *m.Data)
	}
}

func BenchmarkMessageLz(b *testing.B) {
	var data []byte
	f := fuzz.New().NilChance(0)
	f.Fuzz(&data)
	m := Message{
		Data:            &data,
		CompressionType: LZ4,
	}

	var err error
	for i := 0; i < b.N; i++ {
		err = m.Compress(5)
		if err != nil {
			b.Fatal(err)
		}
		err = m.DeCompress()
		if err != nil {
			b.Fatal(err)
		}
		if fmt.Sprintf("%+v", data) != fmt.Sprintf("%+v", *m.Data) {
			b.Fatal("error, data is corrupted.")
		}
	}

}

func BenchmarkMessageGz(b *testing.B) {
	var data []byte
	f := fuzz.New().NilChance(0)
	f.Fuzz(&data)
	m := Message{
		Data:            &data,
		CompressionType: GZIP,
	}

	var err error
	for i := 0; i < b.N; i++ {
		err = m.Compress(5)
		if err != nil {
			b.Fatal(err)
		}
		err = m.DeCompress()
		if err != nil {
			b.Fatal(err)
		}
		if fmt.Sprintf("%+v", data) != fmt.Sprintf("%+v", *m.Data) {
			b.Fatal("error, data is corrupted.")
		}
	}

}

func TestMessageGz(t *testing.T) {
	var data []byte
	f := fuzz.New().NilChance(0)
	f.Fuzz(&data)
	m := Message{
		Data:            &data,
		CompressionType: GZIP,
	}

	//no compression enabled
	m.Compress(5)

	if fmt.Sprintf("%+v", data) == fmt.Sprintf("%+v", *m.Data) {
		t.Errorf("data is not compressed?")
	}

	m.DeCompress()
	if fmt.Sprintf("%+v", data) != fmt.Sprintf("%+v", *m.Data) {
		t.Errorf("data is corrupted\n%+v\n%+v\n", data, *m.Data)
	}
}

func TestMessageLz(t *testing.T) {
	var data []byte
	f := fuzz.New().NilChance(0)
	f.Fuzz(&data)
	m := Message{
		Data:            &data,
		CompressionType: LZ4,
	}

	//no compression enabled
	m.Compress(5)

	if fmt.Sprintf("%+v", data) == fmt.Sprintf("%+v", *m.Data) {
		t.Errorf("data is not compressed?")
	}

	m.DeCompress()
	if fmt.Sprintf("%+v", data) != fmt.Sprintf("%+v", *m.Data) {
		t.Errorf("data is corrupted\n%+v\n%+v\n", data, *m.Data)
	}
}
