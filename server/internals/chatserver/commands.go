package chatserver

import (
	"net"
	"strings"
)

// HandleCommand processes commands for user
func(cs *ChatServer) HandleCommand(conn net.Conn, command string) {
	command = strings.TrimSpace(command)
	args := strings.Split(command, " ")

	switch args[0] {
	case "exit":
		cs.RemoveClient(conn, cs.globalClients[conn])
		conn.Write([]byte("You have exited the server.\n"))
		conn.Close()
		
	}
}