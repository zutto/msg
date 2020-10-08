package msg

import (
	"errors"
	"sort"
)

type Box struct {
	envelopes []*Envelope
}

func NewBox() *Box {
	b := &Box{}
	return b
}

func (b *Box) InsertEnvelope(e *Envelope) error {
	b.envelopes = append(b.envelopes, e)
	return nil
}

func (b *Box) RecornstructMessage() (*Envelope, error) {
	if !b.IsFull() {
		return nil, errors.New("Not enough packets to reconstruct message.")
	}
	b.Sort()
	base := b.envelopes[0]
	base.Headers.TotalPackets = 0 //might be wrong?
	base.Headers.PacketIndex = 0
	for _, v := range b.envelopes {
		if err := (*v).Message.DeCompress(); err != nil {
			return nil, err
		}

		base.Message.AppendData(*(*v).Message.Data)
	}

	return base, nil
}

func (b *Box) Sort() {
	sort.Slice(b.envelopes, func(i, j int) bool {
		return b.envelopes[i].Headers.PacketIndex < b.envelopes[j].Headers.PacketIndex
	})
}

func (b *Box) IsFull() bool {
	//No messages..
	//	b.Sort()
	if len(b.envelopes) < 1 {
		return false
	}

	if len(b.envelopes) < int(b.envelopes[0].Headers.TotalPackets) {
		return false
	}
	return true
}
