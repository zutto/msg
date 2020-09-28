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

var staticdata []byte = []byte{1, 5, 12, 5, 1, 4, 3, 4, 1, 55, 0, 0, 0, 0}
var staticdata_checksum []byte = []byte{78, 243, 23, 210, 217, 64, 251, 17, 233, 29, 12, 219, 182, 243, 82, 220, 204, 22, 151, 78}

func TestMessageStatic(t *testing.T) {
	m := Message{
		Data: &staticdata,
	}

	//no compression enabled
	m.Compress(0)

	ck, err := m.Checksum()
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%+v", ck) != fmt.Sprintf("%+v", staticdata_checksum) {
		t.Fatal(fmt.Sprintf("checksum mismatch, expected %+v, got %+v", ck, staticdata_checksum))
	}

	if fmt.Sprintf("%+v", staticdata) != fmt.Sprintf("%+v", *m.Data) {
		t.Errorf("data is corrupted\n%+v\n%+v\n", staticdata, *m.Data)
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
