package msg

import "crypto/sha1"

//EnvelopeLabels contain optional user defined ID and Feature for the packet.
//TODO - limits for the sizes to ensure that the labels cannot be larger than the packet max size itself.
type EnvelopeLabels struct {
	MessageId      []byte
	MessageFeature []byte
}

//Generate generates the EnvelopeLables for the envelope
func (el *EnvelopeLabels) Generate() (*[]byte, error) {
	data := append(el.MessageId, el.MessageFeature...)

	return &data, nil
}

//GetLengths returns the lengths of MessageId and Messagefeature.
func (el *EnvelopeLabels) GetLengths() (uint16, uint16) {
	return uint16(len(el.MessageId)), uint16(len(el.MessageFeature))
}

//Len returns combined length of the MessageId and MessageFeature.
func (el *EnvelopeLabels) Len() int {
	return len(el.MessageId) + len(el.MessageFeature)
}

//Parse parses input data, data can be read from padded starting position.
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

//Checksum returns SHA-1 checksum for the EnvelopeLabels
func (el *EnvelopeLabels) Checksum() ([]byte, error) {
	data, err := el.Generate()
	if err != nil {
		return nil, err
	}

	output := sha1.Sum(*data)
	return output[:], nil

}
