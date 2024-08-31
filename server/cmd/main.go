package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// ChatServer struct encapsulates the shared resources and the mutex
type ChatServer struct {
	mu            sync.Mutex
	globalClients map[net.Conn]string
	usernames     map[string]bool
}

// NewChatServer creates a new instance of ChatServer
func NewChatServer() *ChatServer {
	return &ChatServer{
		globalClients: make(map[net.Conn]string),
		usernames:     make(map[string]bool),
	}
}

// AddClient safely adds a client to the globalClients map
func (cs *ChatServer) AddClient(conn net.Conn, username string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.globalClients[conn] = username
	cs.usernames[username] = true
}

// RemoveClient safely removes a client from the globalClients map
func (cs *ChatServer) RemoveClient(conn net.Conn, username string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	delete(cs.globalClients, conn)
	delete(cs.usernames, username)
}

// BroadcastMessage safely sends a message to all connected clients except the sender
func (cs *ChatServer) BroadcastMessage(message string, sender net.Conn) {
	for client, username := range cs.globalClients {
		if client != sender {
			cs.mu.Lock()
			_, err := client.Write([]byte(message + "\n"))
			cs.mu.Unlock()
			if err != nil {
				log.Println("Error broadcasting message to", username, ":", err.Error())
			}
		}
	}
}

// HandleConnection handles the connection for each client
func (cs *ChatServer) HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	var username string
	for {
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
		log.Println(username)

		if _, exists := cs.usernames[username]; exists {
			conn.Write([]byte("Username already taken. Please enter a different username:\n"))
		} else {
			cs.AddClient(conn, username)
			conn.Write([]byte(fmt.Sprintf("Your username is %s\n", username)))
			break
		}
	}

	mode, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Client disconnected.")
		cs.RemoveClient(conn, username)
		return
	}

	mode = strings.TrimSpace(mode)
	log.Println("User", username, "entered mode:", mode)
	switch mode {
	case "global":
		cs.handleGlobalChat(conn, reader, username)
	default:
		conn.Write([]byte("Invalid mode. Type 'global' or 'group'.\n"))
	}
}

// handleGlobalChat manages the global chat for the client
func (cs *ChatServer) handleGlobalChat(conn net.Conn, reader *bufio.Reader, username string) {
	coolJoinMessage := fmt.Sprintf("\033[1;36m%s has entered the chat! 🎉\033[0m", username)
	cs.BroadcastMessage(coolJoinMessage, conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println(username, "disconnected")
			cs.BroadcastMessage(fmt.Sprintf("\033[1;31m%s left the chat!\033[0m", username), conn)
			cs.RemoveClient(conn, username)
			return
		}
		formattedMessage := fmt.Sprintf("\033[1;32m%s > %s\033[0m", username, strings.TrimSpace(message))
		cs.BroadcastMessage(formattedMessage, conn)
	}
}

func main() {
	chatServer := NewChatServer()

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
