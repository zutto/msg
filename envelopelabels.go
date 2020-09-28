package msg

import "crypto/sha1"

type EnvelopeLabels struct {
	MessageId      []byte
	MessageFeature []byte
}

func (el *EnvelopeLabels) Generate() (*[]byte, error) {
	data := append(el.MessageId, el.MessageFeature...)

	return &data, nil
}

func (el *EnvelopeLabels) GetLengths() (uint16, uint16) {
	return uint16(len(el.MessageId)), uint16(len(el.MessageFeature))
}

func (el *EnvelopeLabels) Len() int {
	return len(el.MessageId) + len(el.MessageFeature)
}

func (el *EnvelopeLabels) Parse(eh EnvelopeHeaders, data *[]byte, padding int) error {
	var read int = padding + 0
	if eh.MessageIdLength > 0 {
		el.MessageId = (*data)[read : read+int(eh.MessageIdLength)]
		read += int(eh.MessageIdLength)
	}

	if eh.MessageFeatureLength > 0 {
		el.MessageFeature = (*data)[read : read+int(eh.MessageFeatureLength)]
		read += int(eh.MessageFeatureLength)
	}

	return nil
}

func (el *EnvelopeLabels) Checksum() ([]byte, error) {
	data, err := el.Generate()
	if err != nil {
		return nil, err
	}

	output := sha1.Sum(*data)
	return output[:], nil

}
