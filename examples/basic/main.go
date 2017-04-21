package main

import (
	"log"
	"net"

	mux "github.com/nimona/go-nimona-mux"
)

func main() {
	a, b := net.Pipe()

	mpa, _ := mux.New(a)
	mpb, _ := mux.New(b)

	mes := []byte("Hello world")
	go func() {
		s, err := mpb.Accept()
		if err != nil {
			log.Println(err)
		}

		_, err = s.Write(mes)
		if err != nil {
			log.Println(err)
		}

		err = s.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	s, err := mpa.NewStream()
	if err != nil {
		log.Println(err)
	}

	buf := make([]byte, len(mes))
	n, err := s.Read(buf)
	if err != nil {
		log.Println(err)
	}

	if n != len(mes) {
		log.Println("read wrong amount")
	}

	if string(buf) != string(mes) {
		log.Println("got bad data")
	}

	s.Close()

	// mpa.Close()
	// mpb.Close()
}
