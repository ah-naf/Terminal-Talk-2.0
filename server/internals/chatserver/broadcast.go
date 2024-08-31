package chatserver

import (
	"log"
	"net"
)

// BroadcastMessage safely sends a message to all connected clients except the sender and blocked clients
func (cs *ChatServer) BroadcastMessage(message string, sender net.Conn) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	senderUsername := cs.globalClients[sender]

	for client, username := range cs.globalClients {
		// Skip sending the message if the sender or recipient has blocked the other
		if client != sender && !cs.isBlocked(username, senderUsername) && !cs.isBlocked(senderUsername, username) {
			_, err := client.Write([]byte(message + "\n"))
			if err != nil {
				log.Println("Error broadcasting message to", username, ":", err.Error())
			}
		}
	}
}

