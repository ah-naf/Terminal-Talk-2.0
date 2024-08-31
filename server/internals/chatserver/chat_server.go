package chatserver

import (
	"net"
	"sync"
)

// ChatServer struct encapsulates the shared resources and the mutex
type ChatServer struct {
	mu            sync.Mutex
	globalClients map[net.Conn]string
	usernames     map[string]bool
	blockedBy     map[string]map[string]bool
}

// NewChatServer creates a new instance of ChatServer
func NewChatServer() *ChatServer {
	return &ChatServer{
		globalClients: make(map[net.Conn]string),
		usernames:     make(map[string]bool),
		blockedBy:     make(map[string]map[string]bool),
	}
}
