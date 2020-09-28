package msg

const (
	SIZE      = 1458 //could be 1472 -- jumbo upto 65k, idc
	INIT_BYTE = 1

	PREFIX_LEN          = 8
	HEADER_LEN          = 16
	ENVELOPE_HEADER_LEN = 8
	MESSAGE_PREFIX_LEN  = 2
	TOTAL_HEADER_LEN    = PREFIX_LEN + HEADER_LEN + ENVELOPE_HEADER_LEN + MESSAGE_PREFIX_LEN

	// to be decided wtf these are really
	INIT = iota
	FRAGMENT
	STREAM
	ROUTE
	BROADCAST
	MULTICAST
)

//prefix: 			8 bytes
//headers 			16 bytes
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
	}
	return &e
}

func (e *Envelope) Generate() (*[]byte, int, error) {
	return e.GenerateFromByte(0)
}

func (e *Envelope) GenerateFromByte(n int) (*[]byte, int, error) {
	limit := 0
	if e.Message.Len() > e.MessageSizeLimit {
		limit = e.MessageSizeLimit
	} else {
		limit = e.Message.Len()
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

	//prefix
	e.Prefix.Length = e.Prefix.HeaderLength + uint16(e.EnvelopeLabels.Len()) + uint16(limit)
	prefix, err := e.Prefix.Generate()
	if err != nil {
		return nil, 0, err
	}
	b = append(b, (*prefix)...)

	//headers
	e.Headers.TotalLength = uint64(e.Message.Len())
	e.Headers.CalculatePacketCount(e.MessageSizeLimit)
	headers, err := e.Headers.Generate()
	if err != nil {
		return nil, 0, err
	}

	b = append(b, (*headers)...)

	//envelope headers
	e.EnvelopeHeaders.SetLengths(e.EnvelopeLabels)
	envelopeHeaders, err := e.EnvelopeHeaders.Generate()
	if err != nil {
		return nil, 0, err
	}
	b = append(b, (*envelopeHeaders)...)

	//envelope labels
	envelopeLabels, err := e.EnvelopeLabels.Generate()
	if err != nil {
		return nil, 0, err
	}

	b = append(b, (*envelopeLabels)...)

	//add message
	b = append(b, (*message)...)
	if e.AutoIncrement {
		e.Headers.IncrementIndex()
	}
	return &b, limit, nil
}

func (e *Envelope) GetHeaderSize() int {
	return TOTAL_HEADER_LEN + len(e.EnvelopeLabels.MessageId) + len(e.EnvelopeLabels.MessageFeature)
}

//convenience functions..
func (e *Envelope) AddLabels(id, feature []byte) error {
	e.EnvelopeLabels.MessageId = id
	e.EnvelopeLabels.MessageFeature = feature

	return nil
}

func (e *Envelope) AddMessage(data []byte) error {
	return e.Message.SetData(data)
}

func (e *Envelope) SetCompression(cmp uint8) error {
	e.Message.CompressionType = cmp
	return nil
}

func (e *Envelope) SetVersion(ver uint16) error {
	e.Prefix.Version = ver
	return nil
}
