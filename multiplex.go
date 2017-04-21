package mux

import (
	"bufio"
	"encoding/binary"
	"io"
	"time"
)

const (
	// ProtocolID for multiselect
	ProtocolID = "/nimux/v1.0.0"
)

const (
	NewStream  uint64 = 0
	Receiver   uint64 = 1
	Initiator  uint64 = 2
	CloseLocal uint64 = 3
	Close      uint64 = 4
)

type Mux struct {
	con      io.ReadWriteCloser
	buf      *bufio.Reader
	nextID   uint64
	acceptCh chan *Stream
	streams  map[uint64]*Stream
}

func New(con io.ReadWriteCloser) (*Mux, error) {
	mux := &Mux{
		con:      con,
		buf:      bufio.NewReader(con),
		acceptCh: make(chan *Stream, 10),
		streams:  map[uint64]*Stream{},
	}

	go mux.handleIncoming()

	return mux, nil
}

func (m *Mux) Accept() (*Stream, error) {
	select {
	case stream := <-m.acceptCh:
		return stream, nil
	}
}

func (m *Mux) handleIncoming() {
	defer m.con.Close()
	for {
		ch, cmd, err := m.readHeader()
		if err != nil {
			// TODO Handle error
			return
		}

		b, err := m.readNext()
		if err != nil {
			// TODO Handle error
			return
		}

		switch cmd {
		case NewStream:
			// bump next id
			m.nextID++
			stream := &Stream{
				id:        ch,
				m:         m,
				dataIn:    make(chan []byte, 8),
				initiator: false,
			}
			m.streams[ch] = stream
			m.acceptCh <- stream
		// case Receiver:
		case Receiver, Initiator:
			stream := m.streams[ch]
			for i := 0; i < len(b); i = i + 8 {
				to := i + 8
				if to > len(b) {
					to = len(b)
				}
				stream.dataIn <- b[i:to]
			}
		case CloseLocal:

		case Close:
		}
	}
}

func (m *Mux) readHeader() (uint64, uint64, error) {
	h, err := binary.ReadUvarint(m.buf)
	if err != nil {
		return 0, 0, err
	}
	ch := h >> 3
	cmd := h & 7
	return ch, cmd, nil
}

func (m *Mux) readNext() ([]byte, error) {
	// get length
	l, err := binary.ReadUvarint(m.buf)
	if err != nil {
		return nil, err
	}

	if l == 0 {
		return nil, nil
	}

	buff := make([]byte, l)
	m.buf.Read(buff)

	return buff, nil
}

func (m *Mux) NewStream() (*Stream, error) {
	// bump next id
	m.nextID++
	ch := m.nextID
	h := (ch << 3) | NewStream

	s := &Stream{
		id:        ch,
		dataIn:    make(chan []byte, 8),
		m:         m,
		initiator: true,
	}

	m.streams[ch] = s

	_, err := m.sendMsg(h, []byte{}, time.Time{})
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *Mux) sendMsg(header uint64, data []byte, dl time.Time) (int, error) {
	hdrBuf := make([]byte, 20)
	n := binary.PutUvarint(hdrBuf, header)
	n += binary.PutUvarint(hdrBuf[n:], uint64(len(data)))
	_, err := m.con.Write(hdrBuf[:n])
	if err != nil {
		return 0, err
	}

	if len(data) != 0 {
		_, err = m.con.Write(data)
		if err != nil {
			return 0, err
		}
	}

	return len(data), nil
}
