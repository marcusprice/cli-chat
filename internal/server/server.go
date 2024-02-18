package server

import (
	"bufio"
	"log"
	"net"

	"github.com/google/uuid"
)

type server struct {
	address string
	port    string
	clients map[uuid.UUID]client
}

type client struct {
	id   uuid.UUID
	conn net.Conn
}

func (s *server) addClient(c client) {
	s.clients[c.id] = c
}

func (s *server) acceptConncections(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			// TODO: figure out what to do for connection errors
			log.Print(err)
		} else {
			id := uuid.New()
			client := client{id: id, conn: conn}

			s.addClient(client)
			go s.handleConn(client)
		}
	}
}

func (s *server) broadcastMessage(senderId uuid.UUID, message string) {
	for id, client := range s.clients {
		if senderId != id {
			client.conn.Write([]byte(message))
		}
	}
}

func (s *server) handleConn(c client) {
	reader := bufio.NewReader(c.conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				// EOF gets returned when no more input is available (in other
				// words the conn closed). If not EOF log the error
				log.Println(err)
			}

			delete(s.clients, c.id)
			c.conn.Close()
			break
		} else {
			log.Println(message)
			s.broadcastMessage(c.id, message)
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
	return &server{
		address: address,
		port:    port,
		clients: make(map[uuid.UUID]client),
	}
}
