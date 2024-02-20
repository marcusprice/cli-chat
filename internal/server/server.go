package server

import (
	"bufio"
	"log"
	"net"

	"github.com/google/uuid"
)

type client struct {
	id   uuid.UUID
	conn net.Conn
}

type server struct {
	address string
	port    string
	clients map[uuid.UUID]client
}

func (server server) Run() {
	listener, err := net.Listen("tcp", server.address+":"+server.port)
	if err != nil {
		panic(err)
	}
	log.Printf("listening at address %v on port %v", server.address, server.port)

	server.acceptConncections(listener)
}

func (server *server) acceptConncections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			// TODO: figure out what to do for connection errors
			log.Print(err)
		} else {
			id := uuid.New()
			client := client{id: id, conn: conn}

			server.addClient(client)
			go server.handleConn(client)
		}
	}
}

func (server *server) addClient(client client) {
	server.clients[client.id] = client
}

func (server *server) removeClient(client client) {
	id := client.id.String()
	delete(server.clients, client.id)
	client.conn.Close()
	log.Printf("client id %v disconnected\n", id)
}

func (server *server) handleConn(client client) {
	reader := bufio.NewReader(client.conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				// EOF gets returned when no more input is available (in other
				// words the conn closed). If not EOF log the error
				log.Println(err)
			}

			server.removeClient(client)
			break
		} else {
			if message[0] == '\\' {
				// running command
				log.Printf("command: %v", message)
				client.conn.Write([]byte("Command executed successfully 👍\n"))
			} else {
				log.Print(message)
				server.broadcastMessage(client.id, message)
			}
		}
	}
}

func (server *server) broadcastMessage(senderId uuid.UUID, message string) {
	for id, client := range server.clients {
		if senderId != id {
			client.conn.Write([]byte(message))
		}
	}
}

func NewServer(address string, port string) *server {
	return &server{
		address: address,
		port:    port,
		clients: make(map[uuid.UUID]client),
	}
}
