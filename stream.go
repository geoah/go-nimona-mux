package mux

import "time"

// Stream -
type Stream struct {
	id        uint64
	initiator bool
	dataIn    chan []byte
	m         *Mux

	extra []byte

	wDeadline time.Time
	rDeadline time.Time
}

// Write implements io.Writer
func (s *Stream) Write(b []byte) (n int, err error) {
	cmd := Receiver
	if s.initiator {
		cmd = Initiator
	}
	header := (s.id << 3) | cmd
	return s.m.sendMsg(header, b, s.wDeadline)
}

// Close implements io.Closer
func (s *Stream) Close() error {
	header := (s.id << 3) | CloseLocal // TODO Check if this is the correct cmd to send
	_, err := s.m.sendMsg(header, []byte{}, s.wDeadline)
	return err
}

// Read implements io.Reader
func (s *Stream) Read(bs []byte) (int, error) {
	l := 0

	buff := s.extra
	s.extra = []byte{}
	if len(buff) == 0 {
		buff = <-s.dataIn
	}

	r := len(buff)
	if r > len(bs) {
		r = len(bs)
	}

	copy(bs[l:r], buff[l:r])

	if len(buff) >= len(bs) {
		s.extra = buff[r:]
	}

	return r, nil
}
