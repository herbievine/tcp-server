package main

import (
	"fmt"
	"net"
)

type Message struct {
	from  net.Addr
	bytes []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	msgch      chan Message
}

func newServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	fmt.Println("Listening on", s.listenAddr)

	defer ln.Close()
	defer close(s.msgch)

	go s.accept()

	s.ln = ln

	return nil
}

func (s *Server) accept() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			continue
		}

		go s.read(conn)
	}
}

func (s *Server) read(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			continue
		}

		s.msgch <- Message{
			from:  conn.RemoteAddr(),
			bytes: buf[:n],
		}
	}
}

func main() {
	server := newServer(":6969")

	err := server.start()
	if err != nil {
		fmt.Println(err)
		return
	}

	for msg := range server.msgch {
		fmt.Printf("Message received from (%s): %s\n", msg.from, string(msg.bytes))
	}
}
