package server

import (
	"net"
)

type server struct {
	address     string
	port        string
	connections []*net.Conn
}

func (s *server) AddConn(conn *net.Conn) {
	s.connections = append(s.connections, conn)
}

func (s server) Run() {
	l, err := net.Listen("tcp", s.address+":"+s.port)
	if err != nil {
		panic(err)
	}

	for {
		l.Accept()
	}
}

func NewServer(address string, port string) *server {
	return &server{}
}
