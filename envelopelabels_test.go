package msg

import (
	"fmt"
	"testing"
)

var static_envlab []byte = []byte{1, 2, 3, 4}

func TestEnvelopeLabelsGenerate(t *testing.T) {
	el := EnvelopeLabels{
		MessageId:      []byte{1, 2},
		MessageFeature: []byte{3, 4},
	}
	l1, l2 := el.GetLengths()
	if l1 != 2 || l2 != 2 {
		t.Fatal(fmt.Errorf("getlengts returned incorrect lengts."))
	}

	if l := el.Len(); l != 4 {
		t.Fatal(fmt.Errorf("len returned incorrect length"))
	}

	d, err := el.Generate()
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%+v", *d) != fmt.Sprintf("%+v", static_envlab) {
		t.Fatal(fmt.Errorf("wrong data received on generate"))
	}

}

func TestEnvelopeLabelsParse(t *testing.T) {
	eh := EnvelopeHeaders{}
	eh.Parse(&static_envhed, 0)
	el := EnvelopeLabels{}
	err := el.Parse(eh, &static_envlab, 0)
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%+v", el.MessageId) != fmt.Sprintf("%+v", static_envlab[0:2]) || fmt.Sprintf("%+v", el.MessageFeature) != fmt.Sprintf("%+v", static_envlab[2:4]) {
		t.Fatal(fmt.Errorf("wrong data parsed."))
	}
}
