package main

import (
	"log"
	"net"

	"bufio"

	mux "github.com/nimona/go-nimona-mux"
)

func main() {
	a, b := net.Pipe()

	mpa, _ := mux.New(a)
	mpb, _ := mux.New(b)

	mes := []byte("Hello world!\n")
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

	r := bufio.NewReader(s)
	line, _, err := r.ReadLine()
	if err != nil {
		log.Println(err)
	}

	if len(line) != len(mes)-1 {
		log.Println("read wrong amount", len(line), len(mes))
	}

	if string(line)+"\n" != string(mes) {
		log.Println("got bad data", string(line), string(mes))
	}

	s.Close()

	// mpa.Close()
	// mpb.Close()
}
