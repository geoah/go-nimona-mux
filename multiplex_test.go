package mux

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func handle(m *Mux) {
	for {
		mss, _ := m.Accept()
		fmt.Println("on=", mss.id, " *** Accepted connection")
		go func() {
			mssr := bufio.NewReader(mss)
			for {
				line, err := mssr.ReadString('\n')
				fmt.Println("on=", mss.id, "msg=", strings.TrimSpace(string(line)), "len=", len(line), "err=", err)
			}
		}()
	}
}

func TestBasicStreams(t *testing.T) {
	s, _ := net.Listen("tcp", "127.0.0.1:12701")
	c, _ := net.Dial("tcp", "127.0.0.1:12701")

	ss, _ := s.Accept()
	ms, _ := New(ss)
	go handle(ms)

	// let the dust settle a bit
	time.Sleep(time.Second * 1)

	// create a new mux on the client
	mc, _ := New(c)
	// handle connections to the client
	go handle(mc)

	// open a stream from the client to the server
	mcs, _ := mc.NewStream()
	mcs.Write([]byte("hello 1 on 1 with a very padded thing\n"))

	// open another stream from the client to the server
	mcs2, _ := mc.NewStream()
	mcs2.Write([]byte("hello 1 on 2\n"))

	mcs.Write([]byte("hello 2 on 1\n"))
	mcs.Write([]byte("hello 3 on 1\n"))
	mcs.Write([]byte("hello 4 on 1\n"))
	mcs.Write([]byte("hello 5 on 1\n"))
	mcs.Write([]byte("hello 6 on 1\n"))

	// other way around now
	// open a stream from the server to the client
	mss, _ := ms.NewStream()
	mss.Write([]byte("omfg this fucking worked\n"))

	mcs2.Write([]byte("hello 2 on 2\n"))
	mcs2.Write([]byte("hello 3 on 2\n"))
	mcs2.Write([]byte("hello 4 on 2\n"))
	mcs2.Write([]byte("hello 5 on 2\n"))
	mcs2.Write([]byte("hello 6 on 2\n"))

	for {
		// wait
	}
}
