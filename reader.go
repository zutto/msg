package msg

import (
	"context"
	"fmt"
	"io"
)

//
const (
	no_state = iota
	prefix_state
	headers_state
	envelopeheaders_state
	envelopelabels_state
	message_state
)

type envelopeReadState struct {
	Prefix          bool //1
	Headers         bool //2
	EnvelopeHeaders bool //3
	EnvelopeLabels  bool //4
	Message         bool //5

	STATE    int
	READ_LEN int
}

type StreamReader struct {
	read_state envelopeReadState
	reader     io.Reader
	envelope   *Envelope
	err        error

	ctx       context.Context
	ctxCancel context.CancelFunc
	buffer    []byte
}

func NewStreamReader(e *Envelope, r io.Reader) (*StreamReader, error) {
	var s *StreamReader = new(StreamReader)
	if err := s.Reset(e, r); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *StreamReader) GetState() *envelopeReadState {
	return &s.read_state
}

func (s *StreamReader) SetState(es *envelopeReadState) {
	s.read_state = *es
}

func (s *StreamReader) IsFull() bool {
	if s.read_state.Prefix && s.read_state.Headers && s.read_state.EnvelopeHeaders && s.read_state.EnvelopeLabels && s.read_state.Message {
		return true
	}
	return false
}

func (s *StreamReader) GetEnvelope() *Envelope {
	return s.envelope
}

func (s *StreamReader) readUntilFilled(cont *[]byte) (int, error) {
	var read int = 0
	for read < len(*cont) {
		select {
		case <-s.ctx.Done():
			return read, s.ctx.Err()
		}
		r, err := s.reader.Read(*cont)
		if err != nil {
			s.err = err
			return read + r, s.err
		}

		read += r
	}

	return read, nil
}

//todo FIX does not implement reader{} .. wtf was I thinking
func (s *StreamReader) Read(data []byte) (int, error) {
	var read int = 0
	if s.read_state.Prefix && s.read_state.Headers && s.read_state.EnvelopeHeaders && s.read_state.EnvelopeLabels && s.read_state.Message {
		return 0, fmt.Errorf("envelope is fully read, reset to read more.")
	}

	//prefix
	if s.envelope.Prefix.Init == uint8(0) || !s.read_state.Prefix {
		if s.read_state.STATE != prefix_state {
			s.buffer = make([]byte, PREFIX_LEN)
			s.read_state.STATE = prefix_state
			s.read_state.READ_LEN = 0
		}
		if s.read_state.READ_LEN < PREFIX_LEN {
			for s.read_state.READ_LEN < PREFIX_LEN {
				r, err := s.readUntilFilled(&s.buffer)
				if err != nil {
					s.err = err
					return read + r, s.err
				}
				read += r
				s.read_state.READ_LEN += r
			}
		}
		s.err = s.not_readPrefix(&s.buffer)
		if s.err != nil {
			return read, s.err
		}

	}

	if s.read_state.Prefix && !s.read_state.Headers {
		if s.read_state.STATE != headers_state {
			s.buffer = make([]byte, HEADER_LEN)
			s.read_state.STATE = headers_state
			s.read_state.READ_LEN = 0
		}
		if s.read_state.READ_LEN < HEADER_LEN {
			for s.read_state.READ_LEN < HEADER_LEN {
				r, err := s.readUntilFilled(&s.buffer)
				if err != nil {
					s.err = err
					return read + r, s.err
				}
				read += r
				s.read_state.READ_LEN += r
			}
		}

	}
	return 0, nil
}

func (s *StreamReader) ReadWhole() (int, error) {
	var read int = 0
	if s.read_state.Prefix && s.read_state.Headers && s.read_state.EnvelopeHeaders && s.read_state.EnvelopeLabels && s.read_state.Message {
		return 0, fmt.Errorf("envelope is fully read, reset to read more.")
	}
	//test init
	if s.envelope.Prefix.Init == uint8(0) || !s.read_state.Prefix {
		//read init
		r, err := s.readPrefix()
		if err != nil {
			s.err = err
			return r + read, s.err
		}
		read += r
	}

	//headers
	if !s.read_state.Headers {
		r, err := s.readHeaders()
		if err != nil {
			s.err = err
			return r + read, s.err
		}
		read += r
	}

	//envelope label headers
	if !s.read_state.EnvelopeHeaders {
		r, err := s.readEnvelopeHeaders()
		if err != nil {
			s.err = err
			return r + read, s.err
		}
		read += r
	}

	//envelope labels
	if !s.read_state.EnvelopeLabels {
		r, err := s.readEnvelopeLabels()
		if err != nil {
			return r + read, s.err
		}
		read += r
	}

	//message
	if !s.read_state.Message {
		r, err := s.readMessage()
		if err != nil {
			return r + read, s.err
		}
		read += r
	}

	return read, nil
}

func (s *StreamReader) Reset(e *Envelope, r io.Reader) error {
	s.envelope = e
	s.reader = r
	s.err = nil

	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	s.read_state = envelopeReadState{}
	return nil
}

func (s *StreamReader) Close() error {
	s.ctxCancel()
	return nil
}

func (s *StreamReader) readFull(cont *[]byte) (int, error) {
	var read int = 0
	for read < len(*cont) {
		r, err := s.reader.Read(*cont)
		if err != nil {
			s.err = err
			return read + r, s.err
		}

		read += r
	}

	return read, nil
}

func (s *StreamReader) not_readPrefix(data *[]byte) error {
	if len(*data) < 8 {
		return fmt.Errorf("input data is not full")
	}
	// init byte handling
	if uint8((*data)[0]) != uint8(INIT_BYTE) {
		s.err = fmt.Errorf("INIT BYTE MISMATCH, packet loss? Received: %+v, expected: %+v", uint8((*data)[0]), uint8(INIT_BYTE))
		return s.err
	}

	s.err = s.envelope.Prefix.Parse(data, 0)
	if s.err != nil {
		return s.err
	}

	s.read_state.Prefix = true
	return nil
}

func (s *StreamReader) readPrefix() (n int, err error) {
	// init byte handling
	var initByte []byte = make([]byte, 1)
	n, s.err = s.reader.Read(initByte)
	if uint8(initByte[0]) != uint8(INIT_BYTE) {
		s.err = fmt.Errorf("INIT BYTE MISMATCH, packet loss? Received: %+v, expected: %+v", uint8(initByte[0]), uint8(INIT_BYTE))
		return n, s.err
	}

	//read rest of init packet
	var ilPrefix []byte = make([]byte, PREFIX_LEN-1)
	r, err := s.readFull(&ilPrefix)
	if err != nil {
		s.err = err
		return n + r, s.err
	}
	var full []byte = append(initByte, ilPrefix[:]...)
	s.err = s.envelope.Prefix.Parse(&full, 0)
	if s.err != nil {
		return n + r, s.err
	}

	s.read_state.Prefix = true
	return n + r, nil
}
func (s *StreamReader) not_readHeaders(data *[]byte) error {
	if len(*data) < HEADER_LEN {
		return fmt.Errorf("input data is not full")
	}

	s.err = s.envelope.Headers.Parse(data, 0)
	if s.err != nil {
		return s.err
	}

	s.read_state.Headers = true
	return nil
}

func (s *StreamReader) readHeaders() (int, error) {
	//cont todo
	var header []byte = make([]byte, HEADER_LEN)
	r, err := s.readFull(&header)
	if err != nil {
		s.err = err
		return r, err
	}

	s.err = s.envelope.Headers.Parse(&header, 0)
	if s.err != nil {
		return r, s.err
	}

	s.read_state.Headers = true
	return r, nil
}

func (s *StreamReader) readEnvelopeHeaders() (int, error) {
	var header []byte = make([]byte, ENVELOPE_HEADER_LEN)
	r, err := s.readFull(&header)
	if err != nil {
		s.err = err
		return r, err
	}

	s.err = s.envelope.EnvelopeHeaders.Parse(&header, 0)
	if s.err != nil {
		return r, s.err
	}

	s.read_state.EnvelopeHeaders = true
	return r, nil
}

func (s *StreamReader) readEnvelopeLabels() (int, error) {
	var labels []byte = make([]byte, s.envelope.EnvelopeHeaders.MessageIdLength+s.envelope.EnvelopeHeaders.MessageFeatureLength)
	r, err := s.readFull(&labels)
	if err != nil {
		s.err = err
		return r, s.err
	}
	s.err = s.envelope.EnvelopeLabels.Parse(s.envelope.EnvelopeHeaders, &labels, 0)
	if s.err != nil {
		return r, s.err
	}

	s.read_state.EnvelopeLabels = true
	return r, nil
}

// e.Prefix.Length = e.Prefix.HeaderLength + uint16(e.EnvelopeLabels.Len()) + uint16(limit)
func (s *StreamReader) readMessage() (int, error) {
	//+2 = compressed/compressiontype bytes
	var message []byte = make([]byte, int(s.envelope.Prefix.Length)-int(s.envelope.Prefix.HeaderLength)-s.envelope.EnvelopeLabels.Len()+2)
	r, err := s.readFull(&message)
	if err != nil {
		s.err = err
		return r, s.err
	}
	//fmt.Printf("\nfull msg: %+v", string(message[:]))
	s.err = s.envelope.Message.Parse(&message, 0)
	if s.err != nil {
		return r, s.err
	}

	s.read_state.Message = true
	return r, nil
}
