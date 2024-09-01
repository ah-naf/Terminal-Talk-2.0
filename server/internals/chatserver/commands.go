package chatserver

import (
	"fmt"
	"net"
	"strings"

	"github.com/ah-naf/chat-cli-server/internals/utils"
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
			conn.Write([]byte(utils.FormatErrorMessage("Usage: /block <username>\n")))
            return
		}
		cs.BlockUser(conn, args[1])
	case "unblock":
		if len(args) < 2 {
			conn.Write([]byte(utils.FormatErrorMessage("Usage: /unblock <username>\n")))
            return
		}
		cs.UnblockUser(conn, args[1])

	default:
		conn.Write([]byte(utils.FormatWarningMessage("Invalid command. Type 'help' for more information.\n")))
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
		conn.Write([]byte(utils.FormatErrorMessage("Internal error: user not found.\n")))
		return
	}

	if user == targetUsername {
		conn.Write([]byte(utils.FormatWarningMessage("You cannot block yourself.\n")))
        return
    }

	if _, exists := cs.usernames[targetUsername]; !exists {
		conn.Write([]byte(utils.FormatErrorMessage(fmt.Sprintf("User '%s' does not exist.\n", targetUsername))))
		return
	}

	if _, blocked := cs.blockedBy[user]; !blocked {
		cs.blockedBy[user] = make(map[string]bool)
	}

	cs.blockedBy[user][targetUsername] = true
	conn.Write([]byte(utils.FormatSuccessMessage(fmt.Sprintf("User '%s' has been blocked.\n", targetUsername))))
}

func (cs *ChatServer) UnblockUser(conn net.Conn, targetUsername string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	user, exists := cs.globalClients[conn]
	if !exists {
		conn.Write([]byte(utils.FormatErrorMessage("Internal error: user not found.\n")))
		return
	}

	if user == targetUsername {
		conn.Write([]byte(utils.FormatWarningMessage("You cannot unblock yourself.\n")))
		return
	}

	blockedUser, blocked := cs.blockedBy[user]
	if !blocked {
		conn.Write([]byte(utils.FormatWarningMessage(fmt.Sprintf("%s is not a blocked user.\n", targetUsername))))
		return
	}

	if _, exists := blockedUser[targetUsername]; !exists {
		conn.Write([]byte(utils.FormatWarningMessage(fmt.Sprintf("%s is not a blocked user.\n", targetUsername))))
		return
	}

	// Unblock the user
	delete(blockedUser, targetUsername)
	conn.Write([]byte(utils.FormatSuccessMessage(fmt.Sprintf("User '%s' has been unblocked.\n", targetUsername))))
}
