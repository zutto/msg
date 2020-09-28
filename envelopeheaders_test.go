package msg

import (
	"fmt"
	"testing"
)

var static_envhed []byte = []byte{2, 0, 0, 0, 2, 0, 0, 0}
var envhed_checksum []byte = []byte{248, 135, 214, 21, 61, 141, 104, 19, 72, 239, 194, 168, 165, 45, 209, 213, 252, 88, 20, 18}

func TestEnvHeadGenerate(t *testing.T) {
	el := EnvelopeLabels{
		MessageId:      []byte{1, 2},
		MessageFeature: []byte{1, 2},
	}
	eh := EnvelopeHeaders{}
	eh.SetLengths(el)

	d, err := eh.Generate()
	if err != nil {
		t.Fatal(err)
	}

	ck, err := eh.Checksum()
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%+v", ck) != fmt.Sprintf("%+v", envhed_checksum) {
		t.Fatal(fmt.Sprintf("checksum mismatch, got %+v, expected %+v", ck, envhed_checksum))
	}

	if fmt.Sprintf("%+v", *d) != fmt.Sprintf("%+v", static_envhed) {
		t.Fatal(fmt.Errorf("got incorrect data.\nexpected: %+v\nReceived: %+v", static_envhed, *d))
	}

}

func TestEnvHeadParse(t *testing.T) {
	eh := EnvelopeHeaders{}
	err := eh.Parse(&static_envhed, 0)
	if err != nil {
		t.Fatal(err)
	}

	if eh.MessageIdLength != uint16(2) || eh.MessageFeatureLength != uint16(2) {
		t.Fatal(fmt.Errorf("parsed data incorrect"))
	}
}
