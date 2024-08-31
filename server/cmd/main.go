package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	globalClients = make(map[net.Conn]string)
	usernames = make(map[string]bool)
)

func main() {
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	var username string
	for {
		// Get the username
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)

		if _, exists := usernames[username]; exists {
			conn.Write([]byte("Username already taken. Please enter a different username:\n"))
		} else {
			usernames[username] = true
			conn.Write([]byte(fmt.Sprintf("Your username is %s\n", username)))
			break
		}
	}

	// Get the chat mode
	mode, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Client disconnected.")
		return
	}

	mode = strings.TrimSpace(mode)
	log.Println("User", username, "entered mode:", mode)
	switch mode {
	case "global":
		handleGlobalChat(conn, reader, username)
	default:
		conn.Write([]byte("Invalid mode. Type 'global' or 'group'.\n"))
	}
}

func handleGlobalChat(conn net.Conn, reader *bufio.Reader, username string) {
	globalClients[conn] = username
	// log.Println(username, "joined global chat")
	broadcastMessage(username + " joined the chat", conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println(username, "disconnected")
			broadcastMessage(username + " left the chat", conn)
			removeClient(conn, username)
			return
		}
		formattedMessage := fmt.Sprintf("%s > %s", username, strings.TrimSpace(message))
		broadcastMessage(formattedMessage, conn)
	}
}

func broadcastMessage(message string, sender net.Conn) {
	for client, username := range globalClients {
		if client != sender {
			_, err := client.Write([]byte(message + "\n"))
			if err != nil {
				log.Println("Error broadcasting message to", username, ":", err.Error())
			}
		}
	}
}

func removeClient(conn net.Conn, username string) {
	delete(globalClients, conn)
	delete(usernames, username)
}
