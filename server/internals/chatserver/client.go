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

	// Read the mode
	mode, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Client disconnected.")
		return
	}

	mode = strings.TrimSpace(mode)

	switch mode {
	case "global":
		cs.handleGlobalChat(conn, reader)
	// Future mode handling can go here
	default:
		conn.Write([]byte("Invalid mode. Type 'global' or 'group'.\n"))
	}
}


// handleGlobalChat manages the global chat for the client
func (cs *ChatServer) handleGlobalChat(conn net.Conn, reader *bufio.Reader) {
    // Get the username from the client
    var username string
    for {
        conn.Write([]byte("Enter your username: "))
        usernameInput, err := reader.ReadString('\n')
        if err != nil {
            log.Println("Error reading username:", err)
            return
        }
        username = strings.TrimSpace(usernameInput)

        if _, exists := cs.usernames[username]; exists {
            conn.Write([]byte("Username already taken. Please enter a different username.\n"))
        } else {
            cs.AddClient(conn, username)  // Mutex is handled inside AddClient
            conn.Write([]byte(fmt.Sprintf("Your username is %s\n", username)))
            break
        }
    }

    log.Println("User", username, "has joined the chat")

    // Send the join message to users who haven't blocked this user and whom this user hasn't blocked
    for client, clientUsername := range cs.globalClients {
        if client != conn && !cs.isBlocked(clientUsername, username) && !cs.isBlocked(username, clientUsername) {
            _, err := client.Write([]byte(utils.FormatJoinMessage(username) + "\n"))
            if err != nil {
                log.Println("Error sending join message to", clientUsername, ":", err)
            }
        }
    }

    // Keep the connection open to handle incoming messages
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            log.Println("Error reading message from", username, ":", err)
            log.Println(username, "disconnected")

            // Send the leave message to users who haven't blocked this user and whom this user hasn't blocked
            for client, clientUsername := range cs.globalClients {
                if client != conn && !cs.isBlocked(clientUsername, username) && !cs.isBlocked(username, clientUsername) {
                    _, err := client.Write([]byte(utils.FormatLeaveMessage(username) + "\n"))
                    if err != nil {
                        log.Println("Error sending leave message to", clientUsername, ":", err)
                    }
                }
            }

            cs.RemoveClient(conn, username)  // Mutex is handled inside RemoveClient

            // Optional: Remove user from block lists if you want blocks to reset upon leaving
            // cs.clearBlocksForUser(username)
            return
        }

        trimmedMessage := strings.TrimSpace(message)
        if strings.HasPrefix(trimmedMessage, "/") {
            log.Println("Handling command from", username)
            // Handle command
            cs.HandleCommand(conn, trimmedMessage[1:])
        } else {
            log.Println("Broadcasting message from", username)
            // Normal Chat Message
            formattedMessage := utils.FormatChatMessage(username, trimmedMessage)
            cs.BroadcastMessage(formattedMessage, conn)
        }
    }
}


