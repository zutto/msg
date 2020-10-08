// msg is a byte serialization format for 'packets' or 'messages' or 'envelopes'.
//
// Goal of msg is to be able to store these 'packets' or 'messages' on storage mediums, be able to send them through ip links, fifo, unix pipes. (disclaimer: checksums are not.. yet.. part of the packets, so these are not to be used on unstable/unreliable mediums.)
//
// Each packet contains ID, Feature and data.
// Data can be of any length, and it shall automatically be packed into a sequence of packets that can be later recompiled into one.
package msg

import (
	"crypto/sha1"
	"errors"
)

// Envelope constants.
const (
	//init byte, not very necessary but helps to keep track of the stream.
	INIT_BYTE = 1

	//static values, may change between versions
	PREFIX_LEN          = 8
	HEADER_LEN          = 16 + 20
	ENVELOPE_HEADER_LEN = 8
	MESSAGE_PREFIX_LEN  = 2
	//CHECKSUM_LEN        = 20
	TOTAL_HEADER_LEN = PREFIX_LEN + HEADER_LEN + ENVELOPE_HEADER_LEN + MESSAGE_PREFIX_LEN // + CHECKSUM_LEN

	//default message size - Can be modified to fit use cases, such as but not limited to Tcp jumbo frames, large files, etc.
	//by default, the size is calculated to fit into TCP frame.
	SIZE = 1500 - TOTAL_HEADER_LEN - INIT_BYTE //1445

	//packet types. this is to be changed possibly, or one could make custom types? these are just hints for receiver.
	INIT = iota
	FRAGMENT
	STREAM
	ROUTE
	BROADCAST
	MULTICAST
)

//prefix: 			8 bytes
//headers 			36 bytes
//envelope headers 		16 bytes
//envelope Labels		NN bytes
//message headers (built in) 	2 bytes
//message 			NN bytes

// Envelope structure contains all parts of the envelope
type Envelope struct {
	Prefix          Prefix          // Prefix for every envelope
	Headers         Headers         // headers for the content
	EnvelopeHeaders EnvelopeHeaders // envelope headers
	EnvelopeLabels  EnvelopeLabels  //envelope labels
	Message         Message         // message

	MessageSizeLimit int  //limit of message size (for tcp frames)
	AutoIncrement    bool //automatically increment message index
	IncludeChecksum  bool
}

//NewEnvelope creates new envelope structure.
func NewEnvelope() *Envelope {
	e := Envelope{
		Prefix: Prefix{
			Init:         1,                //static 1
			Version:      1,                //to be changed?
			Length:       0,                //changes
			HeaderLength: TOTAL_HEADER_LEN, //changes maybe later?
			Type:         uint8(FRAGMENT),
		},
		Headers: Headers{
			PacketIndex:  0,
			TotalPackets: 0,
			TotalLength:  0,
		},
		EnvelopeHeaders: EnvelopeHeaders{
			MessageIdLength:      0,
			MessageFeatureLength: 0,
		},
		EnvelopeLabels: EnvelopeLabels{
			//MessageId: []byte,
			//MessageFeature: []byte,
		},
		Message: Message{
			//Message: *[]byte
		},
		MessageSizeLimit: SIZE,
		AutoIncrement:    true,
		IncludeChecksum:  true,
	}
	return &e
}

//Checksum generates checksum from all of the envelope feature checksums
func (e *Envelope) Checksum(from, to int) ([]byte, error) {
	prefixChecksum, err := e.Prefix.Checksum()
	if err != nil {
		return nil, err
	}

	headerChecksum, err := e.Headers.Checksum()
	if err != nil {
		return nil, err
	}

	envelopeHeadersChecksum, err := e.EnvelopeHeaders.Checksum()
	if err != nil {
		return nil, err
	}

	envelopeLabelsChecksum, err := e.EnvelopeLabels.Checksum()
	if err != nil {
		return nil, err
	}

	var messageChecksum []byte
	if from == 0 && to == 0 {
		messageChecksum, err = e.Message.Checksum()
	} else {
		messageChecksum, err = e.Message.ChecksumFromBytes(from, to)
	}
	if err != nil {
		return nil, err
	}
	outputChecksums := [][]byte{prefixChecksum, headerChecksum, envelopeHeadersChecksum, envelopeLabelsChecksum, messageChecksum}
	var data []byte
	for _, v := range outputChecksums {
		data = append(data, v...)
	}

	output := sha1.Sum(data)
	return output[:], nil
}

//Generate generates the envelope from zero.
//This function retuns data, length and error.
func (e *Envelope) Generate() (*[]byte, int, error) {
	return e.GenerateFromByte(0)
}

//GenerateFromByte generates envelope starting from N byte. This function retuns data, length and error.
//This is for sending large, or multipart envelopes.
//	processedBytes := 0
// 	for processedBytes < e.Envelope.Message.Len() {
//		data, n, err := s.envelope.GenerateFromByte(processed)
//		if err != nil {
//			//..handle error
//		}
//		//..send/store/whatever with the data.
//
//		processedBytes += n
//	}
func (e *Envelope) GenerateFromByte(n int) (*[]byte, int, error) {
	limit := 0
	if e.Message.Len()+e.EnvelopeLabels.Len() > e.MessageSizeLimit {
		limit = e.MessageSizeLimit - e.EnvelopeLabels.Len()
	} else {
		limit = e.Message.Len()
	}

	if limit < 1 {
		return nil, 0, errors.New("Envelope labels are longer than the message limit.")
	}

	var b []byte = make([]byte, 0)

	//message - read first to get confirmed size.. little copying and all that.
	message, read, err := e.Message.GetData(n, n+limit)
	if err != nil {
		return nil, 0, err
	}
	if read < limit {
		limit = read
	}
	e.Prefix.Length = e.Prefix.HeaderLength + uint16(e.EnvelopeLabels.Len()) + uint16(limit)
	prefix, err := e.Prefix.Generate()
	if err != nil {
		return nil, 0, err
	}
	b = append(b, (*prefix)...)

	//headers
	e.Headers.TotalLength = uint64(e.Message.Len())
	e.Headers.CalculatePacketCount(e.MessageSizeLimit)

	//envelope headers
	e.EnvelopeHeaders.SetLengths(e.EnvelopeLabels)
	envelopeHeaders, err := e.EnvelopeHeaders.Generate()
	if err != nil {
		return nil, 0, err
	}

	//envelope labels
	envelopeLabels, err := e.EnvelopeLabels.Generate()
	if err != nil {
		return nil, 0, err
	}

	if e.IncludeChecksum {
		checksum, err := e.Checksum(n, n+limit)
		if err != nil {
			return nil, 0, err
		}
		copy(e.Headers.PacketChecksum[:], checksum[:20])
	}
	headers, err := e.Headers.Generate()
	if err != nil {
		return nil, 0, err
	}

	b = append(b, (*headers)...)

	b = append(b, (*envelopeHeaders)...)

	b = append(b, (*envelopeLabels)...)

	//add message
	b = append(b, (*message)...)

	if e.AutoIncrement {
		e.Headers.IncrementIndex()
	}

	return &b, limit, nil
}

//GetHeaderSize Calculates header size of the envelope. Does not include the hidden headers included in Message used for compression.
func (e *Envelope) GetHeaderSize() int {
	return TOTAL_HEADER_LEN + len(e.EnvelopeLabels.MessageId) + len(e.EnvelopeLabels.MessageFeature)
}

//AddLabels is a convenience function to add ID and feature byte slices into the Envelope's EnvelopeLabels struct.
//you can manually insert these by doing something like this:
//
//	e.EnvelopeLabels.MessageId = []byte{"my id"}
//
func (e *Envelope) AddLabels(id, feature []byte) error {
	e.EnvelopeLabels.MessageId = id
	e.EnvelopeLabels.MessageFeature = feature

	return nil
}

//AddMessage is a convenience function to insert data into the envelopes Message.
//You can also insert data manually by doing something like this:
//	myData := []byte{"foo bar"}
//	e.Message.Data = &myData
func (e *Envelope) AddMessage(data []byte) error {
	return e.Message.SetData(data)
}

//SetCompression sets the compression method for data, set to 'NIL' for no compression.
func (e *Envelope) SetCompression(cmp uint8) error {
	e.Message.CompressionType = cmp
	return nil
}

//SetVersions sets the envelope version
//TODO - this is not currently used or validated as this is version 1.
func (e *Envelope) SetVersion(ver uint16) error {
	e.Prefix.Version = ver
	return nil
}
