package chatserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ah-naf/chat-cli-server/internals/utils"
)

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

// HandleConnection handles the connection for each client
func (cs *ChatServer) HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	var username string
	for {
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)

		if _, exists := cs.usernames[username]; exists {
			conn.Write([]byte("Username already taken. Please enter a different username\n"))
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
	coolJoinMessage := utils.FormatJoinMessage(username)
	cs.BroadcastMessage(coolJoinMessage, conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println(username, "disconnected")
			cs.BroadcastMessage(utils.FormatLeaveMessage(username), conn)
			cs.RemoveClient(conn, username)
			return
		}

		trimmedMessage := strings.TrimSpace(message)
		if strings.HasPrefix(trimmedMessage, "/") {
			// Handle command
			cs.HandleCommand(conn, trimmedMessage[1:])
		} else {
			formattedMessage := utils.FormatChatMessage(username, strings.TrimSpace(message))
			cs.BroadcastMessage(formattedMessage, conn)
		}
	}
}
