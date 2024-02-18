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

func (s *server) AddConn(conn *net.Conn) {
	id := uuid.New()
	s.connections[id] = conn
}

func (s *server) AcceptConnections(l net.Listener) {
	for {
		conn, err := l.Accept()

		if err != nil {
			log.Print(err)
		}

		s.AddConn(&conn)
	}
}

func (s server) Run() {
	l, err := net.Listen("tcp", s.address+":"+s.port)
	if err != nil {
		panic(err)
	}

	s.AcceptConnections(l)
}

func NewServer(address string, port string) *server {
	return &server{}
}
