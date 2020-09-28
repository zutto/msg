package msg

import (
	"io"
)

type StreamWriter struct {
	envelope *Envelope
	writer   io.Writer
	err      error
}

func NewStreamWriter(e *Envelope, writer io.Writer) (*StreamWriter, error) {
	var s *StreamWriter = new(StreamWriter)
	if err := s.Reset(e, writer); err != nil {
		return nil, err
	}
	return s, s.err
}

// Write implements io.reader
func (s *StreamWriter) Write(d []byte) (int, error) {
	s.envelope.Message.Data = &d
	s.envelope.Message.Compress(5)

	var processed int = 0
	var round int = 1
	//for processed < len(d) {
	for processed < len((*s.envelope.Message.Data)) {
		data, n, err := s.envelope.GenerateFromByte(processed)
		if err != nil {
			s.err = err
			return n, s.err
		}

		written := 0
		for written < len(*data) {
			wrote, err := s.writer.Write((*data)[written:])
			if err != nil {
				//.... this is nowhere near accurate.... compression fucks this up as well, idk what to do here..
				return processed, s.err
			}
			written += wrote
		}

		processed += n
		round++

	}

	return processed, s.err
}

// Flush implements io.Writer (does nothing)
func (s *StreamWriter) Flush() error {
	return s.err
}

//Reset resets the StreamWriter
func (s *StreamWriter) Reset(e *Envelope, writer io.Writer) error {
	s.envelope = e
	s.writer = writer
	s.err = nil
	return nil
}

//Close - this does nothing, does it even need to do anything?
func (s *StreamWriter) Close() error {

	return nil
}
