package main

import (
	"log"
	"net"

	"github.com/ah-naf/chat-cli-server/internals/chatserver"
)

func main() {
	chatServer := chatserver.NewChatServer()

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalln("Error opening connection:", err.Error())
	}
	defer listener.Close()
	log.Println("Server listening on port 8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err.Error())
			continue
		}

		go chatServer.HandleConnection(conn)
	}
}
