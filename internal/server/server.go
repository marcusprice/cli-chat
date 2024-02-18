package server

import (
	"log"
	"net"

	"github.com/google/uuid"
)

type server struct {
	address     string
	port        string
	connections map[uuid.UUID]*net.Conn
}

func (s *server) addConn(conn *net.Conn) {
	id := uuid.New()
	s.connections[id] = conn
}

func (s *server) acceptConncections(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			// TODO: figure out what to do for connection errors
			log.Print(err)
		} else {
			s.addConn(&conn)
		}
	}
}

func (s server) Run() {
	l, err := net.Listen("tcp", s.address+":"+s.port)
	if err != nil {
		panic(err)
	}
	log.Printf("listening at address %v on port %v", s.address, s.port)

	s.acceptConncections(l)
}

func NewServer(address string, port string) *server {
	return &server{address: address, port: port}
}
