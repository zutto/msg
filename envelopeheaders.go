package msg

import (
	"crypto/sha1"
	"encoding/binary"
)

type EnvelopeHeaders struct {
	MessageIdLength      uint16 //[0:4] length of the message id
	MessageFeatureLength uint16 //[4:8] length of the message feature
}

func (eh *EnvelopeHeaders) Generate() (*[]byte, error) {
	var data []byte = make([]byte, 8)
	binary.LittleEndian.PutUint16(data[0:4], eh.MessageIdLength)
	binary.LittleEndian.PutUint16(data[4:8], eh.MessageFeatureLength)

	return &data, nil
}

func (eh *EnvelopeHeaders) SetLengths(el EnvelopeLabels) {
	eh.MessageIdLength, eh.MessageFeatureLength = el.GetLengths()
	/*	eh.MessageIdLength = uint16(len((*el).MessageId))
		eh.MessageFeatureLength = uint16(len((*el).MessageFeature))
	*/
}

func (eh *EnvelopeHeaders) Parse(data *[]byte, pad int) error {
	eh.MessageIdLength = binary.LittleEndian.Uint16((*data)[pad+0 : pad+4])
	eh.MessageFeatureLength = binary.LittleEndian.Uint16((*data)[pad+4 : pad+8])

	return nil
}

func (eh *EnvelopeHeaders) Checksum() ([]byte, error) {
	data, err := eh.Generate()
	if err != nil {
		return nil, err
	}

	output := sha1.Sum(*data)
	return output[:], nil

}
