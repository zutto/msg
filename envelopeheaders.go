package msg

import (
	"crypto/sha1"
	"encoding/binary"
)

//EnvelopeHeaders contain the headers for the message header, including ID and Feature labels lengths.
type EnvelopeHeaders struct {
	MessageIdLength      uint16 //[0:4] length of the message id
	MessageFeatureLength uint16 //[4:8] length of the message feature
}

//Generate generates the EnvelopeHeaders for the envelope.
func (eh *EnvelopeHeaders) Generate() (*[]byte, error) {
	var data []byte = make([]byte, 8)
	binary.LittleEndian.PutUint16(data[0:4], eh.MessageIdLength)
	binary.LittleEndian.PutUint16(data[4:8], eh.MessageFeatureLength)

	return &data, nil
}

//SetLengths sets the values for the EnvelopeHeaders from EnvelopeLabels.
func (eh *EnvelopeHeaders) SetLengths(el EnvelopeLabels) {
	eh.MessageIdLength, eh.MessageFeatureLength = el.GetLengths()
	/*	eh.MessageIdLength = uint16(len((*el).MessageId))
		eh.MessageFeatureLength = uint16(len((*el).MessageFeature))
	*/
}

//Parse parses input data, data can be read from padded starting position.
func (eh *EnvelopeHeaders) Parse(data *[]byte, pad int) error {
	eh.MessageIdLength = binary.LittleEndian.Uint16((*data)[pad+0 : pad+4])
	eh.MessageFeatureLength = binary.LittleEndian.Uint16((*data)[pad+4 : pad+8])

	return nil
}

//Checksum generates SHA-1 checksum from the EnvelopeHeaders for envelope.
func (eh *EnvelopeHeaders) Checksum() ([]byte, error) {
	data, err := eh.Generate()
	if err != nil {
		return nil, err
	}

	output := sha1.Sum(*data)
	return output[:], nil

}
