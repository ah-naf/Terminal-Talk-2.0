package chatserver

import (
	"fmt"
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
	case "block":
		if len(args) < 2 {
			conn.Write([]byte("Usage: /block <username>\n"))
            return
		}
		cs.BlockUser(conn, args[1])

	default:
		conn.Write([]byte("Invalid command. Type 'help' for more information"))
	}
}

// Check if a user is blocked by another user
func (cs *ChatServer) isBlocked(user, target string) bool {
	if blockedList, exists := cs.blockedBy[user]; exists {
		return blockedList[target]
	}
	return false
}

// BlockUser blocks a specific user
func (cs *ChatServer) BlockUser(conn net.Conn, targetUsername string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	user, exists := cs.globalClients[conn]
	if !exists {
		conn.Write([]byte("Internal error: user not found.\n"))
		return
	}

	if user == targetUsername {
		conn.Write([]byte("You cannot block yourself.\n"))
        return
    }

	if _, exists := cs.usernames[targetUsername]; !exists {
		conn.Write([]byte(fmt.Sprintf("User '%s' does not exist.\n", targetUsername)))
		return
	}

	// Check if the target user has already blocked the current user
	if cs.isBlocked(targetUsername, user) {
		conn.Write([]byte(fmt.Sprintf("User '%s' does not exist.\n", targetUsername)))
		return
	}

	if _, blocked := cs.blockedBy[user]; !blocked {
		cs.blockedBy[user] = make(map[string]bool)
	}

	cs.blockedBy[user][targetUsername] = true
	conn.Write([]byte(fmt.Sprintf("User '%s' has been blocked.\n", targetUsername)))
}