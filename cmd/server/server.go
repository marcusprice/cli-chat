package main

import (
	"cmp"
	"flag"
	"os"

	"github.com/marcusprice/cli-chat/internal/server"
)

func main() {
	addressEnv := os.Getenv("ADDRESS")
	portEnv := os.Getenv("PORT")
	defaultAddress := cmp.Or(addressEnv, "127.0.0.1")
	defaultPort := cmp.Or(portEnv, "42069")

	address := flag.String("address", defaultAddress, "server address")
	port := flag.String("port", defaultPort, "server port")
	flag.Parse()

	server := server.NewServer(*address, *port)
	server.Run()
}
